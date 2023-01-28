package main

import (
	"fmt"
	"log"
	"math/big"
)

func ComputePrice(timePeriod uint64, fileSize uint64) (*big.Int, bool) {
	days := float64(timePeriod) / (60 * 60 * 24)

	mbs := float64(fileSize) / (1024 * 1024)

	k := days * mbs
	price := new(big.Float)
	price, valid := price.SetString(NC.FeeWeiPerMBPerDay)
	if !valid {
		log.Println("Failed to parse FeeWeiPerMBPerDay", NC.FeeWeiPerMBPerDay)
		return nil, false
	}

	price = price.Mul(price, big.NewFloat((k)))

	priceInt := new(big.Int)

	intVal := fmt.Sprintf("%.0f", price)

	priceInt.SetString(intVal, 10)

	baseFee := new(big.Int)
	baseFee, val := baseFee.SetString(NC.FeeBase, 10)
	if !val {
		return nil, false
	}
	priceInt = priceInt.Add(priceInt, baseFee)
	return priceInt, true
}
