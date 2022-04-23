package cmd

import (
	"os"

	"gitee.com/logsharp/alamo/job"
	mylog "gitee.com/logsharp/alamo/log"
	homedir "github.com/mitchellh/go-homedir"
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
			log.Error("Requires at least one job\n")
			os.Exit(1)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := job.ExecuteJobs(args); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initApp)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level (fatal|error|warn|info|trace)")
}

func initApp() {
	mylog.InitLog(logLevel)
	loadCfg()
	job.InitJobs(viper.GetViper())
}

func loadCfg() {
	viper.SetConfigType("yaml")

	if len(cfgFile) > 0 {
		viper.SetConfigFile(cfgFile) // set config file directly
	} else {
		viper.SetConfigName("alamo.yml")
		viper.AddConfigPath(".") // ./alamo.yml
		if home, err := homedir.Expand("~/.alamo"); err == nil {
			viper.AddConfigPath(home)
		}

		// viper.AddConfigPath("$HOME/.alamo") // ~/.alamo/alamo.yml
		viper.AddConfigPath("/etc/alamo") // /etc/alamo/alamo.yml
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Load alamo config faild. ", err)
	} else {
		log.Info("Using config file: ", viper.ConfigFileUsed())
	}
	return
}
