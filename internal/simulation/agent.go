package simulation

/*
 * Classe et méthodes principales de la structure Agent
 * à faire :
 *			//TODO: gérer les orientations
 *			// TODO: Gérer les moments où les agents font du quasi-sur place car ils ne peuvent plus bouger
 *			// TODO: Il arrive encore que certains agents soient bloqués, mais c'est quand il n'y a aucun mouvement possible.
 *			// Il faudrait faire en sorte que les agents bougent et laisse passer les autres
 *
 */

import (
	//"fmt"

	"fmt"
	"log"
	"math/rand"
	alg "metrosim/internal/algorithms"
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
	id          AgentID
	vitesse     time.Duration
	force       int
	politesse   bool
	position    Coord // Coordonnées de référence, width et height on compte width et height à partir de cette position
	departure   Coord
	destination Coord
	behavior    Behavior
	env         *Environment
	syncChan    chan int
	decision    int
	isOn        map[Coord]string // Contenu de la case sur laquelle il se trouve
	stuck       bool
	width       int
	height      int
}

type Behavior interface {
	Percept(*Agent)
	Deliberate(*Agent)
	Act(*Agent)
}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, behavior Behavior, departure, destination Coord, width, height int) *Agent {
	isOn := make(map[Coord]string)
	saveCells(&env.station, isOn, departure, width, height)
	return &Agent{AgentID(id), vitesse, force, politesse, departure, departure, destination, behavior, env, syncChan, Noop, isOn, false, width, height}
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

func IsMovementSafe(path []alg.Node, ag *Agent, env *Environment) bool {
	// Détermine si le movement est faisable
	if len(path) <= 0 {
		return false
	}
	startX, startY := ag.position[1], ag.position[0]
	width, height := ag.width, ag.height

	// on simule son nouvel emplacement
	// et regarde s'il chevauche un autre agent
	for i := path[0].Row(); i < path[0].Row()+ag.height; i++ {
		for j := path[0].Col(); j < path[0].Col()+ag.width; j++ {

			if !(j >= startX && j < startX+width && i >= startY && i < startY+height) && (env.station[i][j] != "B" && env.station[i][j] != "_") {
				// Si on est sur une case non atteignable, en dehors de la zone qu'occupe l'agent avant déplacement, on est bloqué
				return false
			}
		}
	}
	return true
}

func IsAgentBlocking(path []alg.Node, ag *Agent, env *Environment) bool {
	// Détermine si un agent se trouve sur la case à visiter
	// Coordonnée de départ et dimensions du rectangle
	if len(path) <= 0 {
		return false
	}
	startX, startY := ag.position[1], ag.position[0]
	width, height := ag.width, ag.height

	// on simule son nouvel emplacement
	// et regarde s'il chevauche un autre agent
	for i := path[0].Row(); i < path[0].Row()+ag.height; i++ {
		for j := path[0].Col(); j < path[0].Col()+ag.width; j++ {

			if !(j >= startX && j < startX+width && i >= startY && i < startY+height) && env.station[i][j] == "A" {
				// Si on est sur un agent, en dehors de la zone qu'occupe l'agent avant déplacement, on est bloqué
				return true
			}
		}
	}

	return false
}

func (ag *Agent) isStuck() bool {
	// Perception des éléments autour de l'agent pour déterminer si bloqué
	not_acc := 0 // nombre de cases indisponibles autour de l'agent

	// Coordonnée de départ et dimensions du rectangle
	startX, startY := ag.position[1], ag.position[0]
	width, height := ag.width, ag.height

	// Largeur et hauteur du rectangle étendu
	extendedWidth := width + 2   // +2 pour les cases à gauche et à droite du rectangle
	extendedHeight := height + 2 // +2 pour les cases au-dessus et en dessous du rectangle

	count := 0
	// Parcourir les cases autour du rectangle
	for i := startX - 1; i < startX+extendedWidth-1; i++ {
		for j := startY - 1; j < startY+extendedHeight-1; j++ {
			// Éviter les cases à l'intérieur du rectangle
			if i >= startX && i < startX+width && j >= startY && j < startY+height {

				continue
			} else {
				count++
			}
			// Case inaccessible
			if ag.env.station[j][i] == "X" || ag.env.station[j][i] == "Q" || ag.env.station[j][i] == "A" {
				not_acc++

			}
			// fmt.Printf("Border (%d, %d) = %s \n", j, i,ag.env.station[j][i])
		}
	}
	// Si aucune case disponible autour de lui, il est bloqué
	return not_acc == count
}

func (ag *Agent) MoveAgent() {

	// ============ Initialisation des noeuds de départ ======================
	start := *alg.NewNode(ag.position[0], ag.position[1], 0, 0, ag.width, ag.height)
	end := *alg.NewNode(ag.destination[0], ag.destination[1], 0, 0, ag.width, ag.height)
	// ================== Tentative de calcul du chemin =======================
	path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0))
	// ================== Etude de faisabilité =======================
	fmt.Println(path)
	if IsAgentBlocking(path, ag, ag.env) {
		// Si un agent bloque notre déplacement, on attend un temps aléatoire, et reconstruit un chemin en évitant la position
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		path = alg.FindPath(ag.env.station, start, end, path[0])
		time.Sleep(time.Second)
	}
	if IsMovementSafe(path, ag, ag.env) {
		removeAgent(&ag.env.station, ag)
		//ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = ag.isOn
		ag.position[0] = path[0].Row()
		ag.position[1] = path[0].Col()
		saveCells(&ag.env.station, ag.isOn, ag.position, ag.width, ag.height)
		//ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = "A"
		writeAgent(&ag.env.station, ag)
		// ============ Prise en compte de la vitesse de déplacement ======================
		time.Sleep(ag.vitesse * time.Millisecond)
		//fmt.Println(path[0])
		//fmt.Println(ag.position)
	}

}

func removeAgent(matrix *[20][20]string, agt *Agent) {
	// Supprime l'agent de la matrice
	for i := agt.position[0]; i < agt.position[0]+agt.height; i++ {
		for j := agt.position[1]; j < agt.position[1]+agt.width; j++ {
			matrix[i][j] = agt.isOn[Coord{i, j}]
			removeCoord(Coord{i, j}, agt.isOn)
		}
	}
}

func writeAgent(matrix *[20][20]string, agt *Agent) {
	// Ecris un agent dans la matrice
	for i := agt.position[0]; i < agt.position[0]+agt.height; i++ {
		for j := agt.position[1]; j < agt.position[1]+agt.width; j++ {
			matrix[i][j] = "A"
		}
	}
}

func saveCells(matrix *[20][20]string, savedCells map[Coord]string, ref Coord, width, height int) {
	// Enregistrement des valeurs des cellules de la matrice
	for i := ref[0]; i < ref[0]+height; i++ {
		for j := ref[1]; j < ref[1]+width; j++ {
			savedCells[Coord{i, j}] = matrix[i][j]
		}
	}
}

func removeCoord(to_remove Coord, mapping map[Coord]string) {
	// Suppression d'une clé dans une map
	for coord, _ := range mapping {
		if coord[0] == to_remove[0] && coord[1] == to_remove[1] {
			delete(mapping, coord)
		}
	}
}
