package cmd

import (
	"log"
	"github.com/spf13/cobra"
)

type Cmd struct {
	Dir string
}

var command Cmd

var root = &cobra.Command{
	Use:   "mio",
	Short: "I'm mio",
	Long:  "mio will provide a series of audio files!!!!!",
	Run:   defaultCmd,
}

func init () {
	defaultInit(root)
}

func defaultInit(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&command.Dir, "dir", "d", "", "provide a directory where save audio files")
	cmd.MarkFlagRequired("dir")
}

func defaultCmd(cmd *cobra.Command, args []string) {
	dir, _ := cmd.Flags().GetString("dir")
	log.Printf("mio is reading \"%s\"\n", dir)
}

func NewCmd() *Cmd {
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}

	return &command
}
