package simulation

/*
 * Classe et méthodes principales de la structure Agent
 * à faire :
 *			// TODO: Gérer les moments où les agents font du quasi-sur place car ils ne peuvent plus bouger
 *			// TODO: Il arrive encore que certains agents soient bloqués, mais c'est quand il n'y a aucun mouvement possible.
 *			// Il faudrait faire en sorte que les agents bougent et laisse passer les autres
 */

import (
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
	Disappear
	Expel // virer l'agent
	Stop  // arreter l'agent
	GiveInfos
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
	orientation int //0 : vers le haut, 1 : vers la droite, 2 : vers le bas, 3 : vers la gauche (sens de deuxieme case occupé par l'agent)
	path        []alg.Node
	request     *Request
	direction   int //0 : vers le haut, 1 : vers la droite, 2 : vers le bas, 3 : vers la gauche (sens de son deplacement)
	// visitedPanneaux map[alg.Node]bool
	// visiting        *alg.Node
}

type Request struct {
	demandeur chan Request //channel de l'émetteur de la demande
	decision  int
}

type Behavior interface {
	Percept(*Agent)
	Deliberate(*Agent)
	Act(*Agent)
}

func NewRequest(demandeur chan Request, decision int) (req *Request) {
	return &Request{demandeur, decision}
}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, behavior Behavior, departure, destination Coord, width, height int) *Agent {
	isOn := make(map[Coord]string)
	return &Agent{AgentID(id), vitesse, force, politesse, departure, departure, destination, behavior, env, syncChan, Noop, isOn, false, width, height, 3, make([]alg.Node, 0), nil, 0}
}

func (ag *Agent) ID() AgentID {
	return ag.id
}

func (ag *Agent) Start() {
	log.Printf("%s starting...\n", ag.id)
	go ag.listenForRequests()
	go func() {
		var step int
		for {
			step = <-ag.syncChan
			ag.behavior.Percept(ag)
			ag.behavior.Deliberate(ag)
			ag.behavior.Act(ag)
			ag.syncChan <- step
			if ag.decision == Disappear {
				ag.env.RemoveAgent(*ag)
				return
			}
		}
	}()
}

func (ag *Agent) Act(env *Environment) {
	if ag.decision == Noop {
		env.Do(Noop, Coord{})
	}
}

func IsMovementSafe(path []alg.Node, agt *Agent, env *Environment) (bool, int) {
	// Détermine si le movement est faisable

	if len(path) <= 0 {
		return false, agt.orientation
	}
	// Calcul des bornes de position de l'agent avant mouvement
	infRow, supRow, infCol, supCol := calculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	// Si pas encore sur la map, mais agent déja sur la position, on ne peut pas encore apparaître
	if len(agt.isOn) == 0 && len(env.station[agt.path[0].Row()][agt.path[0].Col()]) > 1 {
		return false, agt.orientation
	}
	// Simulation du déplacement
	ag := *agt
	ag.position = Coord{path[0].Row(), path[0].Col()}
	for or := 0; or < 4; or++ {
		rotateAgent(&ag, or)
		safe := true

		// Calcul des bornes de position de l'agent après mouvement

		borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(ag.position, ag.width, ag.height, ag.orientation)
		if !(borneInfCol < 0 || borneInfRow < 0 || borneSupRow > 50 || borneSupCol > 50) {
			for i := borneInfRow; i < borneSupRow; i++ {
				for j := borneInfCol; j < borneSupCol; j++ {
					if !(j >= infCol && j < supCol && i >= infRow && i < supRow) && (env.station[i][j] != "B" && env.station[i][j] != "_" && env.station[i][j] != "W" && env.station[i][j] != "S") {
						// Si on n'est pas sur une case atteignable, en dehors de la zone qu'occupe l'agent avant déplacement, on est bloqué
						safe = false
					}
				}
			}
			if safe {
				return true, or
			}
		}

	}
	return false, agt.orientation
}

