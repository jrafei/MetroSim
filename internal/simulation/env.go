package simulation

import (
	"fmt"
	alg "metrosim/internal/algorithms"
	"sync"
)

//TODO:rajouter les entrées et sorties

type Environment struct {
	sync.RWMutex
	ags              []Agent
	agentCount       int
	station          [50][50]string
	agentsChan       map[AgentID]chan Request
	controlledAgents map[AgentID]bool
	newAgentChan     chan Agent
	metros           []Metro
}

func NewEnvironment(ags []Agent, carte [50][50]string, newAgtCh chan Agent, agtCount int) (env *Environment) {
	agentsCh := make(map[AgentID]chan Request)
	mapControlled := make(map[AgentID]bool)
	for _, ag := range ags {
		mapControlled[ag.id] = false

	}
	return &Environment{ags: ags, agentCount: agtCount, station: carte, agentsChan: agentsCh, controlledAgents: mapControlled, newAgentChan: newAgtCh}
}

func (env *Environment) AddAgent(agt Agent) {
	env.Lock()
	defer env.Unlock()
	env.ags = append(env.ags, agt)
	env.controlledAgents[agt.id] = false
	// ajout du channel de l'agent à l'environnement
	env.agentsChan[agt.id] = make(chan Request, 5)
	env.agentCount++
	env.newAgentChan <- agt
}

func (env *Environment) RemoveAgent(agt Agent) {
	// TODO:gérer la suppression dans simu
	for i := 0; i < len(env.station); i++ {
		if env.ags[i].id == agt.id {
			// Utiliser la syntaxe de découpage pour supprimer l'élément
			env.ags = append(env.ags[:i], env.ags[i+1:]...)
			delete(env.agentsChan, agt.id)
			// Sortir de la boucle après avoir trouvé et supprimé l'élément
			break
		}
	}
	//env.agentCount--
}

func (env *Environment) Do(a Action, c alg.Coord) (err error) {
	env.Lock()
	defer env.Unlock()

	switch a {
	// case Mark:
	// 	if c[0] < 0 || c[0] > 1 || c[1] < 0 || c[1] > 1 {
	// 		return fmt.Errorf("bad coordinates (%f,%f)", c[0], c[1])
	// 	}

	// 	return nil

	case Noop:
		return nil
	}

	return fmt.Errorf("bad action number %d", a)
}

func (env *Environment) PI() float64 {
	env.RLock()
	defer env.RUnlock()

	return 4
}

func (env *Environment) Rect() alg.Coord {
	return alg.Coord{0, 0}
}

func (env *Environment) GetAgentChan(agt_id AgentID) chan Request {
	return env.agentsChan[agt_id]
}

func existAgent(c string) bool {
	return c != "X" && c != "E" && c != "S" && c != "W" && c != "Q" && c != "_" && c != "B"
}

func calculDirection(depart alg.Coord, arrive alg.Coord) int {
	if depart[0] == arrive[0] {
		if depart[1] > arrive[1] {
			return 3 //Gauche
		} else {
			return 1 //droite
		}
	} else {
		if depart[0] > arrive[0] {
			return 0 //haut
		} else {
			return 2 //bas
		}
	}
}
