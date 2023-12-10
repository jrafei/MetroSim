package main

import (
	simulation "metrosim/internal/simulation"
	"time"
)

func main() {
	s := simulation.NewSimulation(1, -1, 600*time.Second)
	//go simulation.StartAPI(s)
	s.Run()
}
