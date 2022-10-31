//go:build windows

package ls

import (
	"os/user"
	"syscall"

	"github.com/profclems/glab/pkg/tableprinter"
)

func (l *LS) _display(table *tableprinter.TablePrinter, name string, dir *Dir) (int, error) {
	stat := dir.Info.Sys().(*syscall.Win32FileAttributeData)
	usr, err := user.Current()
	if err != nil {
		return 0, err
	}

	blocks := stat.FileSizeLow
	nlink := 1

	timeStr := dir.Info.ModTime().UTC().Format("Jan 02 15:04")

	table.AddRow(dir.Info.Mode(), nlink, usr.Username, usr.Gid, l.evaluateFileAndDirSize(*dir), timeStr, name)

	return int(blocks), nil
}
