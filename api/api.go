package api

import (
	"fmt"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/joho/godotenv"
)

type Api struct {
    client *marketdata.Client
}

func NewApi() (*Api, error) {
    err := godotenv.Load(".env")

    if err != nil {
        return nil, err
    }

    key := os.Getenv("APIKEY")
    private := os.Getenv("APIPRIVATE")

    client := marketdata.NewClient(marketdata.ClientOpts{
        APIKey: key,
        APISecret: private,
    })

    return &Api{client}, nil
}

func (self *Api) GetOpeningPrice(stock string) (float64, error) {
    bar, err := self.client.GetLatestBar(stock, marketdata.GetLatestBarRequest{})

    if err != nil {
        fmt.Println(err)
        return 0.0, err
    }

    return float64(bar.Open), nil
}

func (self *Api) GetPrice(stock string) (float64, error) {
    bar, err := self.client.GetLatestBar(stock, marketdata.GetLatestBarRequest{})

    if err != nil {
        fmt.Println(err)
        return 0.0, err
    }

    return float64(bar.VWAP), nil
}

// TODO: work on this, will be needed with expansion
func (self *Api) GetRateLimit() float64 {
    //return (float64(60) / float64(200)) * 1000 / float64(len(stocks))
    return (float64(60) / float64(200)) * 1000
}

func (self *Api) GetLast10(stock string) ([]float64, error){
    prices, err := self.client.GetBars(stock, marketdata.GetBarsRequest{
        Start: time.Now().Add(time.Duration(-45) * time.Hour),
        End: time.Now().Add(time.Duration(-15) * time.Hour),
        TimeFrame: marketdata.OneMin,
        TotalLimit: 30,
    })

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
        return []float64{}, err
    }

    values := make([]float64, 0)

    for _, val := range prices {
        values = append(values, val.VWAP)
    }
    
    return values, nil
}

