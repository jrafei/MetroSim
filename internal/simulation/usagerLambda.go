package simulation

import (
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"time"
)

type UsagerLambda struct {
	requete *req.Request
	// once    sync.Once
}

func (ul *UsagerLambda) Percept(ag *Agent) {
	//ul.once.Do(func() { ul.setUpAleaDestination(ag) }) // la fonction setUp est executé à la premiere appel de la fonction Percept()
	/*fmt.Println("[UsagerLambda, Percept] ", ag.id)
	if ( ag.env.controlledAgents[ag.id]){
		fmt.Println("[UsagerLambda, Percept] ", ag.id, " est controlé")
	}
	*/
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		//TODELETEfmt.Printf("Requete recue par l'agent lambda %s : %d \n ", ag.id, ag.request.Decision())
		ul.requete = ag.request

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
		case Expel:
			ag.decision = Expel
			return
		case Disappear:
			ag.decision = Disappear
			return
		case EnterMetro:
			//TODELETEfmt.Println("[UsagerLambda, Deliberate] EnterMetro %s", ag.id)
			ag.decision = EnterMetro
			return
		case Wait:
			ag.decision = Wait
			return
		case Move:
			ag.decision = Move
			return
		case YouHaveToMove:
			ag.decision = TryToMove
			return
		default:
			ag.decision = Move
			return
		}
	} else if (ag.position != ag.departure && ag.position == ag.destination) && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si l'agent est arrivé à sa destination et qu'il est sur une sortie
		//fmt.Println(ag.id, "disappear")
		ag.decision = Disappear
	} else if ag.stuck { // si l'agent est bloqué
		ag.decision = Wait
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
		//fmt.Printf("[UsagerLambda, Act] agent %s est disparu \n",ag.id)
		ag.env.RemoveAgent(ag)
	case EnterMetro:
		//fmt.Printf("[UsagerLambda, Act] agent %s entre dans le Metro \n",ag.id)
		ag.env.RemoveAgent(ag)
		//fmt.Printf("Demandeur d'entrer le metro : %s \n",ul.requete.Demandeur())
		ul.requete.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], ACK)
	case Expel:
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		//TODELETEfmt.Printf("[UsagerLambda, Act] position de l'agent %s =  (%d , %d) \n", ag.id, ag.position[0], ag.position[1])
		//TODELETEfmt.Printf("[UsagerLambda, Act] destination de l'agent %s = (%d , %d) \n", ag.id, ag.destination[0], ag.destination[1])
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()
		//TODELETEfmt.Printf("[UsagerLambda, Act] J'ai bougé %s , ma position = (%d , %d)\n", ag.id, ag.position[0], ag.position[1])
	case TryToMove:
		//fmt.Printf("Je suis %s est-ce que j'ai bougé? %t \n", ag.id, movement)
		if ag.ShiftAgent() {
			ul.requete.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Done)
		} else {
			ul.requete.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Noop)
		}
	}

	//ag.request = nil // la requete est traitée
	if ag.request != nil && ag.decision == ag.request.Decision() {
		ag.request = nil
		ul.requete = nil
	} // la requete est traitée
}

func (ul *UsagerLambda) SetUpDestination(ag *Agent) {
	//fmt.Println("[UsagerLambda, setUpAleaDestination] setUpAleaDestination")
	choix_voie := rand.Intn(len(ag.env.metros))                       // choix de la voie de métro aléatoire
	dest_porte := rand.Intn(len(ag.env.metros[choix_voie].way.gates)) // choix de la porte de métro aléatoire
	ag.destination = ag.env.metros[choix_voie].way.gates[dest_porte]
}

func isControlledAgt(ag *Agent) bool {
	return ag.env.controlledAgents[ag.id]
}
