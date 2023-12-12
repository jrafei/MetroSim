package simulation

/*
 * Classe et méthodes principales de la structure Agent
 * à faire :
 *			// TODO: Gérer les moments où les agents font du quasi-sur place car ils ne peuvent plus bouger
 *			// TODO: Il arrive encore que certains agents soient bloqués, mais c'est quand il n'y a aucun mouvement possible.
 *			// Il faudrait faire en sorte que les agents bougent et laisse passer les autres
 *			// TODO: vérifier map playground, destination en (0,0)
 */

import (
	//"fmt"
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
	id              AgentID
	vitesse         time.Duration
	force           int
	politesse       bool
	position        Coord // Coordonnées de référence, width et height on compte width et height à partir de cette position
	departure       Coord
	destination     Coord
	behavior        Behavior
	env             *Environment
	syncChan        chan int
	decision        int
	isOn            map[Coord]string // Contenu de la case sur laquelle il se trouve
	stuck           bool
	width           int
	height          int
	orientation     int
	path            []alg.Node
	visitedPanneaux map[Coord]bool
}

type Behavior interface {
	Percept(*Agent)
	Deliberate(*Agent)
	Act(*Agent)
}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, behavior Behavior, departure, destination Coord, width, height int) *Agent {
	isOn := make(map[Coord]string)
	saveCells(&env.station, isOn, departure, width, height, 0)
	visitedPanneaux := make(map[Coord]bool, len(env.panneaux[env.zones[destination]]))
	for _, panneau := range env.panneaux[env.zones[destination]] {
		visitedPanneaux[panneau] = false
	}
	return &Agent{AgentID(id), vitesse, force, politesse, departure, departure, destination, behavior, env, syncChan, Noop, isOn, false, width, height, 0, make([]alg.Node, 0), visitedPanneaux}
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

	// ================== Tentative de calcul du chemin =======================
	if len(ag.path) == 0 {
		start, end := ag.generatePathExtremities()
		// Recherche d'un chemin si inexistant
		path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0))
		ag.path = path
	}
	// ================== Etude de faisabilité =======================
	// fmt.Println(ag.position,path[0])
	if IsAgentBlocking(ag.path, ag, ag.env) {
		start, end := ag.generatePathExtremities()
		// Si un agent bloque notre déplacement, on attend un temps aléatoire, et reconstruit un chemin en évitant la position
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		path := alg.FindPath(ag.env.station, start, end, ag.path[0])
		time.Sleep(time.Second)
		ag.path = path
	}
	if IsMovementSafe(ag.path, ag, ag.env) {
		removeAgent(&ag.env.station, ag)
		rotateAgent(ag, ag.path[0].Or())
		//ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = ag.isOn
		ag.position[0] = ag.path[0].Row()
		ag.position[1] = ag.path[0].Col()
		if len(ag.path) > 1 {
			ag.path = ag.path[1:]
		} else {
			ag.path = nil
		}
		saveCells(&ag.env.station, ag.isOn, ag.position, ag.width, ag.height, ag.orientation)
		//ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = "A"
		writeAgent(&ag.env.station, ag)
		// ============ Prise en compte de la vitesse de déplacement ======================
		time.Sleep(ag.vitesse * time.Millisecond)
		//fmt.Println(path[0])
		//fmt.Println(ag.position)
	}

}

func (ag *Agent) generatePathExtremities() (alg.Node, alg.Node) {
	start := *alg.NewNode(ag.position[0], ag.position[1], 0, 0, ag.width, ag.height)
	destination := ag.findDestination()
	end := *alg.NewNode(destination[0], destination[1], 0, 0, ag.width, ag.height)
	return start, end
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
			matrix[i][j] = "A"
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
		if equalCoord(&coord, &to_remove) {
			delete(mapping, coord)
		}
	}
}

func equalCoord(coord1, coord2 *Coord) bool {
	return coord1[0] == coord2[0] && coord1[1] == coord2[1]
}

func rotateAgent(agt *Agent, orientation int) {
	agt.orientation = orientation
}

func (ag *Agent) findDestination() Coord {
	destinationZone := ag.env.zones[ag.destination]
	if destinationZone != ag.env.zones[ag.position] {
		// Si on n'est pas dans la zone de la destination , on va s'orienter par un panneau
		//estimDistPos := alg.Heuristic(ag.position[0], ag.position[1], *alg.NewNode(ag.destination[0], ag.destination[1], 0, 0, 0, 0))
		for _, panneau := range ag.env.panneaux[destinationZone] {
			// On se rapproche du panneau menant à la zone
			//estimDistPan := alg.Heuristic(panneau[0], panneau[1], *alg.NewNode(ag.destination[0], ag.destination[1], 0, 0, 0, 0))
			if !ag.visitedPanneaux[panneau] {
				//TODO:revoir la mise à jour, peut-être à faire lorsqu'on se situe au niveau de panneau, pas avant
				//TODO:trouver une meilleure heuristique
				ag.visitedPanneaux[panneau] = true
				return panneau
			}
		}
	}
	// Sinon, on tente d'aller directement à la destination
	return ag.destination
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
