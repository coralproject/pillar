package statsd

import (
	"github.com/cactus/go-statsd-client/statsd"
	"log"
	"time"
)

func main() {
}

type SD struct {
	Client statsd.Statter
}

var (
	client statsd.Statter
)

func connect(url string) statsd.Statter {

	if client != nil {
		return client
	}

	client, err := statsd.NewBufferedClient(url, "pillar", 300*time.Millisecond, 0)

	if err != nil {
		log.Printf("Error connecting to Statsd [%v]", err)
		return nil
	}

	log.Printf("Connected to Statsd")
	return client
}

func NewSD(url string) *SD {

	sd := SD{nil}

	client := connect(url)
	if client == nil {
		return &sd
	}

	sd.Client = client

	return &sd
}

func (sd *SD) Close() {

	if sd.Client == nil {
		return
	}

	sd.Client.Close()
}
