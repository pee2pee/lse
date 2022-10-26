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

type Flags struct {
	A         bool // ls -a
	D         bool // ls -d
	G         bool // ls --group
	L         bool // ls -l
	Q         bool // ls --quote
	R         bool // ls -R
	T         bool // ls -t
	Reverse   bool // ls -r
	AlmostAll bool // ls -A
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

type Dirs []Dir

func (d Dirs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Dirs) Len() int {
	return len(d)
}

func (d Dirs) Less(i, j int) bool {
	return true
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
	var d Dirs

	// list dotfile if -a is specified
	if l.A {
		var dotFiles = []string{".", ".."}
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
		if !isHiddenPath(entry.Name(), l.A || l.AlmostAll) {
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

	if l.G {
		dirs, fileDirs := getFilesAndDirs(d)
		if l.T {
			sort.SliceStable(dirs, func(i, j int) bool {
				return dirs[i].Info.ModTime().After(dirs[j].Info.ModTime())
			})

			sort.SliceStable(fileDirs, func(i, j int) bool {
				return fileDirs[i].Info.ModTime().After(fileDirs[j].Info.ModTime())
			})
		}
		d = append(dirs, fileDirs...)
	}

	if l.Reverse {
		sort.Sort(sort.Reverse(d))
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
		if dir.IsDir() && !isHiddenPath(dir.Name(), l.A || l.AlmostAll) {
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

func getFilesAndDirs(d []Dir) (dirs []Dir, fileDirs []Dir) {
	for _, file := range d {
		if file.Info.IsDir() {
			dirs = append(dirs, file)
		} else {
			fileDirs = append(fileDirs, file)
		}
	}
	return dirs, fileDirs
}
