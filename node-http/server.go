package main

import (
	"github.com/dejavuzhou/evefs/node-http/routers"
)

func main() {
	routers.Router.Run(":7777")
}
