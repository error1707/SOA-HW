package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	ChatView     = "chat"
	ChatInput    = "chat_input"
	UserListView = "user_list"
	ControlInput = "control_panel"
	GameLog      = "game_log"
)

var sizeX, sizeY int

func ScaleX(scale float64) int {
	return int(float64(sizeX) * scale)
}

func ScaleY(scale float64) int {
	return int(float64(sizeY) * scale)
}

func SetUpGUI() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	g.SetManagerFunc(Layout)
	if err := InitKeyBindings(g); err != nil {
		return nil, err
	}
	return g, nil
}

func Layout(g *gocui.Gui) error {
	sizeX, sizeY = g.Size()
	v, err := g.SetView(ChatView, 0, 0, ScaleX(.34), ScaleY(.74)-2)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Title = "Chat"
	v.Wrap = true

	v, err = g.SetView(ChatInput, 0, ScaleY(.74)-2, ScaleX(.34), ScaleY(.74))
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Title = "Chat Input"
	v.Editable = true

	v, err = g.SetView(UserListView, 0, ScaleY(.75), ScaleX(.34), ScaleY(.99))
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Title = "User List"
	v.Wrap = true

	v, err = g.SetView(ControlInput, ScaleX(.35), ScaleY(.99)-2, ScaleX(.99), ScaleY(.99))
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Title = "Control Input"
	v.Editable = true

	v, err = g.SetView(GameLog, ScaleX(.35), 0, ScaleX(.99), ScaleY(.97)-2)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	v.Title = "Game Log"
	v.Wrap = true
	return nil
}

func InitKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	g.Mouse = true

	for _, viewName := range []string{ChatView, ChatInput, UserListView, ControlInput, GameLog} {
		if err := g.SetKeybinding(viewName, gocui.MouseLeft, gocui.ModNone, setActiveView); err != nil {
			return err
		}
	}

	for _, viewName := range []string{ChatInput, ControlInput} {
		if err := g.SetKeybinding(viewName, gocui.MouseRelease, gocui.ModNone, setInputCursor); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding(ChatInput, gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		return err
	}

	if err := g.SetKeybinding(ControlInput, gocui.KeyEnter, gocui.ModNone, processCommand); err != nil {
		return err
	}

	// Enable scrolling
	//for _, viewName := range []string{ChatView, UserListView, GameLog} {
	//	if err := g.SetKeybinding(viewName, 'a', gocui.ModNone,
	//		func(g *gocui.Gui, v *gocui.View) error {
	//			v.Autoscroll = true
	//			return nil
	//		}); err != nil {
	//		return err
	//	}
	//
	//	if err := g.SetKeybinding(viewName, gocui.MouseWheelDown, gocui.ModNone,
	//		func(g *gocui.Gui, v *gocui.View) error {
	//			return scrollView(v, -1)
	//		}); err != nil {
	//		return err
	//	}
	//	if err := g.SetKeybinding(viewName, gocui.MouseWheelUp, gocui.ModNone,
	//		func(g *gocui.Gui, v *gocui.View) error {
	//			return scrollView(v, 1)
	//		}); err != nil {
	//		return err
	//	}
	//}
	return nil
}

func setActiveView(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(v.Name())
	if v.Name() == ChatInput || v.Name() == ControlInput {
		g.Cursor = true
	} else {
		g.Cursor = false
	}
	return err
}

func setInputCursor(_ *gocui.Gui, v *gocui.View) error {
	sizeX, _ := v.Size()
	_, y := v.Cursor()
	x := len(v.Buffer())
	if x > sizeX {
		x = sizeX
	}
	if x != 0 {
		x--
	}
	err := v.SetCursor(x, y)
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(_ *gocui.Gui, v *gocui.View) error {
	if len(v.Buffer()) == 0 {
		return nil
	}
	Chat <- v.Buffer()
	x, y := v.Cursor()
	x -= len(v.Buffer()) - 1
	v.Clear()
	err := v.SetCursor(x, y)
	if err != nil {
		return err
	}
	return nil
}

func processCommand(_ *gocui.Gui, v *gocui.View) error {
	if len(v.Buffer()) == 0 {
		return nil
	}
	ProcessCommand(v.Buffer())
	x, y := v.Cursor()
	x -= len(v.Buffer()) - 1
	v.Clear()
	err := v.SetCursor(x, y)
	if err != nil {
		return err
	}
	return nil
}

func WriteToLog(text string) {
	GUI.Update(func(gui *gocui.Gui) error {
		v, err := GUI.View(GameLog)
		if err != nil {
			return err
		}
		_, err = v.Write([]byte(text))
		if err != nil {
			return err
		}
		return nil
	})
}

func WriteToChat(text string, player Player) {
	GUI.Update(func(gui *gocui.Gui) error {
		v, err := gui.View(ChatView)
		if err != nil {
			return err
		}
		switch player.status {
		case ALIVE:
			v.Write([]byte(fmt.Sprintf("%s >>> %s", Green(player.username), text)))
		case DEAD:
			v.Write([]byte(fmt.Sprintf("%s >>> %s", Red(player.username), text)))
		}
		return nil
	})
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var (
	Red   = Color("\033[1;31m%s\033[0m")
	Green = Color("\033[1;32m%s\033[0m")
)

func UpdateUserList() {
	GUI.Update(func(gui *gocui.Gui) error {
		v, err := gui.View(UserListView)
		if err != nil {
			return err
		}
		v.Clear()
		for _, player := range Players {
			switch player.status {
			case ALIVE:
				v.Write([]byte(Green(player.username) + ", "))
			case DEAD:
				v.Write([]byte(Red(player.username) + ", "))
			}
		}
		return nil
	})
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}
