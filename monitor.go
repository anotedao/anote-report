package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
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

func (m *Monitor) loadBalances() {
	var as []*Address
	db.Find(&as)

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		// logTelegram(err.Error())
	}

	for _, a := range as {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if strings.HasPrefix(a.Address, "3A") {
			ao := proto.MustAddressFromString(a.Address)

			balance, _, err := cl.Addresses.Balance(c, ao)
			if err != nil {
				log.Println(err)
				// logTelegram(err.Error())
			}

			if a.Balance > balance.Balance &&
				!isNode(a.Address) &&
				a.Address != "3A9y1Zy78DDApbQWXKxonXxci6DvnJnnNZD" &&
				a.Address != "3ANzidsKXn9a1s9FEbWA19hnMgV9zZ2RB9a" {
				logTelegram("Suspicious activity: " + a.Address)
			}

			a.Balance = balance.Balance
			db.Save(a)
		}
	}
}

func (m *Monitor) start() {
	for {
		h := getHeight()

		for i := uint64(1); i <= h; i++ {
			m.processBlock(i)
		}

		log.Println("Done loading blocks.")

		m.loadBalances()

		log.Println("Done loading balances.")

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() {
	m := &Monitor{}
	go m.start()
}
