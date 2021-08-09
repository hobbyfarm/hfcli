package scenario

import hf "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"

// StepWithID is an intermediate structure to read contents and order them

type StepWithID struct {
	Weight int
	Step   hf.ScenarioStep
}

type FilenameWithContent struct {
	FileName string
	Content  []byte
}

const (
	DefaultKeepAliveDuration = "10m"
)