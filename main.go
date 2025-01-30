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
	"processchronicle/internal/register"
)

var (
	processList  []internal.ProcessList
	guiComponent []internal.GuiComponent
	processName  []string
)

func init() {
	// get filter data
	data, err := register.ReadForJson()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Read record file complete")
	// traverse all processes
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		name, _ := p.Name()
		path, _ := p.Exe()
		pid := p.Pid
		if name != "" && path != "" && filterCheck(data, path) {
			processList = append(processList, internal.ProcessList{
				Pid:  pid,
				Name: name,
				Path: path,
			})
			log.Printf("PID: %d, Name: %s\n", pid, path)
		}
	}
}

func main() {
	guiComponent, box := initComponent()
	go internal.CheckProcess(guiComponent)
	myApp := app.New()
	myWindow := myApp.NewWindow("Choice Widgets")
	myWindow.SetCloseIntercept(func() {
		internal.RequestChan <- struct{}{}
		data := <-internal.ResponseChan
		record, err := register.ReadForJson()
		if err != nil {
			log.Fatal(err)
		}
		for i, r := range record.Register {
			if r.Alias == data.Alias {
				originalTotalTime, err := strconv.Atoi(r.TotalTime)
				if err != nil {
					log.Fatal(err)
				}
				record.Register[i].TotalTime = fmt.Sprintf("%d", originalTotalTime+data.Counter)

				err = register.WriteForJson(record)
				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}
		log.Print(data)
		myApp.Quit()
	})
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)
	myWindow.SetContent(box)
	myWindow.ShowAndRun()
}

func initComponent() ([]internal.GuiComponent, *fyne.Container) {
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
	registerComponent("PathLabel", pathLabel)

	registerButton := widget.NewButton("註冊程式", func() {
		err := register.RegisterForJson(pathLabel.Text)
		if err != nil {
			log.Fatal(err)
		}
	})
	registerComponent("RegisterComponent", registerButton)

	timerLabel := widget.NewLabel("")
	timerLabel.TextStyle = fyne.TextStyle{
		Bold:      true,
		Monospace: true,
	}
	registerComponent("TimerLabel", timerLabel)

	box := container.NewVBox(
		container.NewHBox(content, comboContainer, layout.NewSpacer()),
		container.NewHBox(pathLabel),
		container.NewHBox(registerButton),
		container.NewHBox(layout.NewSpacer(), timerLabel, layout.NewSpacer()),
	)
	return guiComponent, box
}

func filterCheck(data *register.DataForJson, path string) bool {
	for _, p := range data.Filter {
		if strings.HasPrefix(path, p) {
			return false
		}
	}
	return true
}

func registerComponent(name string, component interface{}) {
	guiComponent = append(guiComponent, internal.GuiComponent{
		Name: name,
		Item: component,
	})
}
