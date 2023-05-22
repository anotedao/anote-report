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
	StartHeight uint64
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

	loadNodes()

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

			balance, _, err := cl.Addresses.BalanceDetails(c, ao)
			if err != nil {
				log.Println(err)
				return
				// logTelegram(err.Error())
			}

			if a.Balance > balance.Effective &&
				!isNode(a.Address) &&
				a.Address != "3A9y1Zy78DDApbQWXKxonXxci6DvnJnnNZD" &&
				a.Address != "3ANmnLHt8mR9c36mdfQVpBtxUs8z1mMAHQW" &&
				a.Address != "3ANzidsKXn9a1s9FEbWA19hnMgV9zZ2RB9a" {
				logTelegram("Suspicious activity: " + a.Address)
			}

			if !isNode(a.Address) {
				a.Balance = balance.Effective
				a.New = balance.Effective - balance.Generating
			} else {
				a.Balance = 0
			}
			db.Save(a)
		}
	}
}

func (m *Monitor) start() {
	for {
		h := getHeight()

		for i := m.StartHeight; i <= h; i++ {
			m.processBlock(i)
		}

		m.StartHeight = h

		log.Println("Done loading blocks.")

		m.loadBalances()

		log.Println("Done loading balances.")

		ks := &KeyValue{Key: "lastHeight"}
		db.FirstOrCreate(ks, ks)
		ks.ValueInt = h
		db.Save(ks)

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() {
	m := &Monitor{StartHeight: 1}

	ks := &KeyValue{Key: "lastHeight"}
	db.FirstOrCreate(ks, ks)

	if ks.ValueInt == 0 {
		ks.ValueInt = 1
		db.Save(ks)
	}

	m.StartHeight = ks.ValueInt

	go m.start()
}
