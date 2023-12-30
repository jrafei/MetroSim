package simulation


/*
  Agent qui se dirige vers la porte la plus proche sans trop de monde (bon rapport monde/proximité )
*/

import (
	"fmt"

	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"time"
	"sync"
	"math"
)

type UsagerNormal struct {
	req *req.Request // requete recue par l'agent lambda
	once sync.Once
}

func (ul *UsagerNormal) Percept(ag *Agent) {
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

func (ul *UsagerNormal) Deliberate(ag *Agent) {
	//fmt.Println("[AgentLambda Deliberate] decision :", ul.req.decision)
	if (ul.req != nil ) {
		switch ul.req.Decision() {
			case Stop :
				ag.decision = Wait
				ul.req = nil //demande traitée
				return
			case Expel : // cette condition est inutile car l'usager lambda ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
				//fmt.Println("[AgentLambda, Deliberate] Expel")
				ag.decision = Expel
				ul.req = nil //demande traitée
				return
			case Disappear :
				ag.decision = Disappear
				return
			case Wait :
				ag.decision = Wait
				return
			case EnterMetro :
				ag.decision = EnterMetro
				return
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

func (ul *UsagerNormal) Act(ag *Agent) {
	//fmt.Println("[AgentLambda Act] decision :",ag.decision)
	if ag.decision == Move {
		ag.MoveAgent()
	} else if ag.decision == Wait {
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	} else if ag.decision == Disappear {
		ag.env.RemoveAgent(ag)
	} else if ag.decision == EnterMetro {
		fmt.Println("[UsagerNormal, Act] EnterMetro")
		ag.env.RemoveAgent(ag)
		ul.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], ACK)
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


func (ul *UsagerNormal)setUpDestination(ag *Agent){
	choix_voie := rand.Intn(2) // choix de la voie de métro aléatoire
	dest_porte := (ul.findBestGate(ag, ag.env.metros[choix_voie].way.gates))
	ag.destination = dest_porte
}



func (ul *UsagerNormal) findBestGate(ag *Agent, gates []alg.Coord) alg.Coord {
	gatesDistances := make([]Gate, len(gates))
	for i, gate := range gates {
		dist := alg.Abs(ag.position[0]-gate[0]) + alg.Abs(ag.position[1]-gate[1])
		nbAgents := float64(ag.env.getNbAgentsAround(gate))
		gatesDistances[i] = Gate{Position: gate, Distance: float64(dist), NbAgents: nbAgents}
	}
	fmt.Println("[findBestGate] gates non normalisé : ",gatesDistances)
	normalizedGates, _, _ := normalizeGates(gatesDistances)
	fmt.Println("[findBestGate] gates normalisé : ",normalizedGates)
	var bestGate Gate
	lowestScore := 2.0 // Puisque la somme des scores normalisés ne peut pas dépasser 2

	for _, gate := range normalizedGates {
		score := float64(gate.NbAgents) + gate.Distance
		if score < lowestScore {
			lowestScore = score
			bestGate = gate
		}
	}
	return bestGate.Position
}


// Normalise les valeurs d'un ensemble de portes
func normalizeGates(gates []Gate) ([]Gate, float64, float64) {
    var minAgents, maxAgents float64 = math.MaxFloat64, 0
    var minDistance, maxDistance float64 = math.MaxFloat64, 0

    // Trouver les valeurs max et min pour la normalisation
    for _, gate := range gates {
        if gate.NbAgents > maxAgents {
            maxAgents = gate.NbAgents
        }
        if gate.NbAgents < minAgents {
            minAgents = gate.NbAgents
        }
        if gate.Distance > maxDistance {
            maxDistance = gate.Distance
        }
        if gate.Distance < minDistance {
            minDistance = gate.Distance
        }
    }

    // Normaliser les valeurs
	d_agt := (maxAgents - minAgents) 
	if  d_agt == 0 {
		d_agt = 1.0
	}
	d_dist := (maxDistance - minDistance)
	if d_dist == 0 {
		d_dist = 1.0
	}
	fmt.Println("[normalizeGates] d_dist : ",d_dist)
    for i := range gates {
        gates[i].NbAgents = (gates[i].NbAgents - minAgents) / d_agt
		//fmt.Println("[normalizeGates] gates[i].Distance : ",gates[i].Distance)
		//fmt.Println("[normalizeGates] minDistance : ",minDistance)
		//fmt.Println("[normalizeGates] d_dist : ",d_dist)
        gates[i].Distance = (gates[i].Distance - minDistance) / d_dist
	}
    return gates, float64(maxAgents - minAgents), maxDistance - minDistance
}
