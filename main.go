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

func getBias(symbol string) {
	fmt.Printf("\n--- %s %s ---\n", symbol, time.Now().Format("15:04:05"))
	
	score := 0

	// 1. Whale Flow (Institutional)
	var wf []WhaleFlow
	if fetch("/api/v2/spot/market/whale-net-flow?symbol="+symbol, &wf) == nil && len(wf) > 0 {
		vol, _ := strconv.ParseFloat(wf[0].Volume, 64)
		if vol > 0 { score++; fmt.Println("Whale: BUYing") } else { score--; fmt.Println("Whale: SELLing") }
	}

	// 2. Open Interest (Money Flow)
	var oi OpenInterest
	if fetch("/api/v2/mix/market/open-interest?symbol="+symbol+"&productType=usdt-futures", &oi) == nil {
		fmt.Printf("OI Size: %s\n", oi.OpenInterestList[0].Size)
	}

	// 3. Retail Trap (L/S Ratio)
	var ls []LongShortRatio
	if fetch("/api/v2/mix/market/long-short?symbol="+symbol, &ls) == nil && len(ls) > 0 {
		ratio, _ := strconv.ParseFloat(ls[0].LongShortRatio, 64)
		fmt.Printf("L/S Ratio: %.2f\n", ratio)
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
