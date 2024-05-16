package ui

import (
	"encoding/json"
	"fmt"
	"os"
    "strings"
)

type Config struct {
    Stocks   []Stock   `json:"stocks"`
    Colors   []int     `json:"colors"`
}

type Stock struct {
    Name        string  `json:"name"`
    Position    float32 `json:"position"`
}

func GetConfigs() (*Config, error){
    dat, err := os.ReadFile("./config.json")

    if err != nil {
        fmt.Println("couldn't open file")
        return nil, err
    }

    c := Config{}
    json.Unmarshal(dat, &c)

    for i, stock := range c.Stocks {
        c.Stocks[i].Name = strings.ToUpper(stock.Name)
    }
    
    return &c, nil
}
