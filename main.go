package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const BaseURL = "https://api.bitget.com"

func fetch(path string, target interface{}) error {
	resp, err := http.Get(BaseURL + path)
	if err != nil { return err }
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res BitgetRes
	json.Unmarshal(body, &res)
	return json.Unmarshal(res.Data, target)
}

func drawBalance(ratio float64) {
	const size = 20
	// Rango 0.5 a 1.5
	norm := (ratio - 0.5) / (1.5 - 0.5)
	pos := int(norm * size)

	if pos < 0 { pos = 0 }
	if pos > size { pos = size }

	display := make([]rune, size+1)
	for i := range display {
		display[i] = '-'
	}
	display[size/2] = '|'
	display[pos] = 'X'

	fmt.Printf("Ratio: %.2f | [S] %s [L] | Dev: %+.2f\n", ratio, string(display), ratio-1.0)
}

func getBias(symbol string) {
	fmt.Printf("\n--- %s %s ---\n", symbol, time.Now().Format("15:04:05"))
	
	score := 0

	// 1. Whale Flow (Institutional)
	var fr []FundingRate
	if fetch("/api/v2/mix/market/current-fund-rate?symbol="+symbol+"&productType=usdt-futures", &fr) == nil {
		rate, _ := strconv.ParseFloat(fr[0].FundingRate, 64)
		if rate > 0.01 { score--; fmt.Println("Funding: Overheated Longs") }
		if rate < -0.01 { score++; fmt.Println("Funding: Overheated Shorts") }
	}

	// 2. Open Interest (Money Flow)
	var tv []TakerVolume
	if fetch("/api/v2/mix/market/taker-buy-sell?symbol="+symbol, &tv) == nil && len(tv) > 0 {
		b, _ := strconv.ParseFloat(tv[0].BuyVolume, 64)
		s, _ := strconv.ParseFloat(tv[0].SellVolume, 64)
		if b > s { score++; fmt.Println("Aggro: Taker BUYING") } else { score--; fmt.Println("Aggro: Taker SELLING") }
	}

	// 3. Retail Trap (L/S Ratio)
	var ls []LongShortRatio
	if fetch("/api/v2/mix/market/long-short?symbol="+symbol, &ls) == nil && len(ls) > 0 {
		ratio, _ := strconv.ParseFloat(ls[0].LongShortRatio, 64)
		
		drawBalance(ratio)
		if ratio > 1.2 { score--; fmt.Println("Trap: Retail Over-Long (Wait Sweep Down)") }
		if ratio < 0.8 { score++; fmt.Println("Trap: Retail Over-Short (Wait Sweep Up)") }
	}

	// Result
	fmt.Printf("BIAS SCORE: %d -> ", score)
	switch {
	case score > 0: fmt.Println("BIAS: BULLISH (Target BSL/Distribution Up)")
	case score < 0: fmt.Println("BIAS: BEARISH (Target SSL/Distribution Down)")
	default: fmt.Println("BIAS: NEUTRAL (Accumulation Phase)")
	}
}

func main() {
	for {
		getBias("BTCUSDT")
		time.Sleep(30 * time.Second)
	}
}
