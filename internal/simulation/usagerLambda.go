package simulation

import (
	"fmt"
	"math/rand"
	"time"
)

type UsagerLambda struct{
	req Request
}
/*
func (ul *UsagerLambda) Percept(ag *Agent) {
	// récupérer le channel de l'agent lambda

	chan_agt := ag.env.GetAgentChan(ag.id) 
	select {
	case req := <-chan_agt : //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		ul.req = req
	case <- time.After(time.Second):
		ag.stuck = ag.isStuck()
	}
	//
	ag.stuck = ag.isStuck()
	if ag.stuck {
		return

	}
	//
}


func (ul *UsagerLambda) Deliberate(ag *Agent) {
	if ag.position == ag.destination && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") {
		fmt.Println(ag.id, "disapear")
		ag.decision = Disapear
	} else if ul.req.decision == Wait{
		ag.decision = Wait
	} else if ul.req.decision == Expel{ // cette condition est inutile car l'usager lambda ne peut pas etre expulsé
		ag.decision = Expel
	}
}

func (ul *UsagerLambda) Act(ag *Agent) {
	if ag.decision == Move {
		ag.MoveAgent()
	} else if ag.decision == Wait {
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	} else if ag.decision == Disapear {
		RemoveAgent(&ag.env.station, ag)
	} else {
		ag.destination = ag.departure
		ag.MoveAgent()
	}
}


/*
 * Fonction qui permet de trouver la sortie la plus proche

func (ul *UsagerLambda) findNearestExit(env *Environment) Coord{
	station := env.station
	for i := 0; i < len(station); i++ {
		for j := 0; j < len(station[i]); j++ {
			if station[i][j] == "X" {
				return Coord{i,j}
			}
		}
	}
	return Coord{0,0}
}
*/


func (ul *UsagerLambda) Percept(ag *Agent) {
	ag.stuck = ag.isStuck()
	if ag.stuck {
		return

	}

}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	if ag.position == ag.destination && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") {
		fmt.Println(ag.id, "disapear")
		ag.decision = Disapear
	} else if ag.stuck {
		ag.decision = Wait
	} else {
		ag.decision = Move
	}
}

func (ul *UsagerLambda) Act(ag *Agent) {
	if ag.decision == Move {
		ag.MoveAgent()
	} else if ag.decision == Wait {
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	} else if ag.decision == Disapear {
		RemoveAgent(&ag.env.station, ag)
	}

}

