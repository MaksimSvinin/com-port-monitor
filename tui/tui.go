package tui

import (
	"context"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
)

type Tui struct {
	ui tui.UI
}

type Message struct {
	Message string
	Style   string
}

type SidebarData struct {
	PortName string
	BaudRate uint
	DataBits uint
	StopBits uint
}

func NewUi(sidebarData SidebarData, inputCh chan string, writeCh chan Message) *Tui {
	sidebar := tui.NewVBox(
		tui.NewLabel("Connect info"),
		tui.NewLabel(fmt.Sprintf("port: %s", sidebarData.PortName)),
		tui.NewLabel(fmt.Sprintf("baud rate: %d", sidebarData.BaudRate)),
		tui.NewLabel(fmt.Sprintf("data bits: %d", sidebarData.DataBits)),
		tui.NewLabel(fmt.Sprintf("stop bits: %d", sidebarData.StopBits)),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	serialChat := tui.NewVBox(historyBox, inputBox)
	serialChat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		inputCh <- e.Text()
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, serialChat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	t := tui.NewTheme()
	t.SetStyle("label.normal", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorWhite})

	t.SetStyle("label.Red", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})
	t.SetStyle("label.Green", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorGreen})
	t.SetStyle("label.Blue", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorBlue})
	t.SetStyle("label.SerialRead", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorYellow})
	t.SetStyle("label.Magenta", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorMagenta})

	ui.SetTheme(t)

	go func() {
		for {
			m := <-writeCh
			ui.Update(func() {
				label := tui.NewLabel(m.Message)
				label.SetStyleName(m.Style)

				history.Append(tui.NewHBox(
					label,
					tui.NewSpacer(),
				))
			})
		}
	}()

	return &Tui{
		ui: ui,
	}
}

func (t *Tui) Run(ctx context.Context) {
	if err := t.ui.Run(); err != nil {
		log.Fatal(err)
	}
	ctx.Done()
}
