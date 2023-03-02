package scenario

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
	hf "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ParseScenario(name string, namespace string, path string) (s *hf.Scenario, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	scenarioSpec, err := processScenarioYAML(absPath)
	if err != nil {
		return s, err
	}

	steps, err := processContents(absPath)

	if err != nil {
		return s, err
	}

	scenarioSpec.Steps = steps
	s = &hf.Scenario{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: *scenarioSpec,
	}

	annotations := make(map[string]string)
	annotations["managedBy"] = "hfcli"
	s.Annotations = annotations
	s.Spec.Id = name
	if s.Spec.Name == "" {
		s.Spec.Name = base64.StdEncoding.EncodeToString([]byte(name))
	}

	if s.Spec.KeepAliveDuration == "" {
		s.Spec.KeepAliveDuration = DefaultKeepAliveDuration
	}
	return s, nil
}

func processScenarioYAML(absPath string) (s *hf.ScenarioSpec, err error) {
	scenarioFilePath := filepath.Join(absPath, "scenario.yml")
	_, err = os.Stat(scenarioFilePath)
	if err != nil {
		return s, err
	}

	scenarioFileContent, err := os.ReadFile(scenarioFilePath)
	if err != nil {
		return nil, err
	}

	s = &hf.ScenarioSpec{}
	err = yaml.Unmarshal(scenarioFileContent, s)
	// need to b64 encode the name and description
	s.Name = base64.StdEncoding.EncodeToString([]byte(s.Name))
	s.Description = base64.StdEncoding.EncodeToString([]byte(s.Description))
	return s, err
}

func processContents(absPath string) (steps []hf.ScenarioStep, err error) {
	contentDir := filepath.Join(absPath, "content")
	contentDirPath, err := os.Stat(contentDir)
	if err != nil {
		return steps, err
	}

	if !contentDirPath.IsDir() {
		err = fmt.Errorf("%s is not a directory", contentDirPath)
		return steps, err
	}

	files, err := os.ReadDir(contentDir)
	if err != nil {
		return steps, err
	}

	steps, err = readFiles(contentDir, files)

	return steps, err
}

func readFiles(path string, files []os.DirEntry) (steps []hf.ScenarioStep, err error) {

	var filesWithContent []FilenameWithContent

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(path, file.Name())
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return steps, err
			}
			f := FilenameWithContent{
				FileName: filePath,
				Content:  fileContent,
			}

			filesWithContent = append(filesWithContent, f)
		}
	}

	stepsWithID := []StepWithID{}
	for _, f := range filesWithContent {
		s, err := parseToml(f.Content, f.FileName)
		if err != nil {
			return steps, err
		}

		stepsWithID = append(stepsWithID, s)
	}

	steps = sortContent(stepsWithID)
	return steps, nil
}

type stepMeta struct {
	Title  string `toml:"title"`
	Weight *int   `toml:"weight"`
}

func parseToml(content []byte, fileName string) (s StepWithID, err error) {
	// empty defaults
	title := extractFilename(fileName)
	tw := 1000
	tmp := stepMeta{}
	frontMatter, noTomlContent := extractTOML(content)

	if len(frontMatter) != 0 {
		if _, err := toml.Decode(string(frontMatter), &tmp); err != nil {
			return s, err
		}
	}

	if tmp.Weight != nil {
		tw = *tmp.Weight
	}

	s.Weight = tw
	if tmp.Title != "" {
		title = tmp.Title
	}
	s.Step.Title = base64.StdEncoding.EncodeToString([]byte(title))
	s.Step.Content = base64.StdEncoding.EncodeToString(noTomlContent)
	return s, nil
}

func sortContent(stepsWithID []StepWithID) (steps []hf.ScenarioStep) {
	sort.SliceStable(stepsWithID, func(i, j int) bool {
		return stepsWithID[i].Weight < stepsWithID[j].Weight
	})

	for _, stepWithID := range stepsWithID {
		steps = append(steps, stepWithID.Step)
	}
	return steps
}

func extractTOML(content []byte) (toml []byte, noToml []byte) {
	r := regexp.MustCompile(`(?s)\+\+\+(.*)\+\+\+`)
	tmp := r.Find(content)
	r2 := regexp.MustCompile(`\+\+\+`)
	toml = r2.ReplaceAll(tmp, []byte(""))
	noToml = r.ReplaceAll(content, []byte(""))
	return toml, []byte(strings.Trim(string(noToml), "\n"))
}

func extractFilename(pathToFile string) (fileName string) {
	fileNameArr := strings.Split(pathToFile, "/")
	fileName = fileNameArr[len(fileNameArr)-1]
	return fileName
}

func DumpScenario(s *hf.Scenario, path string) (err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	logrus.Infof("writing scenario %s to %s", s.Spec.Name, absPath)

	err = writeScenarioYAML(s, absPath)
	if err != nil {
		return err
	}

	err = writeContents(s, absPath)

	if err != nil {
		return err
	}

	return nil
}

func writeScenarioYAML(s *hf.Scenario, path string) error {
	logrus.Infof("creating scenario.yml")

	scenarioFilePath := filepath.Join(path, "scenario.yml")

	spec := s.Spec

	// need to b64 decode the name and description
	decodedName, err := base64.StdEncoding.DecodeString(s.Spec.Name)
	if err != nil {
		return err
	}
	spec.Name = string(decodedName)
	decodedDescription, err := base64.StdEncoding.DecodeString(s.Spec.Description)
	if err != nil {
		return err
	}
	spec.Description = string(decodedDescription)
	spec.Steps = nil

	scenarioFileContent, err := yaml.Marshal(spec)

	if err != nil {
		return err
	}

	return os.WriteFile(scenarioFilePath, scenarioFileContent, 0666)
}

func writeContents(s *hf.Scenario, path string) error {
	contentDir := filepath.Join(path, "content")
	contentDirPath, err := os.Stat(contentDir)
	if err != nil {
		err = os.Mkdir(contentDir, 0755)
		if err != nil {
			return err
		}
		contentDirPath, err = os.Stat(contentDir)
		if err != nil {
			return err
		}
	}

	if !contentDirPath.IsDir() {
		err = fmt.Errorf("%s is not a directory", contentDirPath)
		return err
	}

	files, err := os.ReadDir(contentDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			err = os.Remove(filepath.Join(contentDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	for i, step := range s.Spec.Steps {
		weight := i + 1
		filename := filepath.Join(contentDir, fmt.Sprintf("step-%d.md", weight))
		logrus.Infof("creating step %s", filename)

		title, err := base64.StdEncoding.DecodeString(step.Title)
		if err != nil {
			return err
		}

		stepMeta := stepMeta{}
		stepMeta.Title = string(title)
		stepMeta.Weight = &weight

		var firstBuffer bytes.Buffer
		encoder := toml.NewEncoder(&firstBuffer)
		err = encoder.Encode(stepMeta)
		if err != nil {
			return err
		}

		spacer := []byte("+++\n")

		stepFileContent, err := base64.StdEncoding.DecodeString(step.Content)

		if err != nil {
			return err
		}

		content := append(spacer, firstBuffer.Bytes()...)
		content = append(content, spacer...)
		content = append(content, "\n"...)
		content = append(content, stepFileContent...)
		content = append(content, "\n"...)

		err = os.WriteFile(filename, content, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
