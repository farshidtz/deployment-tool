package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"code.linksmart.eu/dt/deployment-tool/manager/model"
	"code.linksmart.eu/dt/deployment-tool/manager/source"
	"github.com/pbnjay/memory"
)

type installer struct {
	logger   chan<- model.Log
	executor *executor
}

func newInstaller(logger chan<- model.Log) installer {
	return installer{
		logger: logger,
	}
}

func (i *installer) evaluate(ann *model.Announcement) bool {
	sizeLimit := memory.TotalMemory() / 2 // TODO calculate this based on the available memory
	return uint64(ann.Size) <= sizeLimit
}

func (i *installer) store(artifacts []byte, taskID string, debug bool) {
	defer func() {
		artifacts = nil // release memory
	}()

	taskDir := fmt.Sprintf("%s/tasks/%s", WorkDir, taskID)
	log.Println("installer: Task work directory:", taskDir)

	// nothing to store
	if len(artifacts) == 0 {
		log.Printf("installer: Nothing to store.")
		// create task with source directory
		err := os.MkdirAll(fmt.Sprintf("%s/%s", taskDir, source.SourceDir), 0755)
		if err != nil {
			i.sendLogFatal(taskID, fmt.Sprintf("error creating source directory: %s", err))
			return
		}
		return
	}

	// create task directory
	err := os.MkdirAll(taskDir, 0755)
	if err != nil {
		i.sendLogFatal(taskID, fmt.Sprintf("error creating task directory: %s", err))
		return
	}

	// decompress and store
	log.Printf("installer: Deploying %d bytes of artifacts.", len(artifacts))
	err = model.DecompressFiles(artifacts, taskDir)
	if err != nil {
		i.sendLogFatal(taskID, fmt.Sprintf("error reading archive: %s", err))
		return
	}
	i.sendLog(taskID, fmt.Sprintf("decompressed archive of %d bytes", len(artifacts)), false, debug)
	//i.sendLog(taskID, stage, model.StageEnd, false, debug)
}

func (i *installer) install(commands []string, taskID string, debug bool) bool {
	// nothing to execute
	if len(commands) == 0 {
		log.Printf("installer: Nothing to execute.")
		i.sendLog(taskID, model.StageEnd, false, debug)
		return true
	}

	log.Printf("installer: Installing task: %s", taskID)
	//i.sendLog(taskID, stage, model.StageStart, false, debug)

	// execute sequentially, return if one fails
	i.executor = newExecutor(taskID, model.StageInstall, i.logger, debug)
	for _, command := range commands {
		success := i.executor.execute(command)
		if !success {
			i.sendLogFatal(taskID, "ended with errors")
			return false
		}
	}

	log.Printf("installer: Install ended.")
	i.sendLog(taskID, model.StageEnd, false, debug)
	return true
}

func (i *installer) sendLog(task, output string, error bool, debug bool) {
	i.logger <- model.Log{task, model.StageInstall, "", output, error, model.UnixTime(), debug}
}

func (i *installer) sendLogFatal(task, output string) {
	log.Printf("installer: %s", output)
	if output != "" {
		i.sendLog(task, output, true, true)
	}
	i.sendLog(task, model.StageEnd, true, true)
}

// clean removed old task directory
func (i *installer) clean(taskID string) {
	log.Println("installer: Removing files for task:", taskID)

	wd := fmt.Sprintf("%s/tasks", WorkDir)

	_, err := os.Stat(wd)
	if err != nil && os.IsNotExist(err) {
		// nothing to remove
		return
	}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Printf("installer: Error reading work dir: %s", err)
		return
	}
	for i := 0; i < len(files); i++ {
		if files[i].Name() != taskID {
			log.Println(files[i].Name(), taskID)
			filename := fmt.Sprintf("%s/%s", wd, files[i].Name())
			log.Printf("installer: Removing: %s", filename)
			err = os.RemoveAll(filename)
			if err != nil {
				log.Printf("installer: Error removing: %s", err)
			}
		}
	}
}

func (r *installer) stop() bool {
	log.Println("installer: Shutting down...")
	success := r.executor.stop()
	return success
}
