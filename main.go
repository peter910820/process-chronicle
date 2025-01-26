package main

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/shirou/gopsutil/v4/process"

	"processchronicle/internal"
)

var filter = []string{
	"C:\\Program Files\\",
	"C:\\Program Files (x86)\\",
	"C:\\Windows\\",
}
var (
	processList []internal.Process
	processName []string
)

func init() {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		name, _ := p.Name()
		path, _ := p.Exe()
		pid := p.Pid
		if name != "" && filterCheck(path) {
			processList = append(processList, internal.Process{
				Name: name,
				Path: path,
			})
			log.Printf("PID: %d, Name: %s\n", pid, path)
		}
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Choice Widgets")
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.SetFixedSize(true)
	myWindow.SetContent(initComponent())
	myWindow.ShowAndRun()
}

func initComponent() *fyne.Container {
	for _, process := range processList {
		processName = append(processName, process.Name)
	}
	combo := widget.NewSelect(processName, func(value string) {
		log.Println("Select set to", value)
	})
	combo.PlaceHolder = "選擇"
	content := widget.NewLabel("新增軟體: ")
	comboContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), combo)
	centeredContainer := container.NewHBox(content, comboContainer, layout.NewSpacer())

	return centeredContainer
}

func filterCheck(path string) bool {
	for _, prefix := range filter {
		if strings.HasPrefix(path, prefix) {
			return false
		}
	}
	return true
}
