package job

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Job struct {
	Name      string
	Command   string
	Arguments []string
	Flags     map[string]string
	EnvVars   map[string]string `mapstructure:"env_vars"`
	WorkDir   string            `mapstructure:"work_dir"`
	Stdin     string
	Stdout    string
	Stderr    string

	PreJobs  []string `mapstructure:"pre_jobs"`  // dependent jobs
	PostJobs []string `mapstructure:"post_jobs"` // executed when successful
	FailJobs []string `mapstructure:"fail_jobs"` // executed when failed
}

func (job *Job) Execute(jobs *Jobs) (err error) {
	log.Trace("Execute job: ", job.Name)
	if err = jobs.ExecuteJobs(job.PreJobs); err != nil {
		return
	}

	if len(job.Command) > 0 {
		if err = job.runCommand(); err != nil {
			jobs.ExecuteJobs(job.FailJobs)
			return
		}
	}

	if err = jobs.ExecuteJobs(job.PostJobs); err != nil {
		return
	}

	log.Trace("Job finished: ", job.Name)
	return
}

func (job *Job) runCommand() (err error) {
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
	cmd := exec.Command(job.Command, args...)

	if len(job.WorkDir) > 0 {
		cmd.Dir = job.WorkDir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if outf, err := job.getIOFile(job.Stdout, true); err != nil {
		return err
	} else if outf != nil {
		defer outf.Close()
		cmd.Stdout = outf
	}

	if len(job.Stderr) > 0 && job.Stdout == job.Stderr {
		cmd.Stderr = cmd.Stdout
	} else {
		if errf, err := job.getIOFile(job.Stderr, true); err != nil {
			return err
		} else if errf != nil {
			defer errf.Close()
			cmd.Stderr = errf
		}
	}

	if inf, err := job.getIOFile(job.Stdin, false); err != nil {
		return err
	} else if inf != nil {
		defer inf.Close()
		cmd.Stdin = inf
	}

	if err = cmd.Start(); err != nil {
		return
	}
	err = cmd.Wait()
	log.Trace("Command finished: ", job.Command)

	return
}

func (job *Job) getIOFile(file string, isOut bool) (io *os.File, err error) {
	if len(file) == 0 {
		return nil, nil
	}

	workDir := job.WorkDir
	if len(workDir) > 0 && !isAbsPath(file) {
		if !strings.HasSuffix(workDir, string(os.PathSeparator)) {
			workDir += string(os.PathSeparator)
		}
		file = workDir + file
	}

	if isOut {
		return os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		return os.Open(file)
	}
}

func isAbsPath(path string) bool {
	// Add paths starting with '\' or '/' on windows system
	return filepath.IsAbs(path) ||
		(runtime.GOOS == "windows" &&
			(strings.HasPrefix(path, "\\") || strings.HasPrefix(path, "/")))
}
