package main

import (
	"log"
	"time"
)

type Monitor struct {
}

func (m *Monitor) start() {
	for {
		log.Println("Tick.")

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() {
	m := &Monitor{}
	go m.start()
}
