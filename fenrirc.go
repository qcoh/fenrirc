package main

import (
	"fenrirc/mondrian"
)

func main() {
	if err := mondrian.Init(); err != nil {
		panic(err)
	}
	defer mondrian.Close()
	NewApplication().Run()
}
