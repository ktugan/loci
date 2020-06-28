package commands

import (
	"fmt"
	"os"

	"github.com/ktugan/loci/localci"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config localci.LociConfig

	configFile string
	profile    string
	command    string

	rootCmd = &cobra.Command{
		Use:   "loci",
		Short: "LoCi helps execute commands in a docker environment.",
		Long:  "An easy tool to execute builds, tests, deploys with the same behavior independent of your environment.",
		Run:   execute,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", ".loci.yml", "The config file to be used.")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "default", "The profile to run loci with.")
	rootCmd.PersistentFlags().StringVar(&command, "command", "", "The command to execute.")
}

func initConfig() {
	//fmt.Println("info: initConfig()")

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		os.Exit(1)
	}

	viper.SetEnvPrefix("LOCI")
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}

	fmt.Printf("command: %s\n", viper.GetString("command"))

	//Now unmarshal into the config itself
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func execute(cmd *cobra.Command, args []string) {
	err := localci.PrepConfig(&config)
	if err != nil {
		panic(err)
	}

	localci.Loci(config)
}

func Cli() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
