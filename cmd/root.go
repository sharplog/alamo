package cmd

import (
	"errors"
	"os"

	"gitee.com/logsharp/alamo/job"
	mylog "gitee.com/logsharp/alamo/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   "alamo [flags] job [job] ...",
	Short: "Execute a job from configuration",
	Long:  "See https://gitee.com/logsharp/alamo for documentation.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires at least one job")
		}
		for _, name := range args {
			if _, ok := job.Jobs[name]; !ok {
				return errors.New("Job does't exist: " + name)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return job.ExecuteJobs(args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initApp)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level (fatal|error|warn|info|trace)")
}

func initApp() {
	mylog.InitLog(logLevel)
	loadCfg()
}

func loadCfg() {
	viper.SetConfigType("yaml")

	if len(cfgFile) > 0 {
		viper.SetConfigFile(cfgFile) // set config file directly
	} else {
		viper.SetConfigName("alamo")
		viper.AddConfigPath(".")            // ./alamo.yml
		viper.AddConfigPath("$HOME/.alamo") // ~/.alamo/alamo.yml
		viper.AddConfigPath("/etc/alamo")   // /etc/alamo/alamo.yml
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Load alamo config faild. ", err)
	} else {
		log.Info("Using config file: ", viper.ConfigFileUsed())
		if err := viper.Sub("jobs").Unmarshal(&job.Jobs); err != nil {
			log.Fatal("Load alamo config faild. ", err)
		}
	}
	return
}
