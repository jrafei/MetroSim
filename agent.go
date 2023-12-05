package simulation

import (
	//"fmt"

	"log"
	"math/rand"
	"time"
)

type Action int64

const (
	Noop = iota
	Mark
	Wait
	Move
)

type Coord [2]int
type AgentID string

type Agent struct {
	id                  AgentID
	vitesse             time.Duration
	force               int
	politesse           bool
	coordHautOccupation Coord
	coordBasOccupation  Coord
	departure           Coord
	destination         Coord
	behavior            Behavior
	env                 *Environment
	syncChan            chan int
	decision            int
	isOn                string // Contenu de la case sur laquelle il se trouve
	stuck               bool
}

type Behavior interface {
	Percept(*Agent)
	Deliberate(*Agent)
	Act(*Agent)
}

type UsagerLambda struct{}

func (ul *UsagerLambda) Percept(ag *Agent) {
	ag.stuck = ag.isStuck()
	if ag.stuck {
		return
	}

}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	if ag.stuck {
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
	}

}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, UpCoord Coord, DownCoord Coord, behavior Behavior, departure, destination Coord) *Agent {
	return &Agent{AgentID(id), vitesse, force, politesse, UpCoord, DownCoord, departure, destination, behavior, env, syncChan, Noop, env.station[UpCoord[0]][UpCoord[1]], false}
}

func (ag *Agent) ID() AgentID {
	return ag.id
}

func (ag *Agent) Start() {
	log.Printf("%s starting...\n", ag.id)

	go func() {
		var step int
		for {
			step = <-ag.syncChan
			ag.behavior.Percept(ag)
			ag.behavior.Deliberate(ag)
			ag.behavior.Act(ag)
			ag.syncChan <- step
		}
	}()
}

func (ag *Agent) Act(env *Environment) {
	if ag.decision == Noop {
		env.Do(Noop, Coord{})
	}
}

func IsMovementSafe(path []Node, env *Environment) bool {
	// Détermine si le movement est faisable
	return len(path) > 0 && (env.station[path[0].row][path[0].col] == "B" || env.station[path[0].row][path[0].col] == "_")
}

func IsAgentBlocking(path []Node, env *Environment) bool {
	// Détermine si un agent se trouve sur la case à visiter
	return len(path) > 0 && env.station[path[0].row][path[0].col] == "A"
}

func (ag *Agent) isStuck() bool {
	// Perception des éléments autour de l'agent pour déterminer si bloqué
	s := 0 // nombre de cases indisponibles autour de l'agent
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			ord := (ag.coordBasOccupation[0] - 1) + i
			abs := (ag.coordBasOccupation[1] - 1) + j

			if ord != ag.coordBasOccupation[0] && abs != ag.coordBasOccupation[1] {
				if ord < 0 || abs < 0 || ord > 19 || abs > 19 {
					s++
				} else if ag.env.station[ord][abs] == "X" || ag.env.station[ord][abs] == "Q" || ag.env.station[ord][abs] == "A" {
					s++
				}
			}
		}
	}
	// Si aucune case disponible autour de lui, il est bloqué
	return s == 8
}

func (ag *Agent) MoveAgent() {
	// TODO: Gérer les moments où les agents font du quasi-sur place car il ne peuvent plus bouger
	// TODO: Parfois, certains agents ne bougent plus,
	// TODO: Il arrive encore que certains agents soient bloqués, mais c'est quand il n'y a aucun mouvement possible.
	// Il faudrait faire en sorte que les agents bougent et laisse passer les autres

	// ============ Initialisation des noeuds de départ ======================
	start := Node{ag.coordBasOccupation[0], ag.coordBasOccupation[1], 0, 0}
	end := Node{ag.destination[0], ag.destination[1], 0, 0}

	// ================== Tentative de calcul du chemin =======================
	path := findPath(ag.env.station, start, end, Node{})

	// ================== Etude de faisabilité =======================
	if IsAgentBlocking(path, ag.env) {
		// Si un agent bloque notre déplacement, on attend un temps aléatoire, et reconstruit un chemin en évitant la position
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		path = findPath(ag.env.station, start, end, path[0])
		time.Sleep(time.Second)
	}

	if IsMovementSafe(path, ag.env) {
		ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = ag.isOn
		ag.isOn = ag.env.station[path[0].row][path[0].col]
		ag.coordBasOccupation[0] = path[0].row
		ag.coordBasOccupation[1] = path[0].col
		ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = "A"

		// ============ Prise en compte de la vitesse de déplacement ======================
		time.Sleep(ag.vitesse)
	}

}
