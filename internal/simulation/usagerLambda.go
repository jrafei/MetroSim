package simulation

import (
	"fmt"

	"math/rand"
	alg "metrosim/internal/algorithms"
	"time"
	"sync"
)

type UsagerLambda struct {
	req *Request // requete recue par l'agent lambda
	once sync.Once
}

func (ul *UsagerLambda) Percept(ag *Agent) {
	ul.once.Do(func(){ul.setUpDestination(ag)}) // la fonction setUp est executé à la premiere appel de la fonction Percept()
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		//print("Requete recue par l'agent lambda : ", ag.request.decision, "\n")
		ul.req = ag.request
	default:
		ag.stuck = ag.isStuck()
		if ag.stuck {
			return

		}
	}
}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	//fmt.Println("[AgentLambda Deliberate] decision :", ul.req.decision)
	if (ul.req != nil ) {
		if ul.req.decision == Stop{
			ag.decision = Wait
			ul.req = nil //demande traitée
			return
		} else if ul.req.decision == Expel { // cette condition est inutile car l'usager lambda ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
				//fmt.Println("[AgentLambda, Deliberate] Expel")
				ag.decision = Expel
				ul.req = nil //demande traitée
				return
		}else if ul.req.decision == Disappear {
			ag.decision = Disappear
			return
		}else if ul.req.decision == Wait {
			ag.decision = Wait
		}else if ul.req.decision == EnterMetro {
			ag.decision = EnterMetro
			}
	}else if (ag.position != ag.departure && ag.position == ag.destination) && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si l'agent est arrivé à sa destination et qu'il est sur une sortie
			//fmt.Println(ag.id, "disappear")
			ag.decision = Disappear
		} else if ag.stuck{ // si l'agent est bloqué
			ag.decision = Wait
				}else {
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
	} else if ag.decision == EnterMetro {
		fmt.Println("[UsagerLambda, Act] EnterMetro")
		RemoveAgent(&ag.env.station, ag)
		ul.req.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], ACK)
	} else if ag.decision == Expel {
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		//fmt.Println("[AgentLambda, Act] destination = ",ag.destination)
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()
	} else {
		// nothing to do
	}
}


func (ul *UsagerLambda)setUpDestination(ag *Agent){
	choix_voie := rand.Intn(2) // choix de la voie de métro aléatoire
	dest_porte := (ag.findBestGate(ag.env.metros[choix_voie].way.gates))
	ag.destination = dest_porte
}