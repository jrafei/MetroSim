package simulation

import (
	"log"
	"math/rand"
	"time"
)

var metro_speed int = 5 // Nombre de seconde de l'entrée en gare

type Metro struct {
	frequency  time.Duration
	stopTime   time.Duration
	freeSpace  int // nombre de cases disponibles dans le métro
	comChannel chan Request
	way        *Way
}

func NewMetro(freq time.Duration, stopT time.Duration, freeS int, way *Way) *Metro {
	return &Metro{
		frequency:  freq,
		stopTime:   stopT,
		freeSpace:  freeS,
		comChannel: make(chan Request),
		way:        way,
	}
}

func (metro *Metro) Start() {
	log.Printf("Metro starting...\n")
	refTime := time.Now()
	go func() {
		//var step int
		for {
			//step = <-metro.syncChan
			if refTime.Add(metro.frequency).Sub(time.Now()) <= time.Duration(metro_speed)*time.Second {
				metro.printMetro()
			}
			if refTime.Add(metro.frequency).Before(time.Now()) {
				metro.dropUsers()
				metro.pickUpUsers()
				metro.removeMetro()
				metro.freeSpace = rand.Intn(10)
				refTime = time.Now()
			}
			//metro.syncChan <- step

		}
	}()
}

func (metro *Metro) pickUpUsers() {
	// Faire monter les usagers dans le métro
	t := time.Now()
	for time.Now().Before(t.Add(metro.stopTime)) {
		if metro.freeSpace > 0 {
			for _, gate := range metro.way.gates {
				go metro.pickUpGate(&gate)
			}
		}
	}
}

func (metro *Metro) pickUpGate(gate *Coord) {
	// Récupérer les usagers à une porte spécifique
	gate_cell := metro.way.env.station[gate[0]][gate[1]]
	if len(gate_cell) > 1 {
		agent := metro.findAgent(AgentID(gate_cell))
		if agent != nil && agent.width*agent.height <= metro.freeSpace && agent.destination == *gate {
			metro.way.env.agentsChan[agent.id] <- *NewRequest(metro.comChannel, Disappear)
			metro.freeSpace--
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

func (metro *Metro) dropUsers() {

}
