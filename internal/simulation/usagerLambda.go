package simulation

import (
	//"fmt"
	"fmt"
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
		print("Requete recue par l'agent lambda : ", ag.request.decision, "\n")
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
	} else if ul.req.decision == Disappear || (ag.position == ag.destination && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S")) {
		fmt.Println(ag.id, "disappear")
		ag.decision = Disappear
	} else if ul.req.decision == Wait {
		ag.decision = Wait
	} else if ul.req.decision == YouHaveToMove {
		movement := ag.MoveAgent()
		if movement {
			ag.decision = 5
		} else {
			ag.decision = 0
		}
	} else {
		ag.decision = Move
	}
}

func (ul *UsagerLambda) Act(ag *Agent) {
	//fmt.Println("[AgentLambda Act] decision :",ag.decision)
	switch ag.decision {
	case Move:
		ag.MoveAgent()
	case Wait: // temps d'attente aléatoire
		n := rand.Intn(2)
		time.Sleep(time.Duration(n) * time.Second)
	case Disappear:
		RemoveAgent(&ag.env.station, ag)
	case Noop:
		//Cas ou un usager impoli demande a un usager de bouger et il refuse
		ag.request.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], 0)
	case Done:
		//Cas ou un usager impoli demande a un usager de bouger et il le fait
		ag.request.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], 5)
	default:
		//age.decision == Expel
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.departure
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()
	}
}
