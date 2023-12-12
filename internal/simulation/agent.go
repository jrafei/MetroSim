package simulation

/*
 * Classe et méthodes principales de la structure Agent
 * à faire :
 *			// TODO: Gérer les moments où les agents font du quasi-sur place car ils ne peuvent plus bouger
 *			// TODO: Il arrive encore que certains agents soient bloqués, mais c'est quand il n'y a aucun mouvement possible.
 *			// Il faudrait faire en sorte que les agents bougent et laisse passer les autres
 *
 */

import (
	"log"
	"math/rand"
	alg "metrosim/internal/algorithms"
	"time"
	"fmt"
)

type Action int64

const (
	Noop = iota
	Mark
	Wait
	Move
	Expel // decision du Controleur
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
	orientation int //0 : vers le haut, 1 : vers la droite, 2 : vers le bas, 3 : vers la gauche
	request     *Request
}

type Request struct {
	demandeur AgentID
	decision int
}


type Behavior interface {
	Percept(*Agent)
	Deliberate(*Agent)
	Act(*Agent)
}

func NewRequest(demandeur AgentID, decision int) (req *Request) {
	return &Request{demandeur, decision}
}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, behavior Behavior, departure, destination Coord, width, height int) *Agent {
	isOn := make(map[Coord]string)
	saveCells(&env.station, isOn, departure, width, height, 0)
	return &Agent{AgentID(id), vitesse, force, politesse, departure, departure, destination, behavior, env, syncChan, Noop, isOn, false, width, height, 0, nil}
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

func IsMovementSafe(path []alg.Node, agt *Agent, env *Environment) bool {
	// Détermine si le movement est faisable

	if len(path) <= 0 {
		return false
	}

	// Simulation du déplacement
	ag := *agt
	ag.position = Coord{path[0].Row(), path[0].Col()}
	rotateAgent(&ag, path[0].Or())

	// Calcul des bornes de position de l'agent avant mouvement
	infRow, supRow, infCol, supCol := calculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	// Calcul des bornes de position de l'agent après mouvement
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(ag.position, ag.width, ag.height, ag.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			if !(j >= infCol && j < supCol && i >= infRow && i < supRow) && (env.station[i][j] != "B" && env.station[i][j] != "_") {
				// Si on est sur un agent, en dehors de la zone qu'occupe l'agent avant déplacement, on est bloqué
				return false
			}
		}
	}
	return true
}

func IsAgentBlocking(path []alg.Node, agt *Agent, env *Environment) bool {
	// Détermine si un agent se trouve sur la case à visiter

	if len(path) <= 0 {
		return false
	}

	// Simulation du déplacement
	ag := *agt
	ag.position = Coord{path[0].Row(), path[0].Col()}
	rotateAgent(&ag, path[0].Or())

	// Calcul des bornes de position de l'agent avant mouvement
	infRow, supRow, infCol, supCol := calculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	// Calcul des bornes de position de l'agent après mouvement
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(ag.position, ag.width, ag.height, ag.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			if !(j >= infCol && j < supCol && i >= infRow && i < supRow) && env.station[i][j] == "A" {
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

	count := 0

	// Calcul des bornes de position de l'agent après mouvement
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(ag.position, ag.width, ag.height, ag.orientation)

	for i := borneInfRow - 1; i < borneSupRow+1; i++ {
		for j := borneInfCol - 1; j < borneSupCol+1; j++ {
			// Éviter les cases à l'intérieur du rectangle
			if i >= borneInfRow && i < borneSupRow && j >= borneInfCol && j < borneSupCol {
				continue
			} else {
				count++
			}
			// Case inaccessible
			if i < 0 || j < 0 || i > 19 || j > 19 || ag.env.station[i][j] == "X" || ag.env.station[i][j] == "Q" || ag.env.station[i][j] == "A" {
				not_acc++

			}
			// fmt.Printf("Border (%d, %d) = %s \n", i, j, ag.env.station[i][j])
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
	timea := time.Now()
	path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0))
	fmt.Println(time.Since(timea))
	// ================== Etude de faisabilité =======================
	// fmt.Println(ag.position,path[0])
	if IsAgentBlocking(path, ag, ag.env) {
		// Si un agent bloque notre déplacement, on attend un temps aléatoire, et reconstruit un chemin en évitant la position
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		path = alg.FindPath(ag.env.station, start, end, path[0])
		time.Sleep(time.Second)
	}
	if IsMovementSafe(path, ag, ag.env) {
		removeAgent(&ag.env.station, ag)
		rotateAgent(ag, path[0].Or())
		//ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = ag.isOn
		ag.position[0] = path[0].Row()
		ag.position[1] = path[0].Col()
		saveCells(&ag.env.station, ag.isOn, ag.position, ag.width, ag.height, ag.orientation)
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

	// Calcul des bornes de position de l'agent
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			matrix[i][j] = agt.isOn[Coord{i, j}]
			removeCoord(Coord{i, j}, agt.isOn)
		}
	}
}

func writeAgent(matrix *[20][20]string, agt *Agent) {
	// Ecris l'agent dans la matrice

	// Calcul des bornes de position de l'agent
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			matrix[i][j] = string(agt.id)
		}
	}

}

func saveCells(matrix *[20][20]string, savedCells map[Coord]string, position Coord, width, height, orientation int) {
	// Enregistrement des valeurs des cellules de la matrice
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(position, width, height, orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
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

func rotateAgent(agt *Agent, orientation int) {
	agt.orientation = orientation
}

func calculateBounds(position Coord, width, height, orientation int) (infRow, supRow, infCol, supCol int) {
	borneInfRow := 0
	borneSupRow := 0
	borneInfCol := 0
	borneSupCol := 0

	// Calcul des bornes de position de l'agent après mouvement
	switch orientation {
	case 0:
		borneInfRow = position[0] - width + 1
		borneSupRow = position[0] + 1
		borneInfCol = position[1]
		borneSupCol = position[1] + height
	case 1:
		borneInfRow = position[0]
		borneSupRow = position[0] + height
		borneInfCol = position[1]
		borneSupCol = position[1] + width
	case 2:
		borneInfRow = position[0]
		borneSupRow = position[0] + width
		borneInfCol = position[1]
		borneSupCol = position[1] + height
	case 3:
		borneInfRow = position[0]
		borneSupRow = position[0] + height
		borneInfCol = position[1] - width + 1
		borneSupCol = position[1] + 1

	}
	return borneInfRow, borneSupRow, borneInfCol, borneSupCol
}
