package simulation

/*
	L'agent à Mobilité Reduite cherche la porte du metro la plus proche de lui
*/

import (
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"sort"
	"time"
)

type MobiliteReduite struct {
	req *req.Request
	//once sync.Once
}

func (mr *MobiliteReduite) Percept(ag *Agent) {
	//mr.once.Do(func(){mr.setUpDestination(ag)}) // on initialise la destination la plus proche, la fonction setUp est executé à la premiere appel à la fonction Percept()
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		//fmt.Printf("Requete recue par l'agent mR : %d \n", ag.request.Decision())
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
	if mr.req != nil {
		switch mr.req.Decision() {
		case Expel: // sinon alors la requete est de type "Viré" cette condition est inutile car MR ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
			//fmt.Println("[AgentLambda, Deliberate] Expel")
			ag.decision = Expel
			return
		case Disappear:
			//TODELETEfmt.Println("[Deliberate]", ag.id, "Disappear cond 1 (requete)")
			ag.decision = Disappear
			return
		case Wait:
			ag.decision = Wait
			return
		case EnterMetro:
			//TODELETEfmt.Println("[MobiliteReduite, Deliberate] EnterMetro")
			ag.decision = EnterMetro
			return
		case YouHaveToMove:
			//fmt.Println("J'essaye de bouger")
			movement := ag.MoveAgent()
			//fmt.Printf("Je suis agent %s Resultat du mouvement de la personne %t \n", ag.id, movement)
			if movement {
				ag.decision = Done
			} else {
				ag.decision = Noop
			}
			return
		default:
			ag.decision = Move
			return
		}
	} else if (ag.position != ag.departure && ag.position == ag.destination) && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si l'agent est arrivé à sa destination et qu'il est sur une sortie
		//fmt.Println("[Deliberate]",ag.id, "Disappear cond 2")
		ag.decision = Disappear
		/*}else if (ag.position != ag.departure && ag.position == ag.destination){
			// si l'agent est arrivé à la porte mais n'a pas reçu une requete du metro pour entrer, il attend
			ag.decision = Wait [A REVOIR]
		}*/
	} else if ag.stuck { // si l'agent est bloqué
		ag.decision = Wait
	} else {
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
		ag.env.RemoveAgent(ag)

	case Expel:
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		//fmt.Println("[AgentLambda, Act] destination = ",ag.destination)
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		//mr.MoveMR(ag)
		ag.MoveAgent()
	case EnterMetro:
		//TODELETEfmt.Printf("[MobiliteReduite, Act %s] EnterMetro \n", ag.id)
		ag.env.RemoveAgent(ag)
		mr.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], ACK)
	}
	ag.request = nil
}

/*
* Fonction qui permet de définir la destination d'un agent à mobilité réduite
 */
func (mr *MobiliteReduite) SetUpDestination(ag *Agent) {
	choix_voie := rand.Intn(len(ag.env.metros)) // choix de la voie de métro aléatoire
	dest_porte := (mr.findNearestGates(ag, ag.env.metros[choix_voie].way.gates))
	//fmt.Println("[MobiliteReduite, setUpDestination] dest_porte = ",dest_porte)
	ag.destination = dest_porte[0].Position
}

func (mr *MobiliteReduite) findNearestGates(ag *Agent, gates []alg.Coord) []Gate {
	var gateDistances []Gate
	// Calcul de la distance pour chaque porte
	for _, gate := range gates {
		dist := alg.Abs(ag.position[0]-gate[0]) + alg.Abs(ag.position[1]-gate[1])
		gateDistances = append(gateDistances, Gate{Position: gate, Distance: float64(dist)})
	}

	// Tri des Coords par distance
	sort.Slice(gateDistances, func(i, j int) bool {
		return gateDistances[i].Distance < gateDistances[j].Distance
	})

	return gateDistances
}
