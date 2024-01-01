package simulation

/*
 * Classe et méthodes principales de la structure Agent
 */

import (
	"fmt"
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"time"
)

type Action int64

const (
	Noop       = iota
	Wait       // Attente
	Move       // Déplacement de l'agent
	EnterMetro //Entrer dans le métro
	Disappear  // Disparition  de l'agent dans la simulation
	Expel      // virer l'agent
	Stop       // arreter l'agent
	ACK        // acquittement
)

type AgentID string

type Agent struct {
	id          AgentID
	vitesse     time.Duration
	force       int
	politesse   bool
	position    alg.Coord // Coordonnées de référence, width et height on compte width et height à partir de cette position
	departure   alg.Coord
	destination alg.Coord
	behavior    Behavior
	env         *Environment
	syncChan    chan int
	decision    int
	isOn        map[alg.Coord]string // Contenu de la case sur laquelle il se trouve
	stuck       bool
	width       int
	height      int
	orientation int //0 : vers le haut, 1 : vers la droite, 2 : vers le bas, 3 : vers la gauche (sens de construction de l'agent)
	path        []alg.Node
	request     *req.Request
	direction   int //0 : vers le haut, 1 : vers la droite, 2 : vers le bas, 3 : vers la gauche (sens de son deplacement)

}

type Behavior interface {
	Percept(*Agent)
	Deliberate(*Agent)
	Act(*Agent)
}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, behavior Behavior, departure, destination alg.Coord, width, height int) *Agent {
	isOn := make(map[alg.Coord]string)
	return &Agent{AgentID(id), vitesse, force, politesse, departure, departure, destination, behavior, env, syncChan, Noop, isOn, false, width, height, 3, make([]alg.Node, 0), nil, 0}
}

func (ag *Agent) ID() AgentID {
	return ag.id
}

func (ag *Agent) Start() {
	//log.Printf("%s starting...\n", ag.id)
	go ag.listenForRequests()
	go func() {
		var step int
		for {
			step = <-ag.syncChan
			ag.behavior.Percept(ag)
			ag.behavior.Deliberate(ag)
			ag.behavior.Act(ag)
			ag.syncChan <- step
			//fmt.Println(ag.id, ag.path)
			if ag.decision == Disappear || ag.decision == EnterMetro {
				ag.env.DeleteAgent(*ag)
				return
			}
		}
	}()
}


