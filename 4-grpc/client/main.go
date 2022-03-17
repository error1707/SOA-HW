package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"os"
)

var (
	GUI *gocui.Gui
)

func main() {
	var err error
	GUI, err = SetUpGUI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on GUI set up: %v\n", err)
	}
	defer GUI.Close()

	if err = GUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
	_ = disconnect([]string{})
}
