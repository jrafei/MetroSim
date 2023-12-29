package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	sim "metrosim/internal/simulation"
)

func StartAPI(sim *sim.Simulation) {
	mux := http.NewServeMux()

	station := func(w http.ResponseWriter, r *http.Request) {
		msg, _ := json.Marshal(sim.Print())
		fmt.Fprintf(w, "%s", msg)
	}

	mux.HandleFunc("/sim", station)

	s := &http.Server{
		Addr:           ":12000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	log.Println("Listening on localhost:8080")
	log.Fatal(s.ListenAndServe())
}
