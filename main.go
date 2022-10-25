package main

import (
	"os"

	"github.com/pee2pee/lse/ls"
	"github.com/pee2pee/lse/ls/color"
	"github.com/spf13/cobra"
)

func main() {
	lsf := ls.LS{
		StdOut: color.NewColorable(os.Stdout),
		Stderr: color.NewColorable(os.Stderr),
	}

	cmd := &cobra.Command{
		Use:   "lse [dir]",
		Short: "A cross platform drop-in replacement for ls",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lsf.Dir = "."
			if len(args) == 1 {
				lsf.Dir = args[0]
			}

			return lsf.ListDir()
		},
	}
	cmd.Flags().BoolVarP(&lsf.A, "all", "a", false, "show all files including hidden files")
	cmd.Flags().BoolVarP(&lsf.D, "directory", "d", false, "show directory structure")
	cmd.Flags().BoolVarP(&lsf.G, "group", "g", false, "group directories before files")
	cmd.Flags().BoolVarP(&lsf.L, "tabular", "l", false, "show detailed directory structure in tabular form")
	cmd.Flags().BoolVarP(&lsf.Q, "quote", "q", false, "enclose entry names in double quotes")
	cmd.Flags().BoolVarP(&lsf.R, "recursive", "R", false, "show all subdirectories encountered")
	cmd.Flags().BoolVarP(&lsf.T, "sort-by-time", "t", false, "show the files in a directory sorted by the time of modification")
	cmd.Flags().BoolVarP(&lsf.Reverse, "reverse", "r", false, "show the files in a directory sorted by the time of modification")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
