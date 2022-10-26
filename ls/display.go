package ls

import (
	"fmt"
	"os"

	"github.com/pee2pee/lse/ls/color"
	"github.com/profclems/glab/pkg/tableprinter"
)

func (l *LS) display(dirs []Dir) (err error) {
	totalBlkSize := 0
	c := color.Color()

	tb := tableprinter.NewTablePrinter()
	tb.Wrap = true
	tb.SetTerminalWidth(color.TerminalWidth(l.StdOut))

	for i := range dirs {
		dir := dirs[i]
		name := dir.Info.Name()

		if l.Q {
			name = fmt.Sprintf("%q", name)
		}

		// set file color depending on type
		switch {
		case dir.Info.IsDir(): //dir
			name = c.BoldCyan(name)
		case dir.Info.Mode()&os.ModeSymlink != 0: //symlink
			name = c.Magenta(name)
		case dir.Info.Mode()&0100 != 0: //executable by user
			name = c.Red(name)
		}

		// display format
		switch {
		case l.One:
			tb.AddRow(name)
		case l.L:
			blkSize, err := l._display(tb, name, &dir)
			if err != nil {
				return err
			}
			totalBlkSize += blkSize
		case !color.IsTerminal(l.StdOut):
			tb.AddRow(name)
		default:
			tb.AddCell(name)
		}
	}

	if !l.L {
		tb.EndRow()
	}

	if l.L {
		_, err = fmt.Fprintln(l.StdOut, "total", totalBlkSize)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(l.StdOut, tb.String())
	return err
}
