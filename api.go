package simulation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartAPI(sim *Simulation) {
	mux := http.NewServeMux()

	pi := func(w http.ResponseWriter, r *http.Request) {
		msg, _ := json.Marshal(sim.env.PI())
		fmt.Fprintf(w, "%s", msg)
	}

	mux.HandleFunc("/pi", pi)

	s := &http.Server{
		Addr:           ":12000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	log.Println("Listening on localhost:8080")
	log.Fatal(s.ListenAndServe())
}