func IsAgentBlocking(path []alg.Node, agt *Agent, env *Environment) bool {
	// Détermine si le movement est faisable
	if len(path) <= 0 {
		return false
	}
	// Calcul des bornes de position de l'agent avant mouvement
	infRow, supRow, infCol, supCol := calculateBounds(agt.position, agt.width, agt.height, agt.orientation)
	// Simulation du déplacement
	ag := *agt
	ag.position = Coord{path[0].Row(), path[0].Col()}
	for or := 0; or < 4; or++ {
		rotateAgent(&ag, or)
		blocking := false
		// Calcul des bornes de position de l'agent après mouvement
		borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(ag.position, ag.width, ag.height, ag.orientation)
		//fmt.Println(ag.id,borneInfRow,borneInfRow, borneSupRow, borneInfCol, borneSupCol)
		if !(borneInfCol < 0 || borneInfRow < 0 || borneSupRow > 50 || borneSupCol > 50) {
			for i := borneInfRow; i < borneSupRow; i++ {
				for j := borneInfCol; j < borneSupCol; j++ {
					if !(j >= infCol && j < supCol && i >= infRow && i < supRow) && len(env.station[i][j]) > 2 {
						// Si on n'est pas sur une case atteignable, en dehors de la zone qu'occupe l'agent avant déplacement, on est bloqué
						blocking = true
					}
				}
			}
			if !blocking {
				// Si on n'a pas trouvé d'agent bloquant pour cette nouvelle position, on retourne faux
				return false
			}
		}
	}
	// Le cas où dans tous les mouvements on est bloqué par un agent
	return true
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
			if i < 0 || j < 0 || i > 19 || j > 19 || ag.env.station[i][j] == "X" || ag.env.station[i][j] == "Q" || len(ag.env.station[i][j]) > 2 {
				not_acc++

			}
			// fmt.Printf("Border (%d, %d) = %s \n", i, j, ag.env.station[i][j])
		}
	}

	// Si aucune case disponible autour de lui, il est bloqué
	return not_acc == count
}

func (ag *Agent) WhichAgent() string {
	if ag.direction == 0 { // vers le haut
		return ag.env.station[ag.position[0]-1][ag.position[1]]
	} else if ag.direction == 1 { // vers la droite
		return ag.env.station[ag.position[0]][ag.position[1]+1]
	} else if ag.direction == 2 { // vers le bas
		return ag.env.station[ag.position[0]+1][ag.position[1]]
	} else { // vers la gauche
		return ag.env.station[ag.position[0]][ag.position[1]-1]
	}
}

func (env *Environment) FindAgentByID(agtId AgentID) *Agent {
	for i := range env.ags {
		if env.ags[i].id == agtId {
			return &env.ags[i]
		}
	}
	return nil
}

func (ag *Agent) MoveAgent() {
	//fmt.Println("[Agent, MoveAgent] destination ", ag.destination)

	// ================== Tentative de calcul du chemin =======================
	if len(ag.path) == 0 {
		start, end := ag.generatePathExtremities()
		// Recherche d'un chemin si inexistant
		path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 2*time.Second)
		ag.path = path
	}

	// ================== Etude de faisabilité =======================
	if IsAgentBlocking(ag.path, ag, ag.env) {

		if ag.politesse {
			start, end := ag.generatePathExtremities()
			// Si un agent bloque notre déplacement, on attend un temps aléatoire, et reconstruit
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 2*time.Second)
			ag.path = path
			return
		} else {
			//Si individu impoli, demande à l'agent devant de bouger
			//On récupère le id de la personne devant
			blockingAgentID := AgentID(ag.WhichAgent())
			//blockingAgent := ag.env.FindAgentByID(blockingAgentID)
			var reqToBlockingAgent *Request
			//var reqToImpoliteAgent *Request
			i := 0
			accept := false
			for !accept && i < 3 {
				//Demande à l'agent qui bloque de se pousser (réitère trois fois s'il lui dit pas possible)
				i += 1
				reqToBlockingAgent = NewRequest(ag.env.agentsChan[ag.id], 3) //Création "Hello, je suis ag.id, move."
				ag.env.agentsChan[blockingAgentID] <- *reqToBlockingAgent    //Envoi requête

				/*
					1. Faire le moment ou blocking agent recoit qqchose sur son canal
					2.


				*/
				/*
					//BlockingAgent cherche si autour de lui c'est vide
					possible, or := IsMovementSafe(blockingAgent.path, blockingAgent, blockingAgent.env)

					if !possible {
						reqToImpoliteAgent = NewRequest(ag.id, 0)
						ag.env.agentsChan[ag.id] <- *reqToImpoliteAgent
					} else {
						//Bouge sur la case possible
						accept = true
						coordBlockingAgent := blockingAgent.position
						//Gérer le déplacement de Ag et de BlockingAgent + déplacement en fonction de la force !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
					}
				*/
			}
		}
	}

	// ================== Déplacement si aucun problème ou si blockingAgent se pousse =======================
	safe, or := IsMovementSafe(ag.path, ag, ag.env)
	if safe {
		if len(ag.isOn) > 0 {
			RemoveAgent(&ag.env.station, ag)
		}
		rotateAgent(ag, or)
		ag.direction = calculDirection(ag.position, Coord{ag.path[0].Row(), ag.path[0].Col()})
		ag.position[0] = ag.path[0].Row()
		ag.position[1] = ag.path[0].Col()
		if len(ag.path) > 1 {
			ag.path = ag.path[1:]
		} else {
			ag.path = nil
		}
		saveCells(&ag.env.station, ag.isOn, ag.position, ag.width, ag.height, ag.orientation)
		writeAgent(&ag.env.station, ag)
		// ============ Prise en compte de la vitesse de déplacement ======================
		time.Sleep(ag.vitesse * time.Millisecond)
	}

}

