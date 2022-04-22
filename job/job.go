package job

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Job struct {
	Command   string
	PreJobs   []string          `mapstructure:"pre_jobs"`
	PostJobs  []string          `mapstructure:"post_jobs"`
	EnvVars   map[string]string `mapstructure:"env_vars"`
	Flags     map[string]string
	Arguments []string
}

var Jobs map[string]Job

func ExecuteJobs(names []string) (err error) {
	for _, name := range names {
		if err = executeJob(name); err != nil {
			return
		}
	}
	return
}

func executeJob(name string) (err error) {
	job, ok := Jobs[name]
	if !ok {
		return errors.New("Job does't exist: " + name)
	}

	log.Trace("Execute job: ", name)
	if err = ExecuteJobs(job.PreJobs); err != nil {
		return
	}

	if err = runCommand(job); err != nil {
		return
	}

	if err = ExecuteJobs(job.PostJobs); err != nil {
		return
	}

	return
}

func runCommand(job Job) (err error) {
	var args []string

	for k, v := range job.Flags {
		if len(v) > 0 {
			args = append(args, k, v)
		} else {
			args = append(args, k)
		}
	}

	for _, v := range job.Arguments {
		args = append(args, v)
	}

	for k, v := range job.EnvVars {
		os.Setenv(k, v)
	}

	log.Trace("Run command: ", job.Command)

	// 非文件命令的执行，加cmd /c
	cmd := exec.Command(job.Command, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 错误输出

	err = cmd.Run()
	if stdout.Len() > 0 {
		fmt.Print(string(stdout.Bytes()))
	}
	if stderr.Len() > 0 {
		fmt.Fprintf(os.Stderr, string(stderr.Bytes()))
	}

	return
}
