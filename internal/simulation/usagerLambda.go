package simulation

import (
	"math/rand"
	"time"
)

type UsagerLambda struct{
	req Request
}

func (ul *UsagerLambda) Percept(ag *Agent) {
	// récupérer le channel de l'agent lambda
	chan_agt := ag.env.GetAgentChan(ag.id) 
	select {
	case req := <-chan_agt : //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		ul.req = req
	case <- time.After(time.Second):
		ag.stuck = ag.isStuck()
	}
}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	if ul.req.decision == Wait{
		ag.decision = Wait
	} else if ul.req.decision == Expel{
		ag.decision = Expel
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
	} else {
		ag.destination = ul.findNearestExit(ag.env)
		ag.MoveAgent()
	}
}


/*
 * Fonction qui permet de trouver la sortie la plus proche
*/
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
