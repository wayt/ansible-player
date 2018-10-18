package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	envJobFile      = "JOB_FILE"
	envLogDirectory = "LOG_DIR"
)

type Job struct {
	Name    string `json:"name"`
	JobId   string `json:"job_id"`
	Error   error  `json:"error"`
	Git     string `json:"-"`
	Command string `json:"-"`
	logFile string
	log     io.WriteCloser
	Logs    string `json:"logs"`
}

func GetJob(name string) (*Job, error) {

	data, err := ioutil.ReadFile(os.Getenv(envJobFile))
	if err != nil {
		return nil, err
	}

	jobs := make(map[string]Job)
	if err := yaml.Unmarshal(data, jobs); err != nil {
		return nil, err
	}

	for jobName, job := range jobs {
		if jobName == name {
			job.Name = jobName
			return &job, nil
		}
	}

	return nil, nil
}

func GetJobLogs(id string) ([]byte, error) {

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.log", os.Getenv(envLogDirectory), id))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

func (j *Job) parseLogs() {

	data, _ := ioutil.ReadFile(j.logFile)

	j.Logs = string(data)
}

func (j *Job) Run() error {

	defer j.parseLogs()

	j.JobId = fmt.Sprintf("%s-%x", j.Name, sha1.Sum([]byte(j.Name+time.Now().String())))

	j.logFile = fmt.Sprintf("%s/%s.log", os.Getenv(envLogDirectory), j.JobId)
	log, err := os.OpenFile(j.logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		j.Error = err
		return err
	}
	defer log.Close()

	j.log = log

	fmt.Fprintf(log, "New build %s - %s\n", j.JobId, time.Now().String())

	workdir := fmt.Sprintf("/tmp/%s", j.JobId)
	defer os.RemoveAll(workdir)

	if err := j.clone(workdir); err != nil {
		j.Error = err
		return err
	}

	if err := j.command(workdir); err != nil {
		j.Error = err
		return err
	}

	return nil
}

func (j *Job) clone(workdir string) error {

	command := fmt.Sprintf(`git clone -v %s %s`, j.Git, workdir)

	fmt.Fprintf(j.log, "sh -c '%s'\n", command)
	cmd := exec.Command("sh", "-c", command)

	cmd.Stdout = j.log
	cmd.Stderr = j.log

	return cmd.Run()
}

func (j *Job) command(workdir string) error {

	command := fmt.Sprintf(`cd %s && %s`, workdir, j.Command)

	fmt.Fprintf(j.log, "sh -c '%s'\n", command)
	cmd := exec.Command("sh", "-c", command)

	cmd.Stdout = j.log
	cmd.Stderr = j.log

	return cmd.Run()
}
