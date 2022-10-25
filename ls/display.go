package ls

import (
	"fmt"

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
		if dir.Info.IsDir() {
			name = c.Cyan(name)
		}

		if !l.L {
			tb.AddCell(name)
			continue
		}

		blkSize, err := l._display(tb, name, &dir)
		if err != nil {
			return err
		}

		totalBlkSize += blkSize
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
