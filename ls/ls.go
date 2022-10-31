package ls

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	One       bool // ls -1
	X         bool // ls -X
	H         bool //ls -h
	DirSize   bool // new
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

	if l.X {
		dirs, fileDirs := getFilesAndDirs(d)
		sort.SliceStable(fileDirs, func(i, j int) bool {
			return filepath.Ext(fileDirs[i].Info.Name()) < filepath.Ext(fileDirs[j].Info.Name())
		})
		d = append(dirs, fileDirs...)
	}

	if l.DirSize {
		// dirs, _ := getFilesAndDirs(d)
		// dirSizes := make(map[string]int64)
		// for _, dir := range dirs {
		// 	dirSize := calculateDirSize(dir.Info.Name(), "")
		// 	// dirSizes[dir.Info.Name()] = dirSize
		// }

		// fmt.Fprintln(l.StdOut, dirSizes)
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

func getFilesAndDirs(d []Dir) (dirs, fileDirs []Dir) {
	for _, file := range d {
		if file.Info.IsDir() {
			dirs = append(dirs, file)
		} else {
			fileDirs = append(fileDirs, file)
		}
	}
	return dirs, fileDirs
}

func (l *LS) minifySize(size int64) (sizeString string) {
	if !l.H {
		return strconv.Itoa(int(size))
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	if size == 0 {
		sizeString = fmt.Sprintf("%d%s", size, units[size])
		return
	}

	const expValue int = 3

	for i := len(units) - 1; i >= 0; i-- {
		divisor := math.Pow(10, float64(expValue*i))
		result := float64(size) / divisor
		if result > 1.0 {
			if int64(result) == size {
				sizeString = fmt.Sprintf("%d%s", size, units[i])
			} else {
				sizeString = fmt.Sprintf("%.1f%s", result, units[i])
			}
			break
		}
	}
	return
}

func calculateDirSize(dirPath, currentDir string) error {
	// cMap := make(map[string][]fs.DirEntry)
	// dirs, err := os.ReadDir(dirPath)
	// if err != nil {
	// 	return err
	// }

	// for _, dir := range dirs {
	// 	if dir.IsDir() {
	// 		p := filepath.Join(currentDir, dir.Name())
	// 		cMap[p] = calculateDirSize(p)
	// 	}
	// // }

	return nil
}

func calculateSubDirSize(dirs []fs.DirEntry) error {
	var d []Dir
	for _, entry := range dirs {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		d = append(d, Dir{
			Info: info,
			// Path: filepath.Join(l.Dir, info.Name()),
			Path: "",
		})
	}
	return nil
}

func calculateSubFileSize(childFile, parent string) (size int64) {
	return
}
