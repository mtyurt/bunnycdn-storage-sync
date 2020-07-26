package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mtyurt/bunnycdn-storage-sync/api"
	"github.com/mtyurt/bunnycdn-storage-sync/syncer"
)

func main() {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "Show the difference and exit")
	flag.Parse()
	src := flag.Arg(0)
	zoneName := flag.Arg(1)
	apiKey := os.Getenv("BCDN_APIKEY")
	if apiKey == "" {
		fmt.Println("API key must be set with BCDN_APIKEY")
		os.Exit(-1)
	}
	fmt.Println("Walking path:", src)
	storage := api.BCDNStorage{ZoneName: zoneName, APIKey: apiKey}
	syncerService := syncer.BCDNSyncer{storage, dryRun}
	err := syncerService.Sync(src)
	fmt.Println(err)
}
