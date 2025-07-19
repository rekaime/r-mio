package cmd

import (
    "github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   	"server",
	Short: 	"run the server",
	Long:  	`mio will run the server`,
	Run: 	defaultCmd,
}

func init() {
	root.AddCommand(serverCmd)
	defaultInit(serverCmd)
}