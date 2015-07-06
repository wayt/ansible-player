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
	Writer  io.Writer
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

func (j *Job) Run() error {

	workdir := fmt.Sprintf("/tmp/%x", sha1.Sum([]byte(j.Name+time.Now().String())))
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

	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()

	j.Writer.Write(output)

	return err
}

func (j *Job) command(workdir string) error {

	command := fmt.Sprintf(`cd %s && %s`, workdir, j.Command)

	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()

	j.Writer.Write(output)

	return err
}
