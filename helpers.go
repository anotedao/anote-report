package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
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
	if !strings.HasPrefix(str, "3A") {
		return s
	}

	if !contains(s, str) {
		s = append(s, str)
	}
	return s
}

type AddressResponse struct {
	Address      string  `json:"address"`
	Balance      uint64  `json:"balance"`
	BalanceFloat float64 `json:"balance_float"`
}

func getCallerInfo() (info string) {

	// pc, file, lineNo, ok := runtime.Caller(2)
	_, file, lineNo, ok := runtime.Caller(2)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	// funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // The Base function returns the last element of the path
	return fmt.Sprintf("%s:%d: ", fileName, lineNo)
}

func logTelegram(message string) {
	message = "anote-report:" + getCallerInfo() + url.PathEscape(url.QueryEscape(message))

	_, err := http.Get(fmt.Sprintf("http://localhost:5002/log/%s", message))
	if err != nil {
		log.Println(err)
	}
}

func isNode(address string) bool {
	for _, node := range nodes {
		if node == address {
			return true
		}
	}

	return false
}

func loadNodes() {
	nodes = []string{}

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}, ApiKey: " "})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	addr := proto.MustAddressFromString(NodesListAddress)

	de, _, err := cl.Addresses.AddressesData(ctx, addr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	for _, node := range de {
		nodes = append(nodes, node.ToProtobuf().GetStringValue())
	}
}
