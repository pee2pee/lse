package ls

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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
	Q bool // ls --quote
	R bool // ls -R
	T bool // ls -t
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

	if l.R {
		return l.listDirRecursively("")
	}
	return l.nonRecursiveListing()
}

func (l *LS) listDir(dirs []fs.DirEntry) error {
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

	if l.T {
		sort.SliceStable(d, func(i, j int) bool {
			return d[i].Info.ModTime().After(d[j].Info.ModTime())
		})
	}

	if l.G && !l.T {
		var dirs []Dir
		var fileDirs []Dir
		for _, file := range d {
			if file.Info.IsDir() {
				dirs = append(dirs, file)
			} else {
				fileDirs = append(fileDirs, file)
			}
		}
		d = append(dirs, fileDirs...)
	}

	return l.display(d)
}

func (l *LS) nonRecursiveListing() error {
	dirs, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}
	return l.listDir(dirs)
}

// listDirRecursively list all subdirectories encountered from the folder
func (l *LS) listDirRecursively(path string) error {
	if path == "" {
		path = l.Dir
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	err = l.listDir(dirs)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir.IsDir() && !isHiddenPath(dir.Name(), l.A) {
			p := filepath.Join(path, dir.Name())

			fmt.Fprintf(l.StdOut, "\n%s:\n", p)

			err := l.listDirRecursively(p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isHiddenPath(path string, forceHidden bool) bool {
	if forceHidden {
		return false
	}
	return path[0] == dotCharacter
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
