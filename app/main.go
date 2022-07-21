package main

import "gitee.com/aurora-engine/aurora"

func main() {
	a := aurora.NewAurora()

	go aurora.Run(a)
}
