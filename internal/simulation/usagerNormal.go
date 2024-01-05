package simulation

/*
  Agent qui se dirige vers la porte la plus proche sans trop de monde (bon rapport monde/proximité )
  On normalise les valeurs de proximité et de monde pour avoir un score entre 0 et 1, on choisit la porte ayant le moins du score (le score est la somme des deux valeurs normalisées)
  Si plusieurs portes ont le même score, on choisit celle ayant le moins de monde
*/

import (
	"math"
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"sort"
	"time"
	//"sync"
)

type UsagerNormal struct {
	req *req.Request // req recue par l'agent lambda
	//once sync.Once
}

func (un *UsagerNormal) Percept(ag *Agent) {
	//un.once.Do(func() { un.SetUpDestination(ag) }) // la fonction setUp est executé à la premiere appel de la fonction Percept()
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		//print("req recue par l'agent lambda : ", ag.request.decision, "\n")
		un.req = ag.request
	default:
		ag.stuck = ag.isStuck()
		if ag.stuck {
			return

		}
	}
}

func (un *UsagerNormal) Deliberate(ag *Agent) {
	//fmt.Println("[AgentLambda Deliberate] decision :", un.req.decision)
	if un.req != nil {
		switch un.req.Decision() {
		case Expel: // cette condition est inutile car l'usager lambda ne peut pas etre expunsé , elle est nécessaire pour les agents fraudeurs
			//fmt.Println("[AgentLambda, Deliberate] Expel")
			ag.decision = Expel
			return
		case Disappear:
			ag.decision = Disappear
			return
		case EnterMetro:
			ag.decision = EnterMetro
			return
		case Wait:
			ag.decision = Wait
			return
		case Move:
			ag.decision = Move
			return
		case YouHaveToMove:
			//TODELETEfmt.Println("[AgentNormal, Deliberate] J'essaye de bouger ", ag.id)
			movement := ag.MoveAgent()
			//fmt.Printf("Je suis agent %s Resultat du mouvement de la personne %t \n", ag.id, movement)
			if movement {
				//TODELETEfmt.Println("[AgentNormal, Deliberate] J'ai bougé ", ag.id)
				ag.decision = Done
			} else {
				ag.decision = Noop
			}
			return
		}
	} else if (ag.position != ag.departure && ag.position == ag.destination) && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si l'agent est arrivé à sa destination et qu'il est sur une sortie
		//fmt.Println(ag.id, "disappear")
		ag.decision = Disappear
	} else if ag.stuck { // si l'agent est bloqué
		ag.decision = Wait
	} else {
		ag.decision = Move
		//un.setUpDestination(ag)
	}
}

func (un *UsagerNormal) Act(ag *Agent) {
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
		//TODELETEfmt.Printf("[UsagerNormal, Act] agent %s entre dans le Metro \n", ag.id)
		ag.env.RemoveAgent(ag)
		//fmt.Printf("Demandeur d'entrer le metro : %s \n",un.req.Demandeur())
		un.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], ACK)
	case Expel:
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		//TODELETEfmt.Printf("[UsagerNormal, Act] destination de l'agent %s = %s \n", ag.id, ag.destination)
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node, 0)
		ag.MoveAgent()

	case Noop:
		//Cas ou un usager impoli demande a un usager de bouger et il refuse
		un.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Noop)
		// nothing to do
	case Done:
		//Cas ou un usager impoli demande a un usager de bouger et il le fait
		un.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Done)
	case TryToMove:
		movement := ag.MoveAgent()
		//TODELETEfmt.Printf("Je suis %s est-ce que j'ai bougé? %t \n", ag.id, movement)
		if movement {
			un.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Done)
		} else {
			un.req.Demandeur() <- *req.NewRequest(ag.env.agentsChan[ag.id], Noop)
		}
	}
	ag.request = nil
}

func (un *UsagerNormal) SetUpDestination(ag *Agent) {
	//t := rand.Intn(10) +1
	//time.Sleep(time.Duration(t) * time.Second) // "cool down"
	//fmt.Println("[UsagerNormal, setUpDestination] setUpDestination")
	choix_voie := rand.Intn(len(ag.env.metros)) // choix de la voie de métro aléatoire
	dest_porte := (un.findBestGate(ag, ag.env.metros[choix_voie].way.gates))
	ag.destination = dest_porte
	//TODELETEfmt.Println("[UsagerNormal, setUpDestination] destination de l'agent ", ag.id, " = ", ag.destination, " son position = ", ag.position)
}

