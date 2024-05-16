package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rdawson46/stock-ticker/api"

	"github.com/guptarohit/asciigraph"
)

type model struct {
    timer       chan timerMsg
    current     int
    stockSym    []string
    stocks      map[string][]float64 
    opens       map[string]float64 
    api         *api.Api
    time        time.Time
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
    border := lipgloss.RoundedBorder()
    border.BottomLeft = left
    border.Bottom = middle
    border.BottomRight = right
    return border

}

var (
    inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
    activeTabBorder = tabBorderWithBottom("┘", " ", "└")
    docStyle = lipgloss.NewStyle().Padding(2, 2, 2, 2)
    highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
    inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 2)
    activeTabStyle = inactiveTabStyle.Copy().Border(activeTabBorder, true)
    windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(1, 2).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func InitialModel() model {
    api, err := api.NewApi()

    config, err := GetConfigs()

    if err != nil {
        fmt.Println("Couldn't load config")
        os.Exit(1)
    }

    if err != nil {
        os.Exit(1)
    }

    // stockSymbols := []string{"CROX", "AAPL", "VOO"}
    stockSymbols := make([]string, len(config.Stocks))

    for i, stock := range config.Stocks {
        stockSymbols[i] = stock.Name
    }

    stocks := make(map[string][]float64)
    opens := make(map[string]float64)

    current := ""

    for _, stock := range stockSymbols {
        if current == "" {
            current = stock
        }
        
        price, err := api.GetOpeningPrice(stock)
        
        if err != nil {
            os.Exit(1)
        }
        
        stocks[stock] = append(make([]float64, 0), price)
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
        time: time.Now(),
    }
}


// only runs these functions on startup
func (m model) Init() tea.Cmd {
    cmds := make([]tea.Cmd, 0)

    for stock := range m.stocks {
        cmds = append(cmds, getStartPrices(m.api, stock))
    }

    return tea.Batch(
        tea.Batch(cmds...),
        timer(m.timer),
        waitForTimer(m.timer),
    )
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case startPricesMsg:
        for _, price := range msg.prices {
            m.stocks[msg.stock] = append(m.stocks[msg.stock], price)
        }
        return m, nil
    case priceMsg:
        m.stocks[msg.stock] = append(m.stocks[msg.stock], msg.price)
        return m, nil

    case timerMsg:
        m.time = time.Now()
        cmds := make([]tea.Cmd, 0)

        for stock := range m.stocks {
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
    /* TODO: 
     * add style for full window
         * fix border around screen
     * conditional coloring, overall coloring
    */
    doc := strings.Builder{}

    var renderedStocks []string

    for i, stock := range m.stockSym {
        var style lipgloss.Style
        isFirst, isLast, isActive := i == 0, i == len(m.stockSym)-1, i == m.current
        if isActive {
            style = activeTabStyle.Copy()
        } else {
            style = inactiveTabStyle.Copy()
        }

        border, _, _, _, _ := style.GetBorder()
        if isFirst && isActive {
            border.BottomLeft = "|"
        } else if isFirst && !isActive {
            border.BottomLeft = "├"
        } else if isLast && isActive {
            border.BottomRight = "|"
        } else if isLast && !isActive {
            border.BottomRight = "┤"
        }

        style = style.Border(border)
        renderedStocks = append(renderedStocks, style.Render(stock))
    }

    symbol, open, prices, currentTime := m.getData()
    main_content := fmt.Sprintf(
        "Stock: %s\nOpening: %f\nCurrent: %f\nTime: %s\n\n",
        symbol,
        open,
        prices[len(m.stocks[m.stockSym[m.current]]) - 1],
        currentTime,
    )

    var g string

    if open > prices[len(prices) - 1] {
        g = asciigraph.Plot(
            prices,
            asciigraph.Precision(6),
            asciigraph.SeriesColors(asciigraph.Red),
            asciigraph.Width(len(prices)),
            asciigraph.Height(15),
        )
    } else {
        g = asciigraph.Plot(
            prices,
            asciigraph.Precision(6),
            asciigraph.SeriesColors(asciigraph.Green),
            asciigraph.Width(len(prices)),
            asciigraph.Height(15),
        )
    }

    row := lipgloss.JoinHorizontal(lipgloss.Top, renderedStocks...)
    doc.WriteString(row)
    doc.WriteString("\n")
    //doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(main_content + g))
    doc.WriteString(windowStyle.Render(main_content + g))

    return docStyle.Render(doc.String())
}


func (m model) getData() (string, float64, []float64, string){
    return m.stockSym[m.current],
           m.opens[m.stockSym[m.current]],
           m.stocks[m.stockSym[m.current]],
           m.time.Format(time.ANSIC)
}
