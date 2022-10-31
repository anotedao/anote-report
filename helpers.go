package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
)

func getHeight() uint64 {
	height := uint64(0)

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		// logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bh, _, err := cl.Blocks.Height(ctx)

	height = bh.Height

	return height
}

func getBlockAddresses(n uint64) []string {
	var as []string

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		// logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	b, _, err := cl.Blocks.At(ctx, n)

	as = myappend(as, b.Generator.String())

	for _, t := range b.Transactions {
		at := AnoteTransaction{}
		trb, err := json.Marshal(t)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal(trb, &at)

		s, err := t.GetSender(55)
		if err != nil {
			log.Println(err)
		}

		as = myappend(as, s.String())

		if len(at.Recipient) > 0 {
			as = myappend(as, at.Recipient)
		}
	}

	return as
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

type AnoteTransaction struct {
	Type            int         `json:"type"`
	Version         int         `json:"version"`
	ID              string      `json:"id"`
	Proofs          []string    `json:"proofs"`
	SenderPublicKey string      `json:"senderPublicKey"`
	AssetID         interface{} `json:"assetId"`
	FeeAssetID      interface{} `json:"feeAssetId"`
	Timestamp       int64       `json:"timestamp"`
	Amount          int         `json:"amount"`
	Fee             int         `json:"fee"`
	Recipient       string      `json:"recipient"`
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func myappend(s []string, str string) []string {
	if !contains(s, str) {
		s = append(s, str)
	}
	return s
}
