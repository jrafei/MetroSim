package simulation

import (
	"fmt"
	"sync"
)

type Environment struct {
	sync.RWMutex
	ags        []Agent
	agentCount int
	station    [20][20]string
	agentsChan map[AgentID]chan Request
}

func NewEnvironment(ags []Agent, carte [20][20]string, agentsCh map[AgentID]chan Request) (env *Environment) {
	return &Environment{ags: ags, agentCount: len(ags), station: carte, agentsChan: agentsCh}
}

func (env *Environment) AddAgent(agt Agent) {
	env.ags = append(env.ags, agt)
	env.agentCount++
}

func (env *Environment) RemoveAgent(agt Agent) {
	for i := 0; i < len(env.station); i++ {
		if env.ags[i].id == agt.id {
			// Utiliser la syntaxe de découpage pour supprimer l'élément
			env.ags = append(env.ags[:i], env.ags[i+1:]...)
			delete(env.agentsChan,agt.id)
			// Sortir de la boucle après avoir trouvé et supprimé l'élément
			break
		}
	}
	env.agentCount--
}

func (env *Environment) Do(a Action, c Coord) (err error) {
	env.Lock()
	defer env.Unlock()

	switch a {
	case Mark:
		if c[0] < 0 || c[0] > 1 || c[1] < 0 || c[1] > 1 {
			return fmt.Errorf("bad coordinates (%f,%f)", c[0], c[1])
		}

		return nil

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

func (env *Environment) Rect() Coord {
	return Coord{0, 0}
}

func (env *Environment) GetAgentChan(agt_id AgentID) chan Request {
	return env.agentsChan[agt_id]
}
