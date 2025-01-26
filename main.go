package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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
	processList  []internal.ProcessList
	guiComponent []internal.GuiComponent
	processName  []string
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
			processList = append(processList, internal.ProcessList{
				Pid:  pid,
				Name: name,
				Path: path,
			})
			log.Printf("PID: %d, Name: %s\n", pid, path)
		}
	}
	go internal.CheckProcess()
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Choice Widgets")
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)
	myWindow.SetContent(initComponent())
	myWindow.ShowAndRun()
}

func initComponent() *fyne.Container {
	for _, process := range processList {
		processName = append(processName, fmt.Sprintf("%s <%d>", process.Name, process.Pid))
	}
	combo := widget.NewSelect(processName, func(value string) {
		re := regexp.MustCompile(`<([^>]+)>`)
		match := re.FindStringSubmatch(value)
		num, err := strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range processList {
			if p.Pid == int32(num) {
				for _, c := range guiComponent {
					if c.Name == "PathLabel" {
						if lbl, ok := c.Item.(*widget.Label); ok {
							lbl.SetText(p.Path)
						}
					}
				}
			}
		}
		log.Println("Select set to", value)
	})
	combo.PlaceHolder = "選擇"
	content := widget.NewLabel("新增軟體: ")
	comboContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), combo)

	pathLabel := widget.NewLabel("")
	pathLabel.TextStyle = fyne.TextStyle{
		Bold:   true,
		Italic: true,
	}
	guiComponent = append(guiComponent, internal.GuiComponent{
		Name: "PathLabel",
		Item: pathLabel,
	})

	box := container.NewVBox(
		container.NewHBox(content, comboContainer, layout.NewSpacer()),
		container.NewHBox(pathLabel),
	)

	return box
}

func filterCheck(path string) bool {
	for _, p := range filter {
		if strings.HasPrefix(path, p) {
			return false
		}
	}
	return true
}
