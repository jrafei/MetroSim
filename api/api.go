package api

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"log"
	sim "metrosim/internal/simulation"
	"net/http"
	"time"
)

func StartAPI(sim *sim.Simulation) {
	mux := http.NewServeMux()
	port := "12000"
	station := func(w http.ResponseWriter, r *http.Request) {
		msg, _ := json.Marshal(sim.Print())
		fmt.Fprintf(w, "%s", msg)
	}

	mux.HandleFunc("/sim", station)
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	s := &http.Server{
		Addr:           ":" + "12000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	log.Println(fmt.Sprintf("Listening on localhost:%s", port))
	log.Fatal(s.ListenAndServe())
}
