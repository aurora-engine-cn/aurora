package main

import "gitee.com/aurora-engine/aurora"

func main() {
	// test
	a := aurora.NewAurora()

	go aurora.Run(a)
}
