package simulation

import (
	//"fmt"

	"math/rand"
	alg "metrosim/internal/algorithms"
	"time"
)

type UsagerLambda struct {
	req Request
}

func (ul *UsagerLambda) Percept(ag *Agent) {
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		ul.req = *ag.request
	default:
		ag.stuck = ag.isStuck()
		if ag.stuck {
			return

		}
	}

}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	//fmt.Println("[AgentLambda Deliberate] decision :", ul.req.decision)

	if ul.req.decision == Stop {
		ag.decision = Wait
	} else if ul.req.decision == Expel { // cette condition est inutile car l'usager lambda ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
		ag.decision = Expel
	} else if ul.req.decision == Disappear || (ag.position != ag.departure && ag.position == ag.destination) && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") {
		ag.decision = Disappear
	} else if ul.req.decision == Wait {
		ag.decision = Wait
	} else {
		ag.decision = Move
	}
}

func (ul *UsagerLambda) Act(ag *Agent) {
	//fmt.Println("[AgentLambda Act] decision :",ag.decision)
	if ag.decision == Move {
		ag.MoveAgent()
	} else if ag.decision == Wait {
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	} else if ag.decision == Disappear {
		RemoveAgent(&ag.env.station, ag)
	} else if ag.decision == Expel {
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.departure
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()
	} else {
		// nothing to wait
	}
}
