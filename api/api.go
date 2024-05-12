package api

import (
	"fmt"
    "os"
	"github.com/joho/godotenv"
    "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type Api struct {
    // use this to prevent api overloading, 200 calls/min
    Count int
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

    return &Api{200, client}, nil
}

func (self *Api) GetOpeningPrice(stock string) (float32, error) {
    bar, err := self.client.GetLatestBar("AAPL", marketdata.GetLatestBarRequest{})

    if err != nil {
        fmt.Println(err)
        return 0.0, err
    }

    return float32(bar.Open), nil
}

func (self *Api) GetPrice(stock string) (float32, error) {
    bar, err := self.client.GetLatestBar("AAPL", marketdata.GetLatestBarRequest{})

    if err != nil {
        fmt.Println(err)
        return 0.0, err
    }

    return float32(bar.VWAP), nil
}

func (self *Api) GetPrices(stock string, count int) ([]float32, error) {
    // TODO: apply timelines
    return []float32 {0}, nil
}

