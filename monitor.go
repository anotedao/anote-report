package main

import (
	"fmt"
	"log"
	"time"
)

type Monitor struct {
}

func (m *Monitor) processBlock(n uint64) {
	as := getBlockAddresses(n)

	for _, a := range as {
		adb := &Address{}
		db.FirstOrCreate(adb, &Address{Address: a})
	}
}

func (m *Monitor) start() {
	for {
		h := getHeight()

		for i := uint64(1); i <= h; i++ {
			m.processBlock(i)
			log.Println(fmt.Sprintf("Done block: %d", i))
		}

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() {
	m := &Monitor{}
	go m.start()
}
