//go:build windows

package ls

import (
	"github.com/profclems/glab/pkg/tableprinter"
	"log"
	"os"
	"os/user"
	"syscall"
)

func (l *LS) _display(table *tableprinter.TablePrinter, name string, dir *Dir) (int, error) {
	//stat := dir.Info.Sys().(*syscall.Win32FileAttributeData)
	usr, err := user.Current()
	if err != nil {
		return 0, err
	}

	var data syscall.ByHandleFileInformation

	f, err := os.Open(dir.Path)
	if err != nil {
		return 0, err
	}
	log.Println(f.Fd(), dir.Path)

	err = syscall.GetFileInformationByHandle(syscall.Handle(f.Fd()), &data)
	if err != nil {
		return 0, err
	}

	blocks := data.FileSizeLow
	size := uint64(data.FileSizeHigh)<<32 | uint64(data.FileSizeLow)
	//nlink := 1

	timeStr := dir.Info.ModTime().UTC().Format("Jan 02 15:04")

	table.AddRow(dir.Info.Mode(), data.NumberOfLinks, usr.Username, usr.Gid, l.minifySize(int64(size)), timeStr, name)

	return int(blocks), nil
}
