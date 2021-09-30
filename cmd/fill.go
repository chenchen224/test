package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gitlab.chenxk.com/test/app"
)

var fillCmd = &cobra.Command{
	Use:   "search",
	Short: "Fill all data",
	Long:  "Fill all data",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(args)
		app.FillArgs(args)
	},
}

func init() {
	rootCmd.AddCommand(fillCmd)
}
