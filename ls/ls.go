package ls

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pee2pee/lse/ls/color"
)

const dotCharacter = 46

var dotFiles = []string{".", ".."}

type Flags struct {
	A bool // ls -a
	D bool // ls -d
	G bool // ls --group
	L bool // ls -l
	R bool // ls -R
}

type LS struct {
	Dir string

	Stderr io.Writer
	StdOut io.Writer
	Color  *color.Palette

	Flags
}

func (l *LS) ListDir() error {
	if l.D {
		return l.showDirStructure()
	}

	if l.G {
		return l.groupdirfirst()
	}

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

func (l *LS) groupdirfirst() error {
	var dirs []string
	var filedirs []string
	files, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else {
			filedirs = append(filedirs, file.Name())
		}
	}
	for _, isDirs := range dirs {
		fmt.Fprintln(l.StdOut, isDirs)
	}
	for _, isFiles := range filedirs {
		fmt.Fprintln(l.StdOut, isFiles)
	}

	return nil
}

func (l *LS) showDirStructure() error {
	file, err := os.Stat(l.Dir)
	if err != nil {
		return err
	}

	p := strings.TrimSuffix(l.Dir, "/")
	if file.IsDir() {
		p = p + "/"
	}

	fmt.Fprintln(l.StdOut, p)
	return nil
}
