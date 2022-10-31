package main

import macaron "gopkg.in/macaron.v1"

func distributionView(ctx *macaron.Context) {
	var as []*Address
	var asr []*AddressResponse

	db.Order("balance desc").Where("balance > 0").Find(&as)

	for _, a := range as {
		ar := &AddressResponse{
			Address: a.Address,
			Balance: a.Balance,
		}
		asr = append(asr, ar)
	}

	ctx.JSON(200, asr)
}
