package internal

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2/widget"

	"github.com/shirou/gopsutil/v4/process"

	"processchronicle/internal/register"
)

var (
	alias        string
	pid          int32 = 0
	counter            = 0
	RequestChan        = make(chan struct{})
	ResponseChan       = make(chan register.ChannelForJson)
)

func CheckProcess(guiComponent []GuiComponent) {
	data, err := register.ReadForJson()
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-RequestChan:
			channelData := register.ChannelForJson{
				Alias:   alias,
				Counter: counter,
			}
			ResponseChan <- channelData
		default:
			if pid == 0 {
				checkProcessOpen(data)
			} else {
				checkProcessClose(guiComponent)
			}
			time.Sleep(1 * time.Second)
		}

	}
}

func checkProcessOpen(data *register.DataForJson) {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		path, _ := p.Exe()
		for _, pro := range data.Register {
			if path == pro.Path {
				pid = p.Pid
				alias = pro.Alias
				log.Println("程式啟動中")
				break
			}
		}

	}
}

func checkProcessClose(guiComponent []GuiComponent) {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("error fetching processes: %v", err)
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
	data, err := register.ReadForJson()
	if err != nil {
		log.Fatal(err)
	}
	for i, p := range data.Register {
		if alias == p.Alias {
			originalTotalTime, err := strconv.Atoi(p.TotalTime)
			if err != nil {
				log.Fatal(err)
			}
			data.Register[i].TotalTime = fmt.Sprintf("%d", originalTotalTime+counter)

			err = register.WriteForJson(data)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}
	for _, c := range guiComponent {
		if c.Name == "TimerLabel" {
			if lbl, ok := c.Item.(*widget.Label); ok {
				lbl.SetText("")
			}
		}
	}
	log.Println("程式關閉")
}
