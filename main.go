package main

import (
    "fmt"
    "github.com/rdawson46/stock-ticker/api"
)

func main() {
    x, err := api.NewApi()

    if err != nil {
        fmt.Println(err)
        return 
    }

    if err != nil {
        fmt.Println("no good")
        return 
    }

}
