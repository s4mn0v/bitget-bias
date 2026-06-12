package main

import "encoding/json"

type BitgetRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type WhaleFlow struct {
	Volume string `json:"volume"`
	Date   string `json:"date"`
}

type OpenInterest struct {
	OpenInterestList []struct {
		Symbol string `json:"symbol"`
		Size   string `json:"size"`
	} `json:"openInterestList"`
}

type LongShortRatio struct {
	Symbol          string `json:"symbol"`
	LongShortRatio  string `json:"longShortRatio"`
	Ts              string `json:"ts"`
}

type FundingRate struct {
	Symbol      string `json:"symbol"`
	FundingRate string `json:"fundingRate"`
}

type TakerVolume struct {
	BuyVolume  string `json:"buyVolume"`
	SellVolume string `json:"sellVolume"`
}
