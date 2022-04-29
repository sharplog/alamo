package job

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Jobs map[string]*Job

func (jobs *Jobs) Init(viper *viper.Viper) {
	if err := viper.Sub("jobs").Unmarshal(jobs); err != nil {
		log.Fatal("Load alamo config faild. ", err)
	}
	for k, v := range *jobs {
		v.Name = k
	}
}

func (jobs *Jobs) ExecuteJobs(names []string) (err error) {
	for _, name := range names {
		if job, ok := (*jobs)[name]; !ok {
			return errors.New("Job does't exist: " + name)
		} else {
			job.Execute(jobs)
		}
	}
	return
}