func (un *UsagerNormal) findBestGate(ag *Agent, gates []alg.Coord) alg.Coord {

	uniquegates := make([]alg.Coord, 0)
	for i, gate := range gates {
		if i+1 < len(gates) && twocloseGate(gate, gates[i+1]) {
			// Si la porte est trop proche d'une autre, on l'ignore, on considère que c'est la même
			continue
		}
		uniquegates = append(uniquegates, gate)
		//gatesDistances[i] = Gate{Position: gate, Distance: float64(dist), NbAgents: nbAgents}
	}

	gatesDistances := make([]Gate, len(uniquegates))
	for i, gate := range uniquegates {
		dist := alg.Abs(ag.position[0]-gate[0]) + alg.Abs(ag.position[1]-gate[1]) // Distance de Manhattan entre l'agent et la porte
		nbAgents := float64(ag.env.getNbAgentsAround(gate))
		gatesDistances[i] = Gate{Position: gate, Distance: float64(dist), NbAgents: nbAgents}
	}

	//fmt.Println("[findBestGate] agent : ",ag.id)
	//fmt.Println("[findBestGate] agent Position : ",ag.position)
	//fmt.Println("[findBestGate] gates non normalisé : ",gatesDistances)
	normalizedGates, _, _ := normalizeGates(gatesDistances)
	//TODELETEfmt.Println("[findBestGate, %s] gates normalisé : ",ag.id ,normalizedGates)

	bestGates := gates_with_lowest_score(normalizedGates)
	bestGate := bestGates[0]
	if len(bestGates) > 1 {
		//on choisit la porte ayant le moins de monde
		for _, gate := range bestGates {
			if gate.NbAgents < bestGate.NbAgents {
				bestGate = gate
			}
		}
	}
	//fmt.Println("[findBestGate] bestGate : ",bestGate)

	//choix de la porte aléatoire parmi les portes adjacentes
	nearGates := ag.env.getNearGateFromGate(bestGate)
	var bestGatePos alg.Coord
	if len(nearGates) > 1 {
		bestGatePos = nearGates[rand.Intn(len(nearGates))]
	} else {
		bestGatePos = nearGates[0]
	}
	return bestGatePos
}

func twocloseGate(gate1 alg.Coord, gate2 alg.Coord) bool {
	if (gate1[0] == gate2[0]) && (gate1[1] == gate2[1]+1 || gate1[1] == gate2[1]-1) {
		return true
	}
	if (gate1[1] == gate2[1]) && (gate1[0] == gate2[0]+1 || gate1[0] == gate2[0]-1) {
		return true
	}
	return false
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
	if d_agt == 0 {
		d_agt = 1.0
	}
	d_dist := (maxDistance - minDistance)
	if d_dist == 0 {
		d_dist = 1.0
	}
	//fmt.Println("[normalizeGates] d_dist : ",d_dist)
	for i := range gates {
		gates[i].NbAgents = (gates[i].NbAgents - minAgents) / d_agt
		//fmt.Println("[normalizeGates] gates[i].Distance : ",gates[i].Distance)
		//fmt.Println("[normalizeGates] minDistance : ",minDistance)
		//fmt.Println("[normalizeGates] d_dist : ",d_dist)
		gates[i].Distance = (gates[i].Distance - minDistance) / d_dist
	}
	return gates, float64(maxAgents - minAgents), maxDistance - minDistance
}

// Calcul du score d'un Gate
func (g Gate) Score() float64 {
	return g.Distance + g.NbAgents
}

// sort_by_score trie une tranche de Gates par ordre croissant de leur score
func sort_by_score(gates []Gate) {
	sort.Slice(gates, func(i, j int) bool {
		return gates[i].Score() < gates[j].Score() // Ordre croissant
	})
}

// gates_with_highest_score renvoie une tranche de Gates ayant le score le moins élevé
func gates_with_lowest_score(gates []Gate) []Gate {
	if len(gates) == 0 {
		return nil
	}

	sort_by_score(gates) // D'abord, on trie les gates

	lowestScore := gates[0].Score() // Le premier gate a le score le moins élevé après le tri
	var lowestScoreGates []Gate

	for _, gate := range gates {
		if gate.Score() == lowestScore {
			lowestScoreGates = append(lowestScoreGates, gate)
		} else {
			break // Puisque les gates sont triés, pas besoin de vérifier plus loin
		}
	}

	return lowestScoreGates
}
