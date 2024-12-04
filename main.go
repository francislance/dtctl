package main

import (
    "fmt"
    "os"

    "dtctl/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
