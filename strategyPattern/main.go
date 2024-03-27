package main

import (
	"fmt"
	"math"
	"strconv"
)

const (
	//Rule
	RULE_0_01 = `rules-0.01`
	RULE_0_31 = `rules-0.31`
	DEFAULT   = `DEFAULT`

	//Product Type
	PRODUCTA = `PRODUCTA`
	PRODUCTB = `PRODUCTB`
	PRODUCTC = `PRODUCTC`
	PRODUCTD = `PRODUCTD`
	PRODUCTE = `PRODUCTE`
	PRODUCTF = `PRODUCTF`
)

type ChargeableWeightData struct {
	GrossWeight  float64
	VolumeWeight float64
	ProductType  string
}

type ruleConfigurable func(data ChargeableWeightData) (chargeAbleWeight float64, volumeWeight float64)

func getRule(rules string) ruleConfigurable {
	var Rules = map[string]ruleConfigurable{
		RULE_0_01: ruleZeroPointZeroOne,
		RULE_0_31: ruleZeroPointThirtyOne,
	}

	if _, ok := Rules[rules]; ok {
		return Rules[rules]
	}

	return ruleDefault
}

func ruleDefault(data ChargeableWeightData) (float64, float64) {
	volumeWeight := data.VolumeWeight
	chargeAbleWeight := data.GrossWeight
	if volumeWeight > chargeAbleWeight {
		chargeAbleWeight = volumeWeight
	}

	chargeAbleWeight = RoundChargeAbleWeight(chargeAbleWeight, data.ProductType)

	return chargeAbleWeight, volumeWeight
}

func ruleZeroPointThirtyOne(data ChargeableWeightData) (float64, float64) {
	volumeweight := data.VolumeWeight
	chargeAble := data.GrossWeight
	if volumeweight > chargeAble {
		chargeAble = volumeweight
	}

	chargeAble = TruncateFloat(chargeAble, 2)

	if chargeAble < 1.31 {
		chargeAble = 1
	} else {
		chargeAble = math.Round(chargeAble + 0.19)
	}
	return chargeAble, volumeweight
}

func ruleZeroPointZeroOne(data ChargeableWeightData) (float64, float64) {
	volumeweight := data.VolumeWeight
	chargeAble := data.GrossWeight
	if volumeweight > chargeAble {
		chargeAble = volumeweight
	}

	chargeAble = TruncateFloat(chargeAble, 2)

	chargeAble = math.Ceil(chargeAble)
	return chargeAble, volumeweight
}

func TruncateFloat(f float64, unit int) float64 {
	trimmedStringFloat := ""

	switch unit {
	case 2:
		trimmedStringFloat = fmt.Sprintf("%.2f", f)
	case 4:
		trimmedStringFloat = fmt.Sprintf("%.4f", f)
	default:
		trimmedStringFloat = fmt.Sprintf("%.2f", f)
	}

	trimmedFloat, _ := strconv.ParseFloat(trimmedStringFloat, 64)
	return trimmedFloat
}

func RoundChargeAbleWeight(chargeableWeight float64, productType string) (res float64) {
	chargeableWeight = TruncateFloat(chargeableWeight, 2)
	switch productType {
	case PRODUCTA:
		fallthrough
	case PRODUCTB:
		fallthrough
	case PRODUCTF:
		fallthrough
	case PRODUCTC:
		chargeableWeight = math.Ceil(chargeableWeight)
	default:
		if chargeableWeight < 1.31 {
			chargeableWeight = 1
		} else {
			chargeableWeight = math.Round(chargeableWeight + 0.19)
		}
	}

	return chargeableWeight
}
