package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rishijain300900/nse-margin-datafetch/pkg/scheduler"
)

func main() {
	scheduler.PerformScheduling()
}
