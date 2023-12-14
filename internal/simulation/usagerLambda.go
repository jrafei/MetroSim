package simulation

import (
	"fmt"
	"math/rand"
	"time"
)

type UsagerLambda struct{}

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
