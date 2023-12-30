package simulation

import (
	"fmt"

	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"time"
	"sync"
)

type UsagerLambda struct {
	requete *req.Request
	once sync.Once
}

func (ul *UsagerLambda) Percept(ag *Agent) {
	ul.once.Do(func(){ul.setUpAleaDestination(ag)}) // la fonction setUp est executé à la premiere appel de la fonction Percept()
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		//print("Requete recue par l'agent lambda : ", ag.request.decision, "\n")
		ul.requete = ag.request
		ag.request = nil // la requete est traitée
	default:
		ag.stuck = ag.isStuck()
		if ag.stuck {
			return
		}
	}
}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	//fmt.Println("[AgentLambda Deliberate] decision :", ul.req.decision)

	if ul.requete != nil {
		switch ul.requete.Decision() {
		case Stop :
		ag.decision = Stop
		case Expel: // cette condition est inutile car l'usager lambda ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
		ag.decision = Expel
		case Disappear :
		ag.decision = Disappear
		case EnterMetro :
		fmt.Println("[UsagerLambda, Deliberate] EnterMetro")
		ag.decision = EnterMetro
		case Wait :
		ag.decision = Wait
		case Move :
		ag.decision = Move
		case YouHaveToMove :
		//fmt.Println("J'essaye de bouger")
		movement := ag.MoveAgent()  
		//fmt.Printf("Je suis agent %s Resultat du mouvement de la personne %t \n", ag.id, movement)
		if movement {
			ag.decision = Done
		} else {
			ag.decision = Noop
		}
		default :
		ag.decision = Move

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
	switch ag.decision {
	case Stop : 
		time.Sleep(time.Duration(5) * time.Second) 
	case Move:
		ag.MoveAgent()
	case Wait: // temps d'attente aléatoire
		n := rand.Intn(2)
		time.Sleep(time.Duration(n) * time.Second)
	case Disappear:
		ag.env.RemoveAgent(ag)
	case EnterMetro :
		fmt.Printf("[UsagerLambda, Act] agent %s entre dans le Metro \n",ag.id)
		ag.env.RemoveAgent(ag)
		fmt.Printf("Demandeur d'entrer le metro : %s \n",ul.requete.Demandeur())
		ul.requete.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], ACK)
	case Expel :
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		fmt.Printf("[AgentLambda, Act] destination de l'agent %s = %s \n",ag.id,ag.destination)
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()

	case Noop :
		//Cas ou un usager impoli demande a un usager de bouger et il refuse
		ag.request.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Noop)
		// nothing to do
	case Done : 
		//Cas ou un usager impoli demande a un usager de bouger et il le fait
		ag.request.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Done)
	case TryToMove :
		movement := ag.MoveAgent()
		fmt.Printf("Je suis %s est-ce que j'ai bougé? %t \n", ag.id, movement)
		if movement {
			ag.request.Demandeur()<- *req.NewRequest(ag.env.agentsChan[ag.id], Done)
		} else {
			ag.request.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Noop)
		}
	}
	ul.requete = nil //demande traitée
}


func (ul *UsagerLambda)setUpAleaDestination(ag *Agent){
	choix_voie := rand.Intn(2) // choix de la voie de métro aléatoire
	dest_porte := rand.Intn(len(ag.env.metros[choix_voie].way.gates)) // choix de la porte de métro aléatoire
	ag.destination = ag.env.metros[choix_voie].way.gates[dest_porte]
}