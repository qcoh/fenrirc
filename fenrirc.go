package main

import (
	"fenrirc/mondrian"
	_ "fenrirc/plugins/rainbow"
)

func main() {
	if err := mondrian.Init(); err != nil {
		panic(err)
	}
	defer mondrian.Close()
	NewApplication().Run()
}
