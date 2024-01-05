package api

import (
	"bytes"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"log"
	req "metrosim/internal/request"
	sim "metrosim/internal/simulation"
	"net/http"
	"time"
)

var simulation *sim.Simulation = nil

func StartAPI() {
	mux := http.NewServeMux()
	port := "12000"

	mux.HandleFunc("/configure", simHandler("configure"))
	mux.HandleFunc("/launch", simHandler("launch"))
	//mux.HandleFunc("/stop", simHandler("stop"))

	s := &http.Server{
		Addr:           ":" + "12000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	log.Println(fmt.Sprintf("Listening on localhost:%s", port))
	log.Fatal(s.ListenAndServe())
}

func decodeConf(r *http.Request) (conf req.Cfg, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &conf)
	fmt.Println(err)
	return
}

func checkMethod(method string, r *http.Request) bool {
	return r.Method == method
}

func simHandler(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		switch action {
		case "configure":

			if simulation != nil && simulation.IsRunning() {
				log.Println("[API] Must stop current sim to configure")
			} else {
				configuration(w, r)

			}

		case "launch":
			//TODELETElog.Println("[API] Launch request")
			if simulation != nil {
				//TODELETElog.Println("[API] running")
				if !simulation.IsRunning() {
					//TODELETElog.Println("[API] running")
					go simulation.Run()
				}
				msg, _ := json.Marshal(simulation.Print())
				//TODELETEfmt.Printf("%d\n", len(msg))
				fmt.Fprintf(w, "%s", msg)
			}
		case "stop":
			// TODOD
		}
	}
}

func configuration(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête
	if !checkMethod("POST", r) {
		return
	}
	conf, err := decodeConf(r)
	if err != nil {
		log.Println("[API] Not able to read configuration data")
		simulation = nil
		return
	}
	log.Println("[API] Received configuration data : ", conf)
	size := len(conf.LeftTopCorners)
	if len(conf.RightBottomCorners) != size || len(conf.Gates) != size || len(conf.GoToLeft) != size || len(conf.Frequency) != size || len(conf.StopTime) != size || len(conf.Capacity) != size {
		simulation = nil
		return
	}
	log.Println("[API] Simulation initialized")
	simulation = sim.NewSimulation(conf)
}
