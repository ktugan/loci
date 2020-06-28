package main

import (
	"fmt"

	"github.com/ktugan/loci"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "loci",
		Short: "LoCi helps execute commands in a docker environment.",
		Long:  "An easy tool to execute builds, tests, deploys with the same behavior independent of your environment.",
		Run:   execute,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", ".loci.yml", "The config file to be used.")

	//fmt.Println("info: init()")
}

func initConfig() {
	//fmt.Println("info: initConfig()")

	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}
}

func execute(cmd *cobra.Command, args []string) {
	config := loci.LoadConfig(configFile)

	err := loci.PrepConfig(&config)
	if err != nil {
		panic(err)
	}

	loci.Loci(config)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
