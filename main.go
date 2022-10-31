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
	cmd.PersistentFlags().BoolP("help", "", false, "help for this command")
	cmd.Flags().BoolVarP(&lsf.One, "force-entry-per-line", "1", false, "(The numeric digit “one”.) Force output to be one entry per line. This is the default when output is not to a terminal. (-l) output, and don't materialize dataless directories when listing them.")
	cmd.Flags().BoolVarP(&lsf.A, "all", "a", false, "show all files including hidden files")
	cmd.Flags().BoolVarP(&lsf.D, "directory", "d", false, "show directory structure")
	cmd.Flags().BoolVarP(&lsf.G, "group", "g", false, "group directories before files")
	cmd.Flags().BoolVarP(&lsf.L, "tabular", "l", false, "show detailed directory structure in tabular form")
	cmd.Flags().BoolVarP(&lsf.Q, "quote", "q", false, "enclose entry names in double quotes")
	cmd.Flags().BoolVarP(&lsf.R, "recursive", "R", false, "show all subdirectories encountered")
	cmd.Flags().BoolVarP(&lsf.T, "sort-by-time", "t", false, "show the files in a directory sorted by the time of modification")
	cmd.Flags().BoolVarP(&lsf.Reverse, "reverse", "r", false, "reverse the order of the sort")
	cmd.Flags().BoolVarP(&lsf.AlmostAll, "almost-all", "A", false, " do not list implied . and ..")
	cmd.Flags().BoolVarP(&lsf.X, "sort-by-extension", "X", false, "sort files by extension")
	cmd.Flags().BoolVarP(&lsf.H, "human-readable", "h", false, "print the sizes of files and directories in standardized formats like 1KB, 2GB, etc.")
	cmd.Flags().BoolVarP(&lsf.DirSize, "dir-size", "P", false, "print the sizes of files and directories in standardized formats like 1KB, 2GB, etc.")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
