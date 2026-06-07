# Bitget Bias

Bitget Bias is a command-line market analysis tool written in Go. It aggregates real-time data from Bitget's Spot and Futures markets to calculate a directional bias score based on institutional flow and retail sentiment.

## Highlights

*   Institutional Whale Flow: Identifies whether "Smart Money" is in an accumulation or distribution phase.
*   Retail Sentiment Analysis: Monitors Long/Short ratios to identify potential retail traps and upcoming liquidity sweeps.
*   Open Interest Tracking: Measures money flow into futures to confirm the strength of market displacement.
*   Automated Scoring: Provides a consolidated Bullish, Bearish, or Neutral bias based on multi-factor weighted analysis.
*   Extensible Architecture: Built with types ready to incorporate Orderbook, Fund Flow, and Fair Value Gap (FVG) data.

## Overview

The core philosophy of Bitget Bias is rooted in the interaction between institutional activity and retail positioning. By comparing Spot Whale Net Flow (Institutional Bias) with the Futures Long/Short Ratio (Retail Sentiment), the tool identifies high-probability zones for market movement.

The logic follows the Accumulation-Manipulation-Distribution (AMD) framework:
1.  Accumulation: Identified through sideways candle structures and Whale Net Flow increases.
2.  Manipulation: Detected when retail sentiment becomes over-extended (High L/S ratio), signaling a potential sweep of Buy-Side Liquidity (BSL) or Sell-Side Liquidity (SSL).
3.  Distribution: Confirmed when Whale Flow and Open Interest (OI) align with the directional price movement.

## Methodology and Data Sources

The application leverages specific Bitget API endpoints to build its bias model:

### Spot Data
*   Whale Flow: Uses `/api/v2/spot/market/whale-net-flow` to determine institutional buying vs. selling.
*   Fund Flow: (Internal Logic) Distinguishes between fish, dolphin, and whale volume to confirm the current phase of the market cycle.

### Futures Data
*   Open Interest (OI): Uses `/api/v2/mix/market/open-interest`. Rising OI during price increases confirms displacement, while falling OI suggests a trend conclusion.
*   Long/Short Ratio: Uses `/api/v2/mix/market/long-short`. Extreme ratios (e.g., > 1.2 or < 0.8) indicate a crowded retail trade likely to be swept by larger market participants.

## Usage

The tool runs a continuous monitoring loop, polling the Bitget API every 30 seconds for the BTCUSDT pair.

### Execution

```bash
./bitget-bias
```

### Example Output

```text
--- BTCUSDT 10:15:30 ---
Whale: BUYing
OI Size: 15420.25
L/S Ratio: 0.72
Trap: Retail Over-Short (Wait Sweep Up)
BIAS SCORE: 2 -> BIAS: BULLISH (Target BSL/Distribution Up)
```

## Technical Reference

The following Bitget data points inform the current and future logic of this tool:

| Category | Endpoint | Indicator |
| :--- | :--- | :--- |
| Institutional | Whale Net Flow | Bias identification (Buy > Sell = Bullish) |
| Sentiment | Long/Short Ratio | Retail Trap identification (Crowd Long = Potential Sweep Down) |
| Momentum | Open Interest | Displacement confirmation (OI Up + Price Up = Strong Trend) |
| Liquidity | Orderbook | Large limit order identification (BSL/SSL zones) |
| Structure | Market Candles | Range detection and Fair Value Gaps (FVG) |

## Installation

### From Source

1.  Ensure you have Go 1.26.2 or higher installed.
2.  Clone the repository:
    ```bash
    git clone https://github.com/youruser/bitget-bias.git
    ```
3.  Navigate to the directory:
    ```bash
    cd bitget-bias
    ```
4.  Build the binary:
    ```bash
    go build -o bitget-bias
    ```

## Development and Contributions

Future updates intend to incorporate WebSocket (WS) channels for real-time execution data:
*   Fill Channel: To detect the difference between wick and body closures for sweep confirmation.
*   Ticker Channel: To monitor spread increases that signal high-volatility manipulation risks.

If you wish to contribute, please submit a pull request or open an issue to discuss proposed changes.

## Disclaimer

This software is for informational purposes only. Trading cryptocurrencies involves significant risk. The bias score is a mathematical representation of specific data points and does not guarantee market direction.
