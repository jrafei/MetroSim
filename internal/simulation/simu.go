package simulation

import (
	"fmt"
	"log"
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
 * _ : Couloir, case libre
 * B: Bridge/Pont, zone accessible
 * valeur de AgentID : Agent
 */
var carte [20][20]string = [20][20]string{
	{"X", "X", "X", "X", "X", "X", "X", "X", "W", "W", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "_", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "_", "X", "X"},
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

func (sim *Simulation) Env() *Environment {
	return &sim.env
}

func NewSimulation(agentCount int, maxStep int, maxDuration time.Duration) (simu *Simulation) {
	simu = &Simulation{}
	simu.agents = make([]Agent, 0, agentCount)
	simu.maxStep = maxStep
	simu.maxDuration = maxDuration

	// Communication entre agents
	mapChan := make(map[AgentID]chan Request)
	simu.env = *NewEnvironment([]Agent{}, carte, mapChan)
	//simu.env = *NewEnvironment([]Agent{}, playground, mapChan)

	// création des agents et des channels
	for i := 0; i < agentCount; i++ {
		// création de l'agent
		
		syncChan := make(chan int)
		//ag := NewAgent(id, &simu.env, syncChan, time.Duration(time.Second), 0, true, Coord{0, 8 + i%2}, Coord{0, 8 + i%2}, &UsagerLambda{}, Coord{0, 8 + i%2}, Coord{12 - 4*(i%2), 18 - 15*(i%2)})

		//ag := NewAgent(id, &simu.env, syncChan, 1000, 0, true, &UsagerLambda{},  Coord{18, 4}, Coord{0, 8}, 2, 1)
		
		
		ag := &Agent{}
		
		if i%2==0{ //Type Agent
			id := fmt.Sprintf("Agent%d", i)
			//NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, behavior Behavior, departure, destination Coord, width, height int)
			ag = NewAgent(id, &simu.env, syncChan, 1000, 0, true, &UsagerLambda{}, Coord{18, 4}, Coord{0, 8}, 1, 1)
		}else{ // Type Controleur
			id := fmt.Sprintf("Controleur%d", i)
			//ag = NewAgent(id, &simu.env, syncChan, 1000, 0, true, &UsagerLambda{}, Coord{1, 8}, Coord{8, 5}, 1, 1)
			ag = NewAgent(id, &simu.env, syncChan, 1000, 0, true, &Controleur{}, Coord{18, 12}, Coord{18, 4}, 1, 1)
		}
		
		
		//ag := NewAgent(id, &simu.env, syncChan, 1000, 0, true, &UsagerLambda{}, Coord{19, 12}, Coord{0, 8}, 2, 1)

		// ajout de l'agent à la simulation
		simu.agents = append(simu.agents, *ag)

		// ajout du channel de synchro
		simu.syncChans.Store(ag.ID(), syncChan)

		// ajout de l'agent à l'environnement
		ag.env.AddAgent(*ag)

		// ajout du channel de l'agent à l'environnement
		simu.env.agentsChan[ag.id] = make(chan Request)
	}



	return simu
}

func (simu *Simulation) Run() {
	// A REVOIR si nécessaire de faire appeler simu.env.pi() 
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
				c, _ := simu.syncChans.Load(agt.ID()) // communiquer les steps aux agents
				c.(chan int) <- step             // /!\ utilisation d'un "Type Assertion"
				time.Sleep(1 * time.Millisecond) // "cool down"
				<-c.(chan int)
			}
		}(agt)
	}

	time.Sleep(simu.maxDuration)

	log.Printf("Fin de la simulation [step: %d, in: %d, out: %d, π: %f]", simu.step, simu.env.PI())
}


func (simu *Simulation) Print_v0() {
	for {
		for i := 0; i < 20; i++ {
			fmt.Println(simu.env.station[i])
		}
		fmt.Println()
		fmt.Println()
		//fmt.Println("============================================================")
		//time.Sleep(time.Second / 4) // 60 fps !
		time.Sleep(500 * time.Millisecond) // 1 fps !
		//fmt.Print("\033[H\033[2J") // effacement du terminal
	}
}
func (simu *Simulation) Print() {
    for {
        for i := 0; i < 20; i++ {
            for j := 0; j < 20; j++ {
                element := simu.env.station[i][j]
                if len(element) > 1 {
                    fmt.Print(string(element[0]) + " ") // Afficher le premier caractère si la longueur est supérieure à 1
                } else {
                    fmt.Print(element+" ")
                }
            }
            fmt.Println()
        }
        fmt.Println()
        //fmt.Println("============================================================")
        //time.Sleep(time.Second / 4) // 60 fps !
        time.Sleep(500 * time.Millisecond) // 1 fps !
        //fmt.Print("\033[H\033[2J") // effacement du terminal
    }
}


func (simu *Simulation) Log() {
	// Not implemented
}
