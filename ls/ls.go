package ls

import (
	"fmt"
	"github.com/pee2pee/lse/ls/color"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

	if l.T {
		return l.sortFilesByTime()
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

func (l *LS) listFileIndex() error {
	dirs, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}

	for _, file := range dirs {
		fileInfo, err := os.Stat(file.Name())
		if err != nil {
			return err
		}

		fmt.Fprintln(l.StdOut, fileInfo)
	}

	return nil
}

func (l *LS) sortFilesByTime() error {
	var modTimes []string
	var sortedFiles []string
	dirs, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}

	for _, file := range dirs {
		fileStat, err := os.Stat(file.Name())
		if err != nil {
			return err
		}
		modTimes = append(modTimes, fileStat.ModTime().String())
	}

	sort.Sort(sort.Reverse(sort.StringSlice(modTimes)))

	for _, v := range modTimes {
		for _, file := range dirs {
			fileStat, err := os.Stat(file.Name())
			if err != nil {
				return nil
			}
			if fileStat.ModTime().String() == v {
				sortedFiles = append(sortedFiles, fileStat.Name())
			}
		}
	}

	for _, v := range sortedFiles {
		fileStat, err := os.Stat(v)

		if err != nil {
			return err
		}

		fileName := fileStat.Name()
		year := fileStat.ModTime().UTC().Year()
		month := fileStat.ModTime().UTC().Month()
		day := fileStat.ModTime().UTC().Day()
		hour := fileStat.ModTime().UTC().Hour()
		minute := fileStat.ModTime().UTC().Minute()
		second := fileStat.ModTime().UTC().Second()
		timeFormat := formatTime(hour, minute, second)
		dateFormat := formatDate(year, int(month), day)

		displayFormat := fmt.Sprintf("%s Date: %s Time: %s", fileName, dateFormat, timeFormat)
		fmt.Fprintln(l.StdOut, displayFormat)
	}

	return nil
}

func formatTime(hour, minute, second int) string {
	var timeOfDay string
	if hour < 12 {
		timeOfDay = "AM"
	} else {
		timeOfDay = "PM"
	}
	return fmt.Sprintf("%s : %s : %s %s", iFormat(hour), iFormat(minute), iFormat(second), timeOfDay)
}

func formatDate(year, month, day int) string {
	return fmt.Sprintf("%s / %s / %s", iFormat(year), iFormat(month), iFormat(day))
}

func iFormat(value int) string {
	var format string
	if value < 10 {
		format = fmt.Sprintf("0%d", value)
	} else {
		format = fmt.Sprintf("%d", value)
	}
	return format
}
