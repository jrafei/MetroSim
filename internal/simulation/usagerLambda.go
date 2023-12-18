package simulation

import (
	"fmt"
	"math/rand"
	"time"
	alg "metrosim/internal/algorithms"
)

type UsagerLambda struct {
	req *Request
}

func (ul *UsagerLambda) Percept(ag *Agent) {
	// récupérer le channel de l'agent lambda
	//fmt.Println("[AgentLambda, Percept] direction ", ag.direction)

	chan_agt := ag.env.GetAgentChan(ag.id) 
	select {
	case req := <-chan_agt : //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		//fmt.Println("[AgentLambda, Percept] Requete recue par l'agent lambda : ", req.decision)
		ul.req = &req
	case <- time.After(time.Second):
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
		} else { // sinon alors la requete est de type "Viré" cette condition est inutile car l'usager lambda ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
				fmt.Println("[AgentLambda, Deliberate] Expel")
				ag.decision = Expel
				ul.req = nil //demande traitée
		}
	}else if ag.position == ag.destination && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") {
			//fmt.Println(ag.id, "disapear")
			ag.decision = Disapear
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
	} else if ag.decision == Disapear {
		RemoveAgent(&ag.env.station, ag)
	} else { //age.decision == Expel
		fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		fmt.Println("[AgentLambda, Act] destination = ",ag.destination)
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node,0)
		ag.MoveAgent()
	}
}
