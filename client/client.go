package main

import (
	"github.com/coralproject/pillar/client/fiddler"
)

func main() {

    //Load WAPO users
    fiddler.LoadActors()

    //Load WAPO comments
    fiddler.LoadComments()
}
