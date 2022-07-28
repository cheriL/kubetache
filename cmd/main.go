package main

import (
	"github.com/cheriL/kubetache/action"
	"time"
)

func main() {
	action.Init(1)

	for  {
		time.Sleep(3000)

		action.Tache("myharbor-harbor-core-c5cd8974-wzth6", "")
	}
}
