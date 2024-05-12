package main

import (
    "time"
	"fmt"
	"os"
	// "github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rdawson46/stock-ticker/api"
)

type model struct {
    stock   string
    open    float32
    current float32
    api     *api.Api
}

func initialModel() model {
    api, err := api.NewApi()

    if err != nil {
        os.Exit(1)
    }

    stock := "crox"

    open, err := api.GetOpeningPrice(stock)

    if err != nil {
        os.Exit(1)
    }

    return model{
        stock: stock,
        open: open,
        current: open,
        api: api,
    }
}

func (m model) Init() tea.Cmd {
    return getPrice(m.api, m.stock)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
    switch msg := msg.(type) {
    case priceMsg:
        m.current = float32(msg)
        return m, nil
    }
    return m, nil
}

func (m model) View() string {
    return ""
}


// TODO: merge these next two functions
func getPrice(api *api.Api, stock string) tea.Cmd {
    return func() tea.Msg {
        price, err := api.GetPrice(stock)

        api.Count--

        if err != nil {
            return errMsg{err}
        }
        return priceMsg(price)
    }
}

func timer(api *api.Api) tea.Cmd {
    return func() tea.Msg {
        time.Sleep(time.Minute)
        api.Count = 200
        return resetMsg(true)
    }
}

type resetMsg bool

type priceMsg float32

type errMsg struct {
    err error
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Println("broke:", err)
        os.Exit(1)
    }
}
