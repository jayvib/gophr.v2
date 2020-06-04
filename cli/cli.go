package cli

import (
	"fmt"
	"github.com/jayvib/golog"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string
var debug bool

var GophrApp = &cobra.Command{
	Use:   "gophr",
	Short: "gophr is an CLI for interacting gophr services",
	Long:  "gophr is a CLI for interacting gophr services",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			golog.SetLevel(golog.DebugLevel)
			golog.Warning("GOPHR: DEBUG MODE")
		}
	},
}

func Execute() {
	if err := GophrApp.Execute(); err != nil {
		golog.Fatal(err)
	}
}

func AddCommands(cmds ...*cobra.Command) {
	GophrApp.AddCommand(cmds...)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	GophrApp.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gophr.yaml)")

	GophrApp.PersistentFlags().BoolVar(&debug, "debug", false, "debugging mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	GophrApp.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gophr.v2" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gophr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
