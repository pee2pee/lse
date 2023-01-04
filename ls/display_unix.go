//go:build !windows

package ls

import (
	"os/user"
	"strconv"
	"syscall"

	"github.com/profclems/glab/pkg/tableprinter"
)

func (l *LS) _display(table *tableprinter.TablePrinter, name string, dir *Dir) (int, error) {
	stat := dir.Info.Sys().(*syscall.Stat_t)

	usr, err := user.LookupId(strconv.Itoa(int(stat.Uid)))
	if err != nil {
		return 0, err
	}

	group, err := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
	if err != nil {
		return 0, err
	}

	timeStr := dir.Info.ModTime().UTC().Format("Jan 02 15:04")

	table.AddRow(dir.Info.Mode(), stat.Nlink, usr.Username, group.Name, l.evaluateFileAndDirSize(*dir), timeStr, name)

	return int(stat.Blocks), nil
}
