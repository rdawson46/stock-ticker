package api

import (
	"fmt"
    "os"
	"github.com/joho/godotenv"
    "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
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

func (self *Api) GetOpeningPrice(stock string) (float32, error) {
    bar, err := self.client.GetLatestBar(stock, marketdata.GetLatestBarRequest{})

    if err != nil {
        fmt.Println(err)
        return 0.0, err
    }

    return float32(bar.Open), nil
}

func (self *Api) GetPrice(stock string) (float32, error) {
    bar, err := self.client.GetLatestBar(stock, marketdata.GetLatestBarRequest{})

    if err != nil {
        fmt.Println(err)
        return 0.0, err
    }

    return float32(bar.VWAP), nil
}

// func for getting lower limit of refresh rate
// important for when use for multiple stocks added
//func (self *Api) GetRateLimit(stocks []string, refreshRate int) int {
func (self *Api) GetRateLimit() float32 {
    //return (float32(60) / float32(200)) * 1000 / float32(len(stocks))
    return (float32(60) / float32(200)) * 1000
}

func (self *Api) GetPrices(stock string, count int) ([]float32, error) {
    // TODO: apply timelines
    return []float32 {0}, nil
}

