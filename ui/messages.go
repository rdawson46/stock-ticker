package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rdawson46/stock-ticker/api"
)

type timerMsg struct{}

type priceMsg struct{
    stock string
    price float64
}

type startPricesMsg struct{
    stock  string
    prices []float64
}

type errMsg struct {
    err error
}

func getPrice(api *api.Api, stock string) tea.Cmd {
    return func() tea.Msg {
        price, err := api.GetPrice(stock)
        if err != nil {
            return errMsg{err}
        }
        return priceMsg{stock, price}
    }
}

func getStartPrices(api *api.Api, stock string) tea.Cmd {
    return func() tea.Msg {
        prices, err := api.GetLast10(stock)
        if err != nil {
            return errMsg{err}
        }
        return startPricesMsg{
            prices: prices,
            stock: stock,
        }
    }
}

func timer(sub chan timerMsg) tea.Cmd {
    return func() tea.Msg {
        for{
            time.Sleep(time.Minute)
            sub <- timerMsg(struct{}{})
        }
    }
}

func waitForTimer(sub chan timerMsg) tea.Cmd {
    return func() tea.Msg {
        return <-sub
    }
}
