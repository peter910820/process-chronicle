package internal

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

var pid int32 = 0

func CheckProcess() {
	for {
		if pid == 0 {
			checkProcessOpen()
		} else {
			checkProcessClose()
		}
		time.Sleep(1 * time.Second)
	}
}

func checkProcessOpen() {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		path, _ := p.Exe()
		if path == "" {
			pid = p.Pid
			log.Println("程式啟動中")
			break
		}
	}
}
func checkProcessClose() {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		if p.Pid == pid {
			return
		}
	}
	pid = 0
	log.Println("程式關閉")
}
