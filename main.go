package main

import (
	"fmt"
	"os"
	"time"

	// "github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rdawson46/stock-ticker/api"
)

// HACK: have counter to register timerMsg
type model struct {
    timer   chan timerMsg
    stock   string
    open    float32
    current float32
    api     *api.Api
    counter int
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
        timer: make(chan timerMsg),
        stock: stock,
        open: open,
        current: open,
        api: api,
        counter: 0,
    }
}


// only runs these functions on startup
func (m model) Init() tea.Cmd {
    return tea.Batch(
        getPrice(m.api, m.stock),
        timer(m.timer),
        waitForTimer(m.timer),
    )
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case priceMsg:
        m.current = float32(msg)
        return m, nil

    case timerMsg:
        m.counter++
        return m, tea.Batch(waitForTimer(m.timer), getPrice(m.api, m.stock))

    case errMsg:
        return m, tea.Quit
         
    case tea.KeyMsg:
        // TODO: add quit with q
        if msg.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    s := fmt.Sprintf("Opening: %f\nCurrent: %f\nCounter: %d\n", m.open, m.current, m.counter)
    return s
}

type timerMsg struct{}

type priceMsg float32

type errMsg struct {
    err error
}

func getPrice(api *api.Api, stock string) tea.Cmd {
    return func() tea.Msg {
        price, err := api.GetPrice(stock)
        if err != nil {
            return errMsg{err}
        }
        return priceMsg(price)
    }
}

func timer(sub chan timerMsg) tea.Cmd {
    return func() tea.Msg {
        for{
            time.Sleep(time.Second)
            sub <- timerMsg(struct{}{})
        }
    }
}

func waitForTimer(sub chan timerMsg) tea.Cmd {
    return func() tea.Msg {
        return <-sub
    }
}


func main() {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())

    if _, err := p.Run(); err != nil {
        fmt.Println("broke:", err)
        os.Exit(1)
    }
}
