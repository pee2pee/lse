package ls

import (
	"fmt"
	"github.com/pee2pee/lse/ls/color"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const dotCharacter = 46

var dotFiles = []string{".", ".."}

type Flags struct {
	A bool // ls -a
	D bool // ls -d
	G bool // ls --group
	L bool // ls -l
	Q bool // ls --quote
	R bool // ls -R
}

type LS struct {
	Dir string

	Stderr io.Writer
	StdOut io.Writer
	Color  *color.Palette

	Flags
}

type Dir struct {
	Path string
	Info fs.FileInfo
}

func (l *LS) ListDir() error {
	if l.D {
		return l.showDirStructure()
	}

	if l.G {
		return l.groupdirfirst()
	}
	if l.Q {
		return l.qoutesEntryNames()
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
	var d []Dir

	// list dotfile if -a is specified
	if l.A {
		for _, file := range dotFiles {
			stat, err := os.Stat(file)
			if err != nil {
				return err
			}
			d = append(d, Dir{
				Info: stat,
				Path: file,
			})
		}
	}

	for _, entry := range dirs {
		if !isHiddenPath(entry.Name(), l.A) {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			d = append(d, Dir{
				Info: info,
				Path: filepath.Join(l.Dir, info.Name()),
			})
		}
	}
	return l.display(d)
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
func (l *LS) qoutesEntryNames() error {
	files, err := os.ReadDir(l.Dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(strconv.Quote(file.Name()))
	}
	return nil
}
