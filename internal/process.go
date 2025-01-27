package internal

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2/widget"

	"github.com/shirou/gopsutil/v4/process"

	"processchronicle/internal/register"
)

var pid int32 = 0
var counter = 0

func CheckProcess(guiComponent []GuiComponent) {
	data, err := register.ReadForJson()
	if err != nil {
		log.Fatal(err)
	}
	for {
		if pid == 0 {
			checkProcessOpen(data)
		} else {
			checkProcessClose(guiComponent)
		}
		time.Sleep(1 * time.Second)
	}
}

func checkProcessOpen(data *register.DataForJson) {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		path, _ := p.Exe()
		if path == data.Register[0].Path {
			pid = p.Pid
			log.Println("程式啟動中")
			break
		}
	}
}
func checkProcessClose(guiComponent []GuiComponent) {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		if p.Pid == pid {
			for _, c := range guiComponent {
				if c.Name == "TimerLabel" {
					if lbl, ok := c.Item.(*widget.Label); ok {
						counter += 1
						lbl.SetText(fmt.Sprintf("%d", counter))
					}
				}
			}
			return
		}
	}
	pid = 0
	log.Println("程式關閉")
}
