package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"encoding/json"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	refreshEvery = 30
	BaseURL = "https://api.bitget.com"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

type tickMsg time.Time

type model struct {
	table    table.Model
	progress progress.Model
	percent  float64
	symbol   string
}

func (m model) Init() tea.Cmd { return tick() }

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" { return m, tea.Quit }
	case tickMsg:
		m.percent += 1.0 / float64(refreshEvery)
		
		if m.percent >= 1.0 {
			m.percent = 0
			m.table.SetRows(fetchRows(m.symbol))
		}
		return m, tick()
	}
	return m, nil
}

func fetchRows(symbol string) []table.Row {
	var rows []table.Row
	score := 0

	var wf []WhaleFlow
	if fetch("/api/v2/spot/market/whale-net-flow?symbol="+symbol, &wf) == nil && len(wf) > 0 {
		v, _ := strconv.ParseFloat(wf[0].Volume, 64)
		sig := "SELL"; if v > 0 { sig = "BUY"; score++ } else { score-- }
		rows = append(rows, table.Row{"Whale Flow", wf[0].Volume, sig})
	}

	var tv []TakerVolume
	if fetch("/api/v2/mix/market/taker-buy-sell?symbol="+symbol, &tv) == nil && len(tv) > 0 {
		b, _ := strconv.ParseFloat(tv[0].BuyVolume, 64)
		s, _ := strconv.ParseFloat(tv[0].SellVolume, 64)
		sig := "AGG. SELL"; if b > s { sig = "AGG. BUY"; score++ } else { score-- }
		rows = append(rows, table.Row{"Taker Aggro", fmt.Sprintf("B:%.2f S:%.2f", b, s), sig})
	}

	var ls []LongShortRatio
	if fetch("/api/v2/mix/market/long-short?symbol="+symbol, &ls) == nil && len(ls) > 0 {
		r, _ := strconv.ParseFloat(ls[0].LongShortRatio, 64)
		sig := "NEUTRAL"
		if r > 1.2 { sig = "TRAP LONG"; score-- } else if r < 0.8 { sig = "TRAP SHORT"; score++ }
		rows = append(rows, table.Row{"L/S Ratio", fmt.Sprintf("%.2f", r), sig})
		rows = append(rows, table.Row{"Balance", drawGauge(r), ""})
	}

	bias := "NEUTRAL"
	if score > 0 { bias = "BULLISH" } else if score < 0 { bias = "BEARISH" }
	rows = append(rows, table.Row{"OVERALL BIAS", fmt.Sprintf("Score: %d", score), bias})
	return rows
}

func fetch(path string, target interface{}) error {
	resp, err := http.Get(BaseURL + path)
	if err != nil { return err }
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res BitgetRes
	json.Unmarshal(body, &res)
	return json.Unmarshal(res.Data, target)
}

func drawGauge(ratio float64) string {
	const size = 15
	norm := (ratio - 0.5) / (1.5 - 0.5)
	pos := int(norm * size)
	if pos < 0 { pos = 0 } else if pos > size { pos = size }
	bar := ""
	for i := 0; i <= size; i++ {
		if i == pos { bar += "X" } else if i == size/2 { bar += "|" } else { bar += "-" }
	}
	return "[" + bar + "]"
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n\n" +
		" Refreshing in:\n " + m.progress.ViewAs(m.percent) + "\n\n" +
		" 'q' to quit. Target: INSTITUTIONAL FLOW\n"
}


func main() {
	columns := []table.Column{
		{Title: "Metric", Width: 15},
		{Title: "Value", Width: 20},
		{Title: "Signal", Width: 15},
	}

	t := table.New(table.WithColumns(columns), table.WithRows(fetchRows("BTCUSDT")), table.WithHeight(6))
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Bold(true)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	p := progress.New(progress.WithDefaultGradient())

	m := model{table: t, progress: p, symbol: "BTCUSDT"}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
