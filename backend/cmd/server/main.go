package main

import (
	"app/internal/network"
	"log"
	"os"
)

func init() {
    // Create `data` folder
    if err := os.MkdirAll("data", os.ModePerm); err != nil {
        log.Fatalf("failed to create data folder: %v", err)
    }

    // Create `data/output.csv` file if it doesn't exist
    if _, err := os.Stat("./data/output.csv"); os.IsNotExist(err) {
        f, err := os.Create("./data/output.csv")
        if err != nil {
            log.Fatalf("failed to create output.csv: %v", err)
        }
        f.Close()
    }
}

func main() {
    go network.Start()

    select {}
}
