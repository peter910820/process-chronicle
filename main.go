package main

import (
	"log"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

var filter = []string{
	"C:\\Program Files\\",
	"C:\\Program Files (x86)\\",
	"C:\\Windows\\",
}

func main() {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error fetching processes: %v", err)
	}
	for _, p := range processes {
		name, _ := p.Name()
		path, _ := p.Exe()
		pid := p.Pid
		if name != "" && filterCheck(path) {
			log.Printf("PID: %d, Name: %s\n", pid, path)
		}
	}
}

func filterCheck(path string) bool {
	for _, prefix := range filter {
		if strings.HasPrefix(path, prefix) {
			return false
		}
	}
	return true
}
