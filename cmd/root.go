package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configFileName string
var configFilePath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pipe",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFileName, "config_name", "c", "config", "config file yml")
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config_path", "p", "./configs", "config file path")
}

func initResource() {
	// err := cain_cfg.InitConfiguration(configFileName, strings.Split(configFilePath, ","), &cfg.CONFIG)

	// if err != nil {
	// 	panic(err)
	// }
}