func (ag *Agent) generatePathExtremities() (alg.Node, alg.Node) {
	// Génère les points extrêmes du chemin de l'agent
	start := *alg.NewNode(ag.position[0], ag.position[1], 0, 0, ag.width, ag.height)
	destination := ag.destination
	end := *alg.NewNode(destination[0], destination[1], 0, 0, ag.width, ag.height)
	return start, end
}

func RemoveAgent(matrix *[50][50]string, agt *Agent) {
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

func writeAgent(matrix *[50][50]string, agt *Agent) {
	// Ecris l'agent dans la matrice

	// Calcul des bornes de position de l'agent
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := calculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			matrix[i][j] = string(agt.id)
		}
	}

}

func saveCells(matrix *[50][50]string, savedCells map[Coord]string, position Coord, width, height, orientation int) {
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
	// Vérifie l'égalité de 2 objets Coord
	return coord1[0] == coord2[0] && coord1[1] == coord2[1]
}

// Fonction utilitaire de rotation
func rotateAgent(agt *Agent, orientation int) {
	agt.orientation = orientation
}

func calculateBounds(position Coord, width, height, orientation int) (infRow, supRow, infCol, supCol int) {
	// Fonction de génération des frontières d'un objet ayant une largeur et une hauteur, en focntion de son orientation
	borneInfRow := 0
	borneSupRow := 0
	borneInfCol := 0
	borneSupCol := 0

	// Calcul des bornes de position de l'agent après mouvement
	switch orientation {
	case 0:
		// Orienté vers le haut
		borneInfRow = position[0] - width + 1
		borneSupRow = position[0] + 1
		borneInfCol = position[1]
		borneSupCol = position[1] + height
	case 1:
		// Orienté vers la droite
		borneInfRow = position[0]
		borneSupRow = position[0] + height
		borneInfCol = position[1]
		borneSupCol = position[1] + width
	case 2:
		// Orienté vers le bas
		borneInfRow = position[0]
		borneSupRow = position[0] + width
		borneInfCol = position[1]
		borneSupCol = position[1] + height
	case 3:
		// Orienté vers la gauche
		borneInfRow = position[0]
		borneSupRow = position[0] + height
		borneInfCol = position[1] - width + 1
		borneSupCol = position[1] + 1

	}
	return borneInfRow, borneSupRow, borneInfCol, borneSupCol
}

func (ag *Agent) listenForRequests() {
	for {
		if ag.request == nil {
			req := <-ag.env.agentsChan[ag.id]
			fmt.Println("Request received by UsagerLambda:", req.decision)
			ag.request = &req
			if req.decision == Disappear {
				return
			}
		}
	}
}
