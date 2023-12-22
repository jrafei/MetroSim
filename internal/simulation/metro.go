package simulation

import (
	"fmt"
	"log"
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"sync"
	"time"
)

/*
 * //TODO:Ajouter la capacité max
 * //TODO:Parfois agents bloqués
 * // Apparition des agents sortant
 */

var metro_speed int = 5 // Nombre de seconde de l'entrée en gare

type Metro struct {
	frequency  time.Duration
	stopTime   time.Duration
	capacity   int
	freeSpace  int // nombre de cases disponibles dans le métro
	comChannel chan req.Request
	way        *Way
}

func NewMetro(freq time.Duration, stopT time.Duration, capacity, freeS int, way *Way) *Metro {
	return &Metro{
		frequency:  freq,
		stopTime:   stopT,
		capacity:   capacity,
		freeSpace:  freeS,
		comChannel: make(chan req.Request),
		way:        way,
	}
}

func (metro *Metro) Start() {
	log.Printf("Metro starting...\n")
	refTime := time.Now()
	//var step int
	for {
		//step = <-metro.syncChan
		if refTime.Add(metro.frequency).Sub(time.Now()) <= time.Duration(metro_speed)*time.Second {
			metro.printMetro()
		}
		if refTime.Add(metro.frequency).Before(time.Now()) {
			metro.dropUsers()
			metro.openGates()
			metro.pickUpUsers()
			metro.closeGates()
			metro.removeMetro()
			metro.freeSpace = rand.Intn(10)
			refTime = time.Now()
		}
		//metro.syncChan <- step

	}
}

func (metro *Metro) pickUpUsers() {
	var wg sync.WaitGroup
	for _, gate := range metro.way.gates {
		wg.Add(1)
		go func(gate alg.Coord) {
			defer wg.Done()
			metro.pickUpGate(&gate, time.Now().Add(metro.stopTime))
		}(gate)
	}

	wg.Wait()
}

func (metro *Metro) pickUpGate(gate *alg.Coord, endTime time.Time) {
	// Récupérer les usagers à une porte spécifique
	for {

		if !time.Now().Before(endTime) {
			return
		} else {
			gate_cell := metro.way.env.station[gate[0]][gate[1]]
			if len(gate_cell) > 1 {
				agent := metro.findAgent(AgentID(gate_cell))
				if agent != nil && agent.width*agent.height <= metro.freeSpace && alg.EqualCoord(&agent.destination, gate) {
					fmt.Println("agent entering metro : ", agent.id, "at gate ", gate)
					metro.way.env.agentsChan[agent.id] <- *req.NewRequest(metro.comChannel, EnterMetro)
					<-metro.comChannel
					metro.freeSpace = metro.freeSpace - agent.width*agent.height
					//fmt.Println("leaving", agent.id)
				}
			}
		}
	}

}

func (metro *Metro) findAgent(agent AgentID) *Agent {
	// Trouver l'adresse de l'agent
	for _, agt := range metro.way.env.ags {
		if agt.id == agent {
			return &agt
		}
	}
	return nil
}

func (metro *Metro) dropUsers() {
	nb := rand.Intn(metro.capacity - metro.freeSpace) // Nombre de cases à vider du métro
	for nb > 0 {
		gate_nb := rand.Intn(len(metro.way.gates)) // Sélection d'une porte aléatoirement
		width := 1                                 //+ rand.Intn(2)
		height := 1                                //+ rand.Intn(2)
		metro.freeSpace = metro.freeSpace + width*height
		nb = nb - width*height
		id := fmt.Sprintf("Agent%d", metro.way.env.agentCount)
		path := metro.way.pathsToExit[gate_nb]
		ag := NewAgent(id, metro.way.env, make(chan int), 200, 0, true, &UsagerLambda{}, metro.way.gates[gate_nb], metro.way.nearestExit[gate_nb], width, height)
		ag.path = path
		metro.way.env.AddAgent(*ag)
		ag.env.writeAgent(ag)
		//log.Println(metro.way.id, nb, metro.way.env.agentCount)
		//fmt.Println("agent leaving metro", ag.id, ag.departure, ag.destination, width, height)
		time.Sleep(500 * time.Millisecond)
	}

}

