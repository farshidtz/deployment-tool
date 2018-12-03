package model

const (
	// Request types
	RequestTargetAll = "ALL"
	RequestTargetID  = "ID"
	RequestTargetTag = "TAG"
	// Stage types
	StageAssemble = "ASSEMBLE"
	StageTransfer = "TRANSFER"
	StageInstall  = "INSTALL"
	StageTest     = "TEST"
	StageRun      = "RUN"
	// Other consts
	PrefixSeparator = "-"
)

type Stages struct {
	//Assemble []string
	Transfer []string `json:"transfer"`
	Install  []string `json:"install"`
	//Test     []string
	Run []string `json:"run"`
}

// Header contains information that is common among task related structs
type Header struct {
	ID    string `json:"id"`
	Debug bool   `json:"debug"`
	Created  int64  `json:"created"`
}

// Announcement carries information about a task
type Announcement struct {
	Header
	Size int `json:"size"`
}

// Task is a struct with all the information for deployment on a target
type Task struct {
	Header
	Stages    Stages `json:"stages"`
	Artifacts []byte `json:"artifacts"`
}

type LogRequest struct {
	IfModifiedSince UnixTimeType
}

func TargetTopic(id string) string {
	return RequestTargetID + PrefixSeparator + id
}

func TargetTag(tag string) string {
	return RequestTargetTag + PrefixSeparator + tag
}
