package ls

import (
	"fmt"
	"io"
	"os"
)

const dotCharacter = 46

var dotFiles = []string{".", ".."}

type Flags struct {
	L bool // ls -l
	A bool // ls -a
	G bool // ls --group
}

type LS struct {
	Dir string

	Stderr io.Writer
	StdOut io.Writer

	Flags
}

func (l *LS) ListDir() error {
	if l.G {
		return l.groupdirfirst()
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
		fmt.Println(isDirs)
	}
	for _, isFiles := range filedirs {
		fmt.Println(isFiles)
	}
	return nil
}
