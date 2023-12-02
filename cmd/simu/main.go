package main

import (
	simulation "metrosim"
	"time"
)

func main() {
	s := simulation.NewSimulation(30, -1, 600*time.Second)
	//go simulation.StartAPI(s)
	s.Run()
}
