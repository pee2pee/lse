//go:build windows

package ls

import (
	"os"
	"os/user"
	"syscall"
	"unsafe"

	"github.com/profclems/glab/pkg/tableprinter"
)

func (l *LS) _display(table *tableprinter.TablePrinter, name string, dir *Dir) (int, error) {
	//stat := dir.Info.Sys().(*syscall.Win32FileAttributeData)
	usr, err := user.Current()
	if err != nil {
		return 0, err
	}

	var data syscall.ByHandleFileInformation
	handle, err := syscall.Open(dir.Path, os.O_RDONLY, 0600)
	if err != nil {
		return 0, err
	}

	err = getFileInformationByHandle(handle, &data)
	if err != nil {
		return 0, err
	}

	blocks := data.FileSizeLow
	//nlink := 1

	timeStr := dir.Info.ModTime().UTC().Format("Jan 02 15:04")

	table.AddRow(dir.Info.Mode(), data.NumberOfLinks, usr.Username, usr.Gid, l.minifySize(dir.Info.Size()), timeStr, name)

	return int(blocks), nil
}

func getFileInformationByHandle(handle syscall.Handle, data *syscall.ByHandleFileInformation) (err error) {
	modKernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetFileInformationByHandle := modKernel32.NewProc("GetFileInformationByHandle")
	r1, _, e1 := syscall.SyscallN(procGetFileInformationByHandle.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(data)), 0)
	if r1 == 0 {
		err = syscall.Errno(e1)
	}
	return
}
