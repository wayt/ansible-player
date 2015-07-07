package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/wayt/happyngine/env"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type Job struct {
	Name    string
	Git     string
	Command string
	logFile string
	log     io.WriteCloser
}

func GetJob(name string) (*Job, error) {

	data, err := ioutil.ReadFile(env.Get("JOB_FILE"))
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

func (j *Job) Logs() ([]byte, error) {

	return ioutil.ReadFile(j.logFile)
}

func (j *Job) Run() error {

	buildId := fmt.Sprintf("%x", sha1.Sum([]byte(j.Name+time.Now().String())))

	j.logFile = fmt.Sprintf("%s/%s.log", env.Get("LOG_DIR"), buildId)
	log, err := os.OpenFile(j.logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer log.Close()

	j.log = log

	fmt.Fprintf(log, "New build %s - %s\n", buildId, time.Now().String())

	workdir := fmt.Sprintf("/tmp/%s", buildId)
	defer os.RemoveAll(workdir)

	if err := j.clone(workdir); err != nil {
		return err
	}

	if err := j.command(workdir); err != nil {
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