func (agt *Agent) IsMovementSafe() (bool, int) {
	// Détermine si le movement est faisable

	if len(agt.path) <= 0 {
		return false, agt.orientation
	}
	// Calcul des bornes de position de l'agent avant mouvement
	infRow, supRow, infCol, supCol := alg.CalculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	// Si pas encore sur la map, mais agent déja sur la position, on ne peut pas encore apparaître
	if len(agt.isOn) == 0 && existAgent(agt.env.station[agt.path[0].Row()][agt.path[0].Col()]) {
		return false, agt.orientation
	}
	// Simulation du déplacement
	ag := *agt
	ag.position = alg.Coord{agt.path[0].Row(), agt.path[0].Col()}
	for or := 0; or < 4; or++ {
		ag.orientation = or
		safe := true

		// Calcul des bornes de position de l'agent après mouvement

		borneInfRow, borneSupRow, borneInfCol, borneSupCol := alg.CalculateBounds(ag.position, ag.width, ag.height, ag.orientation)
		if !(borneInfCol < 0 || borneInfRow < 0 || borneSupRow > 50 || borneSupCol > 50) {
			for i := borneInfRow; i < borneSupRow; i++ {
				for j := borneInfCol; j < borneSupCol; j++ {
					if agt.env.station[i][j] == "O" {
						// Vérification si porte de métro
						metro := findMetro(ag.env, &alg.Coord{i, j})
						if metro != nil && !metro.way.gatesClosed && alg.EqualCoord(&ag.destination, &alg.Coord{i, j}) {
							// On s'assure que les portes ne sont pas fermées et que c'est la destination
							return true, or
						} else {
							safe = false
						}
					}
					if !(j >= infCol && j < supCol && i >= infRow && i < supRow) && (agt.env.station[i][j] != "B" && agt.env.station[i][j] != "_" && agt.env.station[i][j] != "W" && agt.env.station[i][j] != "S") {
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

func (agt *Agent) IsAgentBlocking() bool {
	// Détermine si le movement est faisable
	if len(agt.path) <= 0 {
		return false
	}
	// Calcul des bornes de position de l'agent avant mouvement
	infRow, supRow, infCol, supCol := alg.CalculateBounds(agt.position, agt.width, agt.height, agt.orientation)
	// Simulation du déplacement
	ag := *agt
	ag.position = alg.Coord{agt.path[0].Row(), agt.path[0].Col()}
	for or := 0; or < 4; or++ {
		ag.orientation = or
		blocking := false
		// Calcul des bornes de position de l'agent après mouvement
		borneInfRow, borneSupRow, borneInfCol, borneSupCol := alg.CalculateBounds(ag.position, ag.width, ag.height, ag.orientation)
		//fmt.Println(ag.id,borneInfRow,borneInfRow, borneSupRow, borneInfCol, borneSupCol)
		if !(borneInfCol < 0 || borneInfRow < 0 || borneSupRow > 50 || borneSupCol > 50) {
			for i := borneInfRow; i < borneSupRow; i++ {
				for j := borneInfCol; j < borneSupCol; j++ {
					if !(j >= infCol && j < supCol && i >= infRow && i < supRow) && existAgent(ag.env.station[i][j]) {
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
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := alg.CalculateBounds(ag.position, ag.width, ag.height, ag.orientation)

	for i := borneInfRow - 1; i < borneSupRow+1; i++ {
		for j := borneInfCol - 1; j < borneSupCol+1; j++ {
			// Éviter les cases à l'intérieur du rectangle
			if i >= borneInfRow && i < borneSupRow && j >= borneInfCol && j < borneSupCol {
				continue
			} else {
				count++
			}
			// Case inaccessible
			if i < 0 || j < 0 || i > 19 || j > 19 || ag.env.station[i][j] == "X" || ag.env.station[i][j] == "Q" || ag.env.station[i][j] == "M" || existAgent(ag.env.station[i][j]){
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

func (ag *Agent) MoveAgent() {
	//fmt.Println("[Agent, MoveAgent] destination ", ag.destination)

	// ================== Tentative de calcul du chemin =======================
	if len(ag.path) == 0 || ag.isGoingToExitPath() || (ag.env.station[ag.path[0].Row()][ag.path[0].Col()]=="O"&& !alg.EqualCoord(&ag.destination,&alg.Coord{ag.path[0].Row(),ag.path[0].Col()})) {
		start, end := ag.generatePathExtremities()
		// Recherche d'un chemin si inexistant
		if len(ag.path) > 0 {
			ag.path = alg.FindPath(ag.env.station, start, end, ag.path[0], false, 2*time.Second)
		} else {
			ag.path = alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 2*time.Second)
		}
	}

	// ================== Etude de faisabilité =======================
	if ag.IsAgentBlocking() {

		if ag.politesse {
			start, end := ag.generatePathExtremities()
			// Si un agent bloque notre déplacement, on attend un temps aléatoire, et reconstruit
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			//path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 2*time.Second)
			path := alg.FindPath(ag.env.station, start, end, ag.path[0], false, 2*time.Second)
			ag.path = path
			return
		} else {
			//Si individu impoli, demande à l'agent devant de bouger
			//On récupère le id de la personne devant
			blockingAgentID := AgentID(ag.WhichAgent())
			//blockingAgent := ag.env.FindAgentByID(blockingAgentID)
			var reqToBlockingAgent *req.Request
			//var reqToImpoliteAgent *Request
			i := 0
			accept := false
			for !accept && i < 3 {
				//Demande à l'agent qui bloque de se pousser (réitère trois fois s'il lui dit pas possible)
				i += 1
				reqToBlockingAgent = req.NewRequest(ag.env.agentsChan[ag.id], 3) //Création "Hello, je suis ag.id, move."
				ag.env.agentsChan[blockingAgentID] <- *reqToBlockingAgent        //Envoi requête

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
	safe, or := ag.IsMovementSafe()
	if safe {
		if len(ag.isOn) > 0 {
			ag.env.RemoveAgent(ag)
		}
		ag.orientation = or
		ag.direction = calculDirection(ag.position, alg.Coord{ag.path[0].Row(), ag.path[0].Col()})
		ag.position[0] = ag.path[0].Row()
		ag.position[1] = ag.path[0].Col()
		if len(ag.path) > 1 {
			ag.path = ag.path[1:]
		} else {
			ag.path = nil
		}
		ag.saveCells()
		ag.env.writeAgent(ag)
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

func (agt *Agent) saveCells() {
	// Enregistrement des valeurs des cellules de la matrice
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := alg.CalculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			agt.isOn[alg.Coord{i, j}] = agt.env.station[i][j]
		}
	}
}

func (ag *Agent) listenForRequests() {
	for {
		if ag.request == nil {
			req := <-ag.env.agentsChan[ag.id]
			fmt.Println("Request received by :", ag.id, req.Decision)
			ag.request = &req
		}
		if ag.request.Decision() == Disappear || ag.request.Decision() == EnterMetro {
			return
		}
	}
}

func (ag *Agent) isGoingToExitPath() bool {
	if len(ag.path) > 0 {
		for _, metro := range ag.env.metros {
			for gate_index, gate := range metro.way.gates {
				if alg.EqualCoord(&ag.destination, &gate) {
					// Si la destination est une porte de métro, on va essayer de libérer le chemin des agents sortants
					exit_path := metro.way.pathsToExit[gate_index]
					for _, cell := range exit_path {
						if alg.EqualCoord(&alg.Coord{cell.Row(), cell.Col()}, &alg.Coord{ag.path[0].Row(), ag.path[0].Col()}) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func findMetro(env *Environment, gateToFind *alg.Coord) *Metro {
	for _, metro := range env.metros {
		for _, gate := range metro.way.gates {
			if alg.EqualCoord(&gate, gateToFind) {
				return &metro
			}
		}
	}
	return nil
}
