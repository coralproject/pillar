package service

import (
	"fmt"

	"github.com/coralproject/pillar/pkg/stats"
)

// CalculateStats calculate stats as a service to use in the cron scheduler
func CalculateStats() {

	fmt.Println("Calculating Stats")
	stats.Init()
	stats.Calculate()
}
