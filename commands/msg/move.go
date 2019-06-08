package msg

import (
	"errors"
	"time"

	"git.sr.ht/~sircmpwn/getopt"
	"github.com/gdamore/tcell"

	"git.sr.ht/~sircmpwn/aerc/widgets"
	"git.sr.ht/~sircmpwn/aerc/worker/types"
)

func init() {
	register("mv", Move)
	register("move", Move)
}

func Move(aerc *widgets.Aerc, args []string) error {
	opts, optind, err := getopt.Getopts(args[1:], "p")
	if err != nil {
		return err
	}
	if optind != len(args)-2 {
		return errors.New("Usage: mv [-p] <folder>")
	}
	var (
		createParents bool
	)
	for _, opt := range opts {
		switch opt.Option {
		case 'p':
			createParents = true
		}
	}

	widget := aerc.SelectedTab().(widgets.ProvidesMessage)
	acct := widget.SelectedAccount()
	if acct == nil {
		return errors.New("No account selected")
	}
	msg := widget.SelectedMessage()
	store := widget.Store()
	_, isMsgView := widget.(*widgets.MessageViewer)
	if isMsgView {
		aerc.RemoveTab(widget)
	}
	acct.Messages().Next()
	store.Move([]uint32{msg.Uid}, args[optind+1], createParents, func(
		msg types.WorkerMessage) {

		switch msg := msg.(type) {
		case *types.Done:
			aerc.PushStatus("Messages moved.", 10*time.Second)
		case *types.Error:
			aerc.PushStatus(" "+msg.Error.Error(), 10*time.Second).
				Color(tcell.ColorDefault, tcell.ColorRed)
		}
	})
	return nil
}