func (metro *Metro) printMetro() {

	if metro.way.horizontal {
		waiting_time := time.Duration((metro_speed * 1000) / (metro.way.downRightCoord[1] - metro.way.upLeftCoord[1]))
		if metro.way.goToLeft {
			for y := metro.way.downRightCoord[1]; y >= metro.way.upLeftCoord[1]; y-- {
				for x := metro.way.upLeftCoord[0]; x <= metro.way.downRightCoord[0]; x++ {
					if metro.way.env.station[x][y] == "Q" {
						metro.way.env.station[x][y] = "M"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		} else {
			for y := metro.way.upLeftCoord[1]; y <= metro.way.downRightCoord[1]; y++ {
				for x := metro.way.upLeftCoord[0]; x <= metro.way.downRightCoord[0]; x++ {
					if metro.way.env.station[x][y] == "Q" {
						metro.way.env.station[x][y] = "M"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		}

	} else {
		waiting_time := time.Duration((metro_speed * 1000) / (metro.way.downRightCoord[0] - metro.way.upLeftCoord[0]))
		if metro.way.goToLeft {
			// de bas en haut
			for x := metro.way.downRightCoord[0]; x >= metro.way.upLeftCoord[0]; x-- {
				for y := metro.way.upLeftCoord[1]; y <= metro.way.downRightCoord[1]; y++ {
					if metro.way.env.station[x][y] == "Q" {
						metro.way.env.station[x][y] = "M"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		} else {
			for x := metro.way.upLeftCoord[0]; x <= metro.way.downRightCoord[0]; x++ {
				for y := metro.way.upLeftCoord[1]; y <= metro.way.downRightCoord[1]; y++ {
					if metro.way.env.station[x][y] == "Q" {
						metro.way.env.station[x][y] = "M"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		}

	}

}

func (metro *Metro) removeMetro() {

	if metro.way.horizontal {
		waiting_time := time.Duration((metro_speed * 1000) / (metro.way.downRightCoord[1] - metro.way.upLeftCoord[1]))

		if metro.way.goToLeft {
			for y := metro.way.downRightCoord[1]; y >= metro.way.upLeftCoord[1]; y-- {
				for x := metro.way.upLeftCoord[0]; x <= metro.way.downRightCoord[0]; x++ {
					if metro.way.env.station[x][y] == "M" {
						metro.way.env.station[x][y] = "Q"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		} else {
			for y := metro.way.upLeftCoord[1]; y <= metro.way.downRightCoord[1]; y++ {
				for x := metro.way.upLeftCoord[0]; x <= metro.way.downRightCoord[0]; x++ {
					if metro.way.env.station[x][y] == "M" {
						metro.way.env.station[x][y] = "Q"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		}

	} else {
		waiting_time := time.Duration((metro_speed * 1000) / (metro.way.downRightCoord[0] - metro.way.upLeftCoord[0]))
		if metro.way.goToLeft {
			// de bas en haut
			for x := metro.way.downRightCoord[0]; x >= metro.way.upLeftCoord[0]; x-- {
				for y := metro.way.upLeftCoord[1]; y <= metro.way.downRightCoord[1]; y++ {
					if metro.way.env.station[x][y] == "M" {
						metro.way.env.station[x][y] = "Q"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		} else {
			for x := metro.way.upLeftCoord[0]; x <= metro.way.downRightCoord[0]; x++ {
				for y := metro.way.upLeftCoord[1]; y <= metro.way.downRightCoord[1]; y++ {
					if metro.way.env.station[x][y] == "M" {
						metro.way.env.station[x][y] = "Q"
					}
				}
				time.Sleep(waiting_time * time.Millisecond)
			}
		}

	}
}

func (metro *Metro) openGates() {
	// Début d'autorisation d'entrer dans le métro
	for _, gate := range metro.way.gates {
		metro.way.env.station[gate[0]][gate[1]] = "O"
	}
	metro.way.gatesClosed = false
}

func (metro *Metro) closeGates() {
	// Fin d'autorisation d'entrer dans le métro
	metro.way.gatesClosed = true
	for _, gate := range metro.way.gates {
		if len(metro.way.env.station[gate[0]][gate[1]]) > 1 {
			// On autorise les agents déjà sur la case à rentrer dans le métro
			metro.pickUpGate(&gate, time.Now().Add(time.Duration(1*time.Second)))
		}
		metro.way.env.station[gate[0]][gate[1]] = "G"
	}

}
