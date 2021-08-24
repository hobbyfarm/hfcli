package scenario

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
	s.Spec.Name = base64.StdEncoding.EncodeToString([]byte(name))

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

	scenarioFileContent, err := ioutil.ReadFile(scenarioFilePath)
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
			fileContent, err := ioutil.ReadFile(filePath)
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

func parseToml(content []byte, fileName string) (s StepWithID, err error) {
	type obj struct {
		Title  string `toml:"title"`
		Weight *int   `toml:"weight"`
	}

	// empty defaults
	title := extractFilename(fileName)
	tw := 1000
	tmp := obj{}
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
	r := regexp.MustCompile("(?s)\\+\\+\\+(.*)\\+\\+\\+")
	tmp := r.Find(content)
	r2 := regexp.MustCompile("\\+\\+\\+")
	toml = r2.ReplaceAll(tmp, []byte(""))
	noToml = r.ReplaceAll(content, []byte(""))
	return toml, noToml
}

func extractFilename(pathToFile string) (fileName string) {
	fileNameArr := strings.Split(pathToFile, "/")
	fileName = fileNameArr[len(fileNameArr)-1]
	return fileName
}
