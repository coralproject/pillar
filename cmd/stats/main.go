// It provides teh coral project sponge CLI  platform.
package main

import (
	"os"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/cmd/stats/cmd"
	"github.com/pborman/uuid"
)

func main() {

	uid := uuid.New()

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Error(uid, "main", err, "Unable to execute the command.")
		os.Exit(-1)
	}
}
