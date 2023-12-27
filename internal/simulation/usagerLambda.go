package simulation

import (
	//"fmt"

	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"time"
)

type UsagerLambda struct {
	req req.Request
}

func (ul *UsagerLambda) Percept(ag *Agent) {
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		print("Requete recue par %d : %d", ag.id, ag.request.decision, "\n")
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

	if ul.req.Decision() == Stop {
		ag.decision = Wait
	} else if ul.req.Decision() == Expel { // cette condition est inutile car l'usager lambda ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
		ag.decision = Expel
	} else if ul.req.Decision() == Disappear || (ag.position != ag.departure && ag.position == ag.destination) && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") {
		ag.decision = Disappear
	} else if ul.req.Decision() == EnterMetro {
		ag.decision = EnterMetro
	} else if ul.req.Decision() == Wait {
		ag.decision = Wait
	} else if ul.req.decision == YouHaveToMove {
		fmt.Println("J'essaye de bouger")
		movement := ag.MoveAgent()
		fmt.Printf("Je suis agent %s Resultat du mouvement de la personne %t \n", ag.id, movement)
		if movement {
			ag.decision = Done
		} else {
			ag.decision = Noop
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
		ag.env.RemoveAgent(ag)
	case EnterMetro:
		ag.env.RemoveAgent(ag)
		ul.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], ACK)
	case Expel : 
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.departure
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()
	case Noop :
		//Cas ou un usager impoli demande a un usager de bouger et il refuse
		ag.request.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], 0)
		// nothing to do
	case Done : 
		//Cas ou un usager impoli demande a un usager de bouger et il le fait
		ag.request.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], 5)
	case TryToMove :
		movement := ag.MoveAgent()
		fmt.Printf("Je suis %s est-ce que j'ai bougé? %t \n", ag.id, movement)
		if movement {
			ag.request.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], Done)
		} else {
			ag.request.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], Noop)
		}
	default:
		//age.decision == Expel
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.departure
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()
	}
}
