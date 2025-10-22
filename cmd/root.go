package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

type Cmd struct {}

var command Cmd

var root = &cobra.Command{
	Use:   "mio",
	Short: "I'm mio",
	Long:  "mio will provide a series of audio files!!!!!",
	Run:   defaultCmd,
}

func init() {
	defaultInit(root)
}

func defaultInit(cmd *cobra.Command) {
	
}

func defaultCmd(cmd *cobra.Command, args []string) {
	
}

func NewCmd() *Cmd {
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}

	return &command
}
