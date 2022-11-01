package main

import (
	"fmt"
	"math"

	macaron "gopkg.in/macaron.v1"
)

func distributionView(ctx *macaron.Context) {
	var as []*Address
	var asr []*AddressResponse

	db.Order("balance desc").Where("balance > 0").Find(&as)

	for _, a := range as {
		bf := float64(a.Balance) / float64(MULTI8)
		bf = math.Floor(bf*MULTI8) / MULTI8

		ar := &AddressResponse{
			Address:      a.Address,
			Balance:      a.Balance,
			BalanceFloat: bf,
		}
		asr = append(asr, ar)
	}

	ctx.JSON(200, asr)
}

func distView(ctx *macaron.Context) string {
	var as []*Address
	response := "<html><pre>"

	db.Order("balance desc").Where("balance > 0").Find(&as)

	for i, a := range as {
		bf := float64(a.Balance) / float64(MULTI8)

		response += fmt.Sprintf("%d.&nbsp;&nbsp;&nbsp;%s&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;%.8f\n", i+1, a.Address, bf)
	}

	response += "</pre></html>"

	return response
}
