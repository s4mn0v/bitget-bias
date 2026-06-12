# Bitget Bias

Bitget Bias is a professional command-line market analysis tool built in Go. It aggregates real-time data from Bitget's Futures markets to calculate a directional bias score by weighing institutional aggression against retail sentiment.

## Highlights

*   Institutional Aggression Tracking: Monitors Taker Buy/Sell volume to identify market displacement.
*   Retail Sentiment Analysis: Analyzes Long/Short ratios to identify potential retail traps and liquidity sweeps.
*   Funding Rate Monitoring: Detects overheated positioning that often precedes market reversals.
*   Automated Scoring: Provides a consolidated Bullish, Bearish, or Neutral bias based on multi-factor weighted analysis.
*   AMD Framework: Designed to assist traders using the Accumulation-Manipulation-Distribution methodology.

## Overview

The core philosophy of Bitget Bias is rooted in the interaction between institutional activity and retail positioning. Most retail traders are liquidated at key levels; this tool identifies when retail is "over-extended" in one direction while institutional participants are actively hitting the tape.

By comparing the Taker Buy/Sell Volume (Institutional Aggression) with the Futures Long/Short Ratio (Retail Sentiment), the tool identifies high-probability zones for market movement.

### Key Questions Answered
1.  Is the current trend supported by aggressive buying or selling?
2.  Is retail heavily one-sided, suggesting a potential sweep of their liquidity?
3.  Is the funding rate suggesting an exhausted move?

## Usage

The tool runs a continuous monitoring loop, polling the Bitget API every 30 seconds for the BTCUSDT pair. It provides a real-time visual balance of the Long/Short ratio.

### Execution

```bash
./bitget-bias
```

### Understanding the Output

```text
--- BTCUSDT 10:15:30 ---
Aggro: Taker BUYING
Ratio: 0.72 | [S] -------X-----|----- [L] | Dev: -0.28
Trap: Retail Over-Short (Wait Sweep Up)
BIAS SCORE: 2 -> BIAS: BULLISH (Target BSL/Distribution Up)
```

*   Aggro: Shows whether market orders (Takers) are predominantly buying or selling.
*   Ratio: A visual scale representing retail positioning.
*   Trap: Specific warning when retail is over-leveraged in a specific direction.
*   Bias Score: A numerical sum of all indicators (-3 to +3).

## Technical Reference

The application leverages specific Bitget API endpoints to build its bias model:

| Category | Endpoint | Indicator |
| :--- | :--- | :--- |
| Aggression | /api/v2/mix/market/taker-buy-sell | Institutional Aggression (Market Orders) |
| Sentiment | /api/v2/mix/market/long-short | Retail Trap identification (Crowd Long/Short) |
| Exhaustion | /api/v2/mix/market/current-fund-rate | Detects over-leveraged costs for holding positions |

## Installation

### Prerequisites
*   Go 1.26.2 or higher.

### From Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/s4mn0v/bitget-bias.git
    ```
2.  Navigate to the directory:
    ```bash
    cd bitget-bias
    ```
3.  Run it | Build the binary:
    ```bash
    go run .
    ```
    ```
    
    ```bash
    go build -o bitget-bias
    ```

## Development and Contributions

The tool is designed with an extensible architecture. Future updates intend to incorporate:
*   Spot Whale Net Flow: To compare spot accumulation with futures positioning.
*   WebSocket Integration: For sub-second response times on volume spikes.
*   Orderbook Analysis: Identifying large limit order walls (Sell-Side and Buy-Side Liquidity).

If you wish to contribute, please submit a pull request or open an issue to discuss proposed changes.

## Author

[S4M-N0V](https://github.com/s4mn0v/)

## Disclaimer

This software is for informational purposes only. Trading cryptocurrencies involves significant risk. The bias score is a mathematical representation of specific data points and does not guarantee market direction. Never trade with money you cannot afford to lose.

