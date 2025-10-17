package cmd

import (
	"log"
	"path/filepath"
	"os"

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

func init() {
	defaultInit(root)
}

func defaultInit(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&command.Dir, "dir", "d", "", "provide a directory where save audio files")
	err := cmd.MarkFlagRequired("dir")
	if err != nil {
		log.Fatal(err)
	}
}

func validateDir(dir string) string {
	isAbs := filepath.IsAbs(dir)
	if !isAbs {
		return ""
	}

	cleanPath, err := filepath.Abs(filepath.Clean(dir))
	if err != nil {
		return ""
	}

	_, err = os.Stat(cleanPath)
	if err == nil {
		return cleanPath
	}
	
	return ""
}

func defaultCmd(cmd *cobra.Command, args []string) {
	dir, _ := cmd.Flags().GetString("dir")
	dir = validateDir(dir)
	if dir == "" {
		log.Fatal("dir is invalid")
	}
	err := cmd.Flags().Set("dir", dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("mio is reading \"%s\"\n", dir)
}

func NewCmd() *Cmd {
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}

	return &command
}
