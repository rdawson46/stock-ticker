package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rdawson46/stock-ticker/api"
)

// HACK: have counter to register timerMsg
type model struct {
    // might need to make a slice with symbols and pointer val
    timer       chan timerMsg
    current     int
    stockSym    []string
    stocks      map[string][]float32 // maps to array of prices
    opens       map[string]float32 // change to map
    api         *api.Api
    counter     int // remove
}

func initialModel() model {
    api, err := api.NewApi()

    if err != nil {
        os.Exit(1)
    }

    stockSymbols := []string{"CROX", "AAPL", "VOO"}

    stocks := make(map[string][]float32)
    opens := make(map[string]float32)

    current := ""

    for _, stock := range stockSymbols {
        if current == "" {
            current = stock
        }
        
        price, err := api.GetOpeningPrice(stock)
        
        if err != nil {
            os.Exit(1)
        }
        
        stocks[stock] = append(make([]float32, 0), price)
        opens[stock] = price
    }

    if err != nil {
        os.Exit(1)
    }

    return model{
        timer: make(chan timerMsg),
        current: 0,
        stocks: stocks,
        stockSym: stockSymbols,
        opens: opens,
        api: api,
        counter: 0,
    }
}


// only runs these functions on startup
func (m model) Init() tea.Cmd {
    cmds := make([]tea.Cmd, 0)

    for stock, _ := range m.stocks {
        cmds = append(cmds, getPrice(m.api, stock))
    }

    return tea.Batch(
        tea.Batch(cmds...),
        timer(m.timer),
        waitForTimer(m.timer),
    )
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case priceMsg:
        m.stocks[msg.stock] = append(m.stocks[msg.stock], float32(msg.price))
        return m, nil

    case timerMsg:
        m.counter++
        cmds := make([]tea.Cmd, 0)

        for stock, _ := range m.stocks {
            cmds = append(cmds, getPrice(m.api, stock))
        }

        return m, tea.Batch(tea.Batch(cmds...), waitForTimer(m.timer))

    case errMsg:
        return m, tea.Quit
         
    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }

        switch msg.String() {
        case "q":
            return m, tea.Quit
        case "n":
            m.current = (m.current + 1) % len(m.stockSym)
        case "p":
            if m.current == 0 {
                m.current = len(m.stockSym) - 1
            } else {
                m.current -= 1
            }
        }
    }
    return m, nil
}

func (m model) View() string {
    // TODO: make a pleasing UI
    // add border around screen
    // add tabs
    // center content
    // conditional coloring
    // graphing

    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("9")).
        Bold(true).
        Padding(2).
        Align(lipgloss.Center)

    main_content := fmt.Sprintf(
        "Stock: %s\nOpening: %f\nCurrent: %f\nCounter: %d\n",
        m.stockSym[m.current],
        m.opens[m.stockSym[m.current]],
        m.stocks[m.stockSym[m.current]],
        m.counter,
    )

    return style.Render(main_content)
}

type timerMsg struct{}

type priceMsg struct{
    stock string
    price float32
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

func main() {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())

    if _, err := p.Run(); err != nil {
        fmt.Println("broke:", err)
        os.Exit(1)
    }
}
