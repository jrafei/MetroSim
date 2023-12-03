package simulation

import (
	"fmt"
	"log"
	//"math/rand"
	"sync"
	"time"
)

// Déclaration de la matrice
/*
 * X : Mur, zone inatteignable
 * E : Entrée
 * S : Sortie
 * W : Entrée et Sortie
 * Q : Voie
 * _ : Couloir, zon
 * B: Bridge/Pont, zone accessible
 */
var carte [20][20]string = [20][20]string{
	{"X", "X", "X", "X", "X", "X", "X", "X", "W", "W", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "_", "X"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "_"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "S", "S", "X", "X", "X", "X", "X", "X", "E", "E", "X", "X", "X", "X", "X", "X"},
}

var playground [20][20]string = [20][20]string{
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
}

type Simulation struct {
	env         Environment
	agents      []Agent
	maxStep     int
	maxDuration time.Duration
	step        int // Stats
	start       time.Time
	syncChans   sync.Map
}

func NewSimulation(agentCount int, maxStep int, maxDuration time.Duration) (simu *Simulation) {
	simu = &Simulation{}
	simu.agents = make([]Agent, 0, agentCount)
	simu.maxStep = maxStep
	simu.maxDuration = maxDuration

	simu.env = *NewEnvironment([]Agent{}, carte)
	//simu.env = *NewEnvironment([]Agent{}, playground)

	// création des agents et des channels
	for i := 0; i < agentCount; i++ {
		// création de l'agent
		id := fmt.Sprintf("Agent #%d", i)
		syncChan := make(chan int)
		ag := NewAgent(id, &simu.env, syncChan, 1, 0, true, Coord{0, 8 + i%2}, Coord{0, 8 + i%2}, &UsagerLambda{},Coord{0, 8 + i%2}, Coord{12 - 4*(i%2), 18 - 15*(i%2)})

		//ag := NewAgent(id, &simu.env, syncChan, 1, 0, true, Coord{4,10}, Coord{4,10}, &UsagerLambda{}, Coord{4,10}, Coord{0, 0})

		// ajout de l'agent à la simulation
		simu.agents = append(simu.agents, *ag)

		// ajout du channel de synchro
		simu.syncChans.Store(ag.ID(), syncChan)

		// ajout de l'agent à l'environnement
		ag.env.AddAgent(*ag)
	}

	return simu
}

func (simu *Simulation) Run() {
	log.Printf("Démarrage de la simulation [step: %d, π: %f]", simu.step, simu.env.PI())

	// Démarrage du micro-service de Log
	go simu.Log()
	// Démarrage du micro-service d'affichage
	go simu.Print()

	// Démarrage des agents
	
	var wg sync.WaitGroup
	for _, agt := range simu.agents {
		wg.Add(1)
		go func(agent Agent) {
			defer wg.Done()
			agent.Start()
		}(agt)
	}
	
	// On sauvegarde la date du début de la simulation
	simu.start = time.Now()

	// Lancement de l'orchestration de tous les agents
	// simu.step += 1 // plus de sens
	for _, agt := range simu.agents {
		go func(agt Agent) {
			step := 0
			for {
				step++
				c, _ := simu.syncChans.Load(agt.ID())
				c.(chan int) <- step             // /!\ utilisation d'un "Type Assertion"
				time.Sleep(1 * time.Millisecond) // "cool down"
				<-c.(chan int)
			}
		}(agt)
	}

	time.Sleep(simu.maxDuration)

	log.Printf("Fin de la simulation [step: %d, in: %d, out: %d, π: %f]", simu.step, simu.env.PI())
}

func (simu *Simulation) Print() {
	for {
		for i := 0; i < 20; i++ {
			fmt.Println(simu.env.station[i])
		}
		//time.Sleep(time.Second / 4) // 60 fps !
		time.Sleep(time.Second)    // 1 fps !
		//fmt.Print("\033[H\033[2J") // effacement du terminal
	}
}

func (simu *Simulation) Log() {
	// Not implemented
}
