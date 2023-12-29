package main

import (
	"fmt"
	simulation "metrosim/internal/simulation"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
	"metrosim/api"
)

func main() {
	s := simulation.NewSimulation(20, -1, 600*time.Second)
	go api.StartAPI(s)
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Printf("Nombre de goroutines : %d\n", runtime.NumGoroutine())
			if runtime.NumGoroutine() > 1000 {

				f, err := os.Create("goroutines.pprof")
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not create goroutine profile: %v\n", err)
				}

				pprof.Lookup("goroutine").WriteTo(f, 1)
				f.Close()

			}
		}
	}()

	s.Run()
}
