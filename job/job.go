package job

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Job struct {
	Command string

	// jobs that this job is dependent upon
	PreJobs []string `mapstructure:"pre_jobs"`

	// jobs executed after this job successful
	PostJobs []string `mapstructure:"post_jobs"`

	// jobs executed after this job failed
	FailJobs []string `mapstructure:"fail_jobs"`

	EnvVars   map[string]string `mapstructure:"env_vars"`
	Flags     map[string]string
	Arguments []string
	WorkDir   string `mapstructure:"work_dir"`
	Stdin     string
	Stdout    string
	Stderr    string
}

var Jobs map[string]Job

func InitJobs(viper *viper.Viper) {
	if err := viper.Sub("jobs").Unmarshal(&Jobs); err != nil {
		log.Fatal("Load alamo config faild. ", err)
	}
}

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

	if len(job.Command) > 0 {
		if err = runCommand(job); err != nil {
			ExecuteJobs(job.FailJobs)
			return
		}
	}

	if err = ExecuteJobs(job.PostJobs); err != nil {
		return
	}

	log.Trace("Job finished: ", name)
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
	cmd := exec.Command(job.Command, args...)

	if len(job.WorkDir) > 0 {
		cmd.Dir = job.WorkDir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if outf, err := getIOFile(job.Stdout, job.WorkDir, true); err != nil {
		return err
	} else if outf != nil {
		defer outf.Close()
		cmd.Stdout = outf
	}

	if errf, err := getIOFile(job.Stderr, job.WorkDir, true); err != nil {
		return err
	} else if errf != nil {
		defer errf.Close()
		cmd.Stderr = errf
	}

	if inf, err := getIOFile(job.Stdin, job.WorkDir, false); err != nil {
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

func getIOFile(file string, workDir string, isOut bool) (io *os.File, err error) {
	if len(file) == 0 {
		return nil, nil
	}

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
