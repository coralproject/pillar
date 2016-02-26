package backend

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/newsdev/golympics/data"
// 	"github.com/newsdev/golympics/data/backend/mongodb"
// )

// var (
// 	ObjectTypes = []string{
// 		"bracket",
// 		"competitor_stats",
// 		"configs",
// 		"current_periods",
// 		"current_result",
// 		"extended_infos",
// 		"feed",
// 		"message",
// 		"officials",
// 		"participant",
// 		"periods",
// 		"pool_result",
// 		"ranking_result",
// 		"schedule_unit",
// 		"team",
// 		"unit_actions",
// 		"unit_result",
// 	}
// )

// func testBackendClose(b Backend, t *testing.T) {
// 	if err := b.Close(); err != nil {
// 		t.Error(err)
// 	}
// }

// func testBrackets(b Backend, t *testing.T) {
// 	bracketIterator, err := b.Brackets()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for {
// 		e, err := bracketIterator.Next()
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		if e == nil {
// 			break
// 		}

// 		if b, ok := e.(*data.Bracket); ok {
// 			fmt.Println(b)
// 		} else {
// 			t.Error("wrong data type from iterator")
// 		}
// 	}

// 	if err := bracketIterator.Close(); err != nil {
// 		t.Error(err)
// 	}
// }

// func testCompetitorStats(b Backend, t *testing.T) {
// 	competitorStatsIterator, err := b.CompetitorStats()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for {
// 		e, err := competitorStatsIterator.Next()
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		if e == nil {
// 			break
// 		}

// 		if b, ok := e.(*data.CompetitorStats); ok {
// 			fmt.Println(b)
// 		} else {
// 			t.Error("wrong data type from iterator")
// 		}
// 	}

// 	if err := competitorStatsIterator.Close(); err != nil {
// 		t.Error(err)
// 	}
// }

// func testBackend(b Backend, t *testing.T) {
// 	testBrackets(b, t)
// 	testCompetitorStats(b, t)
// 	testBackendClose(b, t)
// }

// func TestMongoDBBackend(t *testing.T) {
// 	b, err := mongodb.NewMongoDBBackend("mongodb://127.0.0.1:27017", "olympics")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	testBackend(b, t)
// }
