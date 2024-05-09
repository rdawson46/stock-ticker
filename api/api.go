/*

TODO:
* learn the v3 api
* save keys to .env
* way to make calls

*/

package api

import (
    "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
    "github.com/joho/godotenv"
    "os"
)

type Api struct {
    Endpoint *alpaca.Client
}

func NewApi() (*Api, error) {
    err := godotenv.Load("../.env")

    if err != nil {
        return nil, err
    }

    key := os.Getenv("APIKEY")
    private := os.Getenv("APIPRIVATE")

    ep := alpaca.NewClient(alpaca.ClientOpts{
        APIKey: key,
        APISecret: private,
        BaseURL: "https://paper-api.alpaca.markets",
    })

    return &Api{ep}, nil
}

func (l *Api) CheckAccount() error {
    _, err := l.Endpoint.GetAccount()

    return err
}
