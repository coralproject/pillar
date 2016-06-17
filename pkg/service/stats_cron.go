package service

import (
	"log"

	"github.com/coralproject/pillar/pkg/stats"
)

// CalculateStats calculate stats as a service to use in the cron scheduler
func CalculateStats() {

	log.Printf("New scheduled job - Calculating Stats!\n")

	stats.Init()
	stats.Calculate()
}
