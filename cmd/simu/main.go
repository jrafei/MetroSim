package main

import (
	simulation "metrosim"
	"time"
)

func main() {
	s := simulation.NewSimulation(40, -1, 600*time.Second)
	//go simulation.StartAPI(s)
	s.Run()
}
