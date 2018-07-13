package main

import (
	"code.linksmart.eu/dt/deployment-tool/model"
)

type registry struct {
	taskDescriptions []TaskDescription
	//tasks            []model.Task
	targets map[string]*model.Target
}

type TaskDescription struct {
	Stages     Stages
	Activation model.Activation
	Target     DeploymentTarget
	Log        model.Log

	DeploymentInfo DeploymentInfo
}

type Stages struct {
	Assemble []string
	Transfer []string
	Install  []string
	Test     []string
}

type DeploymentTarget struct {
	Tags []string
}

type DeploymentInfo struct {
	TaskID          string
	Created         string
	TransferSize    int
	MatchingTargets []string
}
