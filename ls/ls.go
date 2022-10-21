package ls

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const dotCharacter = 46

var dotFiles = []string{".", ".."}

type Flags struct {
	L bool // ls -l
	A bool // ls -a
	R bool // ls -R
}

type LS struct {
	Dir string

	Stderr io.Writer
	StdOut io.Writer

	Flags
}

func (l *LS) ListDir() error {
	if l.R {
		return l.listDirRecursively()
	}
	return l.nonRecursiveListing()
}

func (l *LS) nonRecursiveListing() error {
	dirs, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}

	// list dotfile if -a is specified
	l.lsDotfiles()

	for _, entry := range dirs {
		if !isHiddenPath(entry.Name(), l.A) {
			fmt.Fprintln(l.StdOut, entry.Name())
		}
	}
	return nil
}

// lsDotfiles list hidden files and paths if -a is specified
func (l *LS) lsDotfiles() {
	if l.A {
		for _, file := range dotFiles {
			fmt.Fprintln(l.StdOut, file)
		}
	}
}

// listDirRecursively list all subdirectories encountered from the folder
func (l *LS) listDirRecursively() error {

	err := filepath.Walk(l.Dir,

		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if !isHiddenPath(path, l.A) {
				fmt.Fprint(l.StdOut, path, "  ")
				return err
			}
			return nil
		})

	if err != nil {

		return err
	}
	return nil
}

func isHiddenPath(path string, forceHidden bool) bool {
	if forceHidden {
		return false
	}
	return path[0] == dotCharacter
}
