package job

import (
	"errors"
	"os"
	"os/exec"
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

	return exec.Command(job.Command, args...).Run()
}

// TODO
// 1. 非文件命令的执行
// 2. 命令的输出
