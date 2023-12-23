package simulation

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
	alg "metrosim/internal/algorithms"
)


type MobiliteReduite struct {
	req *Request
	once sync.Once
}


func (mr *MobiliteReduite) Percept(ag *Agent) {
	mr.once.Do(func(){mr.setUpDestination(ag)}) // la fonction setUp est executé à la premiere appel de la fonction Percept()
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		fmt.Printf("Requete recue par l'agent mR : %d \n", ag.request.decision)
		mr.req = ag.request
	default:
		ag.stuck = ag.isStuck()
		if ag.stuck {
			return
		}
	}
}

func (mr *MobiliteReduite) Deliberate(ag *Agent) {
	//fmt.Println("[AgentLambda Deliberate] decision :", ul.req.decision)
	if (mr.req != nil ) {
		if mr.req.decision == Stop{
			ag.decision = Wait
			mr.req = nil //demande traitée
			return
		} else if mr.req.decision == Expel { // sinon alors la requete est de type "Viré" cette condition est inutile car MR ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
				//fmt.Println("[AgentLambda, Deliberate] Expel")
				ag.decision = Expel
				mr.req = nil //demande traitée
				return
			}else if mr.req.decision == Disappear {
				fmt.Println("[Deliberate]",ag.id, "Disappear cond 1 (requete)")
				ag.decision = Disappear
				mr.req = nil
				return
			}else if mr.req.decision == Wait {
					ag.decision = Wait
					mr.req = nil
					return
			}else if mr.req.decision == EnterMetro {
					fmt.Println("[MobiliteReduite, Deliberate] EnterMetro")
					ag.decision = EnterMetro
					mr.req = nil
					return
			}
	}else if (ag.position != ag.departure && ag.position == ag.destination) && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si l'agent est arrivé à sa destination et qu'il est sur une sortie
			fmt.Println("[Deliberate]",ag.id, "Disappear cond 2")
			ag.decision = Disappear
		}else if (ag.position != ag.departure && ag.position == ag.destination){
			// si l'agent est arrivé à la porte mais n'a pas reçu une requete du metro pour entrer, il attend
			ag.decision = Wait
		} else if ag.stuck{ // si l'agent est bloqué
			ag.decision = Wait
			}else {
			ag.decision = Move
			}	
}

func (mr *MobiliteReduite) Act(ag *Agent) {
	//fmt.Println("[AgentLambda Act] decision :",ag.decision)
	switch ag.decision {
	case Move:
		//mr.MoveMR(ag)
		ag.MoveAgent()
	case Wait:
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	case Disappear:
		RemoveAgent(&ag.env.station, ag)
		
	case Expel : 
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		//fmt.Println("[AgentLambda, Act] destination = ",ag.destination)
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node,0)
		//mr.MoveMR(ag)
		ag.MoveAgent()
	case EnterMetro :
		fmt.Printf("[MobiliteReduite, Act %s] EnterMetro \n", ag.id)
		RemoveAgent(&ag.env.station, ag)
		mr.req.demandeur <- *NewRequest(ag.env.agentsChan[ag.id], ACK)
	}
}

/*
* Fonction qui permet de définir la destination d'un agent à mobilité réduite
*/
func (mr *MobiliteReduite)setUpDestination(ag *Agent){
	choix_voie := rand.Intn(2) // choix de la voie de métro aléatoire
	dest_porte := (ag.findNearestGates(ag.env.metros[choix_voie].way.gates))
	//fmt.Println("[MobiliteReduite, setUpDestination] dest_porte = ",dest_porte)
	ag.destination = dest_porte[0].Position
}
