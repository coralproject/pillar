package statsd

import (
	"github.com/cactus/go-statsd-client/statsd"
	"log"
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

	// A buffered client, which sends multiple stats in one packet, is
	// recommended when your server supports it (better performance).
	// client, err := statsd.NewBufferedClient("127.0.0.1:8125", "test-client", 300*time.Millisecond, 0)

	if client != nil {
		return client
	}

	client, err := statsd.NewClient(url, "pillar")

	if err != nil {
		log.Printf("Error connecting to Statsd [%v]", err)
		return nil
	}

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
