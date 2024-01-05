package simulation

/*
 * Classe et méthodes principales de la structure Agent
 */

import (
	//"fmt"

	//"log"
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"time"
)

type Action int64

const (
	Noop       = iota //No opération, utiliser pour refuser un mouvement 0
	Wait              // Attente
	Move              // Déplacement de l'agent
	EnterMetro        //Entrer dans le métro
	TryToMove
	YouHaveToMove //Utiliser par un usager impoli pour forcer un déplacement
	Done
	Disappear // Disparition  de l'agent dans la simulation
	Expel     // virer l'agent //8
	Stop      // arreter l'agent
	ACK       // acquittement
)

type AgentID string

type Agent struct {
	id      AgentID
	vitesse time.Duration
	//force       int
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
	SetUpDestination(ag *Agent)
}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, politesse bool, behavior Behavior, departure, destination alg.Coord, width, height int) *Agent {
	isOn := make(map[alg.Coord]string)
	direct := initDirection(departure, len(env.station[0]))
	return &Agent{AgentID(id), vitesse, politesse, departure, departure, destination, behavior, env, syncChan, Noop, isOn, false, width, height, 3, make([]alg.Node, 0), nil, direct}
}

func (ag *Agent) ID() AgentID {
	return ag.id
}

func (ag *Agent) Start() {

	//TODELETElog.Printf("%s starting...\n", ag.id)
	go ag.listenForRequests()

	// si c'est un controlleur on lance le timer de durée de vie
	if ag.id[0] == 'C' {
		//fmt.Println("[Start()] C'est un controleur")
		ag.behavior.(*Controleur).startTimer()
	}

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
	if len(agt.isOn) == 0 && len(agt.env.station[agt.path[0].Row()][agt.path[0].Col()]) > 1 {
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
		if !(borneInfCol < 0 || borneInfRow < 0 || borneSupRow > len(agt.env.station[0]) || borneSupCol > len(agt.env.station[1])) {
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
						//fmt.Println("[IsMovementSafe]case inaccessible :",agt.id)
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
		if !(borneInfCol < 0 || borneInfRow < 0 || borneSupRow > len(agt.env.station[0]) || borneSupCol > len(agt.env.station[1])) {
			for i := borneInfRow; i < borneSupRow; i++ {
				for j := borneInfCol; j < borneSupCol; j++ {
					if !(j >= infCol && j < supCol && i >= infRow && i < supRow) && len(ag.env.station[i][j]) > 2 {
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
			if i < 0 || j < 0 || i > len(ag.env.station[0])-1 || j > len(ag.env.station[0])-1 || ag.env.station[i][j] == "X" || ag.env.station[i][j] == "Q" || len(ag.env.station[i][j]) > 2 {
				not_acc++

			}
			// fmt.Printf("Border (%d, %d) = %s \n", i, j, ag.env.station[i][j])
		}
	}

	// Si aucune case disponible autour de lui, il est bloqué
	return not_acc == count
}

func (ag *Agent) NextCell() string {
	//0 : vers le haut, 1 : vers la droite, 2 : vers le bas, 3 : vers la gauche (sens de son deplacement)
	switch ag.direction {
	case 0: // vers le haut
		if ag.position[0]-1 >= 0 && ag.position[0]-1 < len(ag.env.station[0]) {
			return ag.env.station[ag.position[0]-1][ag.position[1]]
		}
	case 1: // vers la droite
		if ag.position[1]+1 >= 0 && ag.position[1]+1 < len(ag.env.station[1]) {
			return ag.env.station[ag.position[0]][ag.position[1]+1]
		}
	case 2: // vers le bas
		if ag.position[0]+1 >= 0 && ag.position[0]+1 < len(ag.env.station[0]) {
			return ag.env.station[ag.position[0]+1][ag.position[1]]
		}
	default: //vers la gauche
		if ag.position[1]-1 >= 0 && ag.position[1]-1 < len(ag.env.station[1]) {
			return ag.env.station[ag.position[0]][ag.position[1]-1]
		}
	}
	return "X"
}

func (agt *Agent) MyNextCellIsSafe() bool {

	// Simulation du déplacement
	ag := *agt
	switch ag.direction {
	case 0: //haut
		ag.position = alg.Coord{ag.position[0] - 1, ag.position[1]}
	case 1: //droite
		ag.position = alg.Coord{ag.position[0], ag.position[1] + 1}
	case 2: //gauche
		ag.position = alg.Coord{ag.position[0] + 1, ag.position[1]}
	case 3: //bas
		ag.position = alg.Coord{ag.position[0], ag.position[1] - 1}
	}

	if !(ag.position[1] < 0 || ag.position[0] < 0 || ag.position[0] > len(agt.env.station[0]) || ag.position[1] > len(agt.env.station[1])) {
		i := ag.position[0]
		j := ag.position[1]
		if agt.env.station[i][j] != "B" && agt.env.station[i][j] != "_" && agt.env.station[i][j] != "W" && agt.env.station[i][j] != "S" {
			// Si on n'est pas sur une case atteignable, en dehors de la zone qu'occupe l'agent avant déplacement, on est bloqué
			//fmt.Println("[IsMovementSafe]case inaccessible :",agt.id)
			return false
		} else {
			return true
		}
	}
	return false
}

func (ag *Agent) ShiftAgent() bool {
	//fmt.Printf("ShiftAgent")
	storeDirection := ag.direction //enregistrer l'orientation initiale de l'agent

	for i := 0; i < 4; i++ {
		ag.direction = (storeDirection + i) % 4
		safe := ag.MyNextCellIsSafe()
		if safe { //  Deplacement possible safe=true

			switch ag.direction {
			case 0: //haut
				ag.position = alg.Coord{ag.position[0] - 1, ag.position[1]}
			case 1: //droite
				ag.position = alg.Coord{ag.position[0], ag.position[1] + 1}
			case 2:
				ag.position = alg.Coord{ag.position[0] + 1, ag.position[1]}
			case 3:
				ag.position = alg.Coord{ag.position[0], ag.position[1] - 1}
			}

			//Début du déplacement
			ag.env.Lock()
			x := ag.position[0]
			y := ag.position[1]

			if len(ag.isOn) > 0 {
				// Suppression de l'agent
				ag.env.station[x][y] = ag.isOn[alg.Coord{x, y}]
				alg.RemoveCoord(alg.Coord{x, y}, ag.isOn)
			}

			// Enregistrement des valeurs précédentes de la matrice
			ag.isOn[alg.Coord{x, y}] = ag.env.station[x][y]

			//Ecriture agent dans la matrice (déplacement)
			ag.env.station[x][y] = string(ag.id)
			ag.env.Unlock()
			//Fin du déplacement

			// Prise en compte de la vitesse de déplacement
			time.Sleep(ag.vitesse * time.Millisecond)
			//TODELETEfmt.Printf("J'ai bougé")
			return true
		}
	}

	//si on peut aller nulle part, on se remet dans notre config initiale
	ag.direction = storeDirection
	//fmt.Printf("Je me bouge pas")
	return false
}

func (ag *Agent) MoveAgent() bool {
	//fmt.Printf("[MoveAgent, %s ] direction = %d \n",ag.id, ag.direction)
	// ================== Tentative de calcul du chemin =======================
	if len(ag.path) == 0 ||
		ag.isGoingToExitPath() ||
		(ag.env.station[ag.path[0].Row()][ag.path[0].Col()] == "O" && !alg.EqualCoord(&ag.destination, &alg.Coord{ag.path[0].Row(), ag.path[0].Col()})) {
		start, end := ag.generatePathExtremities()
		// Recherche d'un chemin si inexistant
		if len(ag.path) > 0 {
			ag.path = alg.FindPath(ag.env.station, start, end, ag.path[0], false, 2*time.Second)
		} else {
			ag.path = alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 2*time.Second)
		}
	}

	// ================== Vérification si déplacement possible =======================
	if ag.IsAgentBlocking() {
		//fmt.Printf("[MoveAgent, %s ] %s est bloqué\n",ag.id, ag.id)
		if ag.politesse {
			start, end := ag.generatePathExtremities()
			// Si un agent bloque notre déplacement, on attend un temps aléatoire, et reconstruit
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			//path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 2*time.Second)
			path := alg.FindPath(ag.env.station, start, end, ag.path[0], false, 2*time.Second)
			ag.path = path
			return false
		} else {
			//Si individu impoli, demande à l'agent devant de bouger
			//On récupère le id de la personne devant
			if existAgent(ag.NextCell()) {
				blockingAgentID := AgentID(ag.NextCell())
				//blockingAgent := ag.env.FindAgentByID(blockingAgentID)
				var reqToBlockingAgent *req.Request
				//var reqToImpoliteAgent *Request
				i := 0
				accept := false
				for !accept && i < 3 {
					//Demande à l'agent qui bloque de se pousser (réitère trois fois s'il lui dit pas possible)
					i += 1
					//TODELETEfmt.Printf("[MoveAgent, %s] You have to move %s for the %d time \n", ag.id, blockingAgentID, i)
					reqToBlockingAgent = req.NewRequest(ag.env.agentsChan[ag.id], YouHaveToMove) //Création "Hello, je suis ag.id, move."
					ag.env.agentsChan[blockingAgentID] <- *reqToBlockingAgent                    //Envoi requête
					repFromBlockingAgent := <-ag.env.agentsChan[ag.id]                           //Attend la réponse

					if repFromBlockingAgent.Decision() == Done { //BlockingAgent lui a répondu Done, il s'est donc poussé
						//TODELETEfmt.Printf("okay i will move agent %s \n", ag.id)
						accept = true
					}
				}
				if !accept {
					//TODELETEfmt.Printf("i can't move agent %s \n", ag.id)
					return false //il ne peut pas bouger, il s'arrête
				}
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
			//fmt.Println("[MoveAgent]Path : ", ag.path[0])
			ag.path = ag.path[1:]
		} else {
			ag.path = nil
		}
		ag.saveCells()
		ag.env.writeAgent(ag)
		// ============ Prise en compte de la vitesse de déplacement ======================
		time.Sleep(ag.vitesse * time.Millisecond)
		return true
	}
	return false
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
			//TODELETEfmt.Printf("[listenForRequests] Request received by :%s , decision : %d \n", ag.id, req.Decision())
			ag.request = &req
			if ag.request.Decision() == Expel {
				//TODELETEfmt.Println("[listenForRequests] Expel received by :", ag.id)
			}
			if ag.request != nil && (ag.request.Decision() == Disappear || ag.request.Decision() == EnterMetro) {
				return
			}
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

/*
 * Méthode qui envoie la valeur de case en face de l'agent
 */
func (ag *Agent) getFaceCase() string {
	switch {
	case ag.direction == 0: // vers le haut
		if (ag.position[0] - 1) < 0 {
			return "X" // si le controleur est au bord de la station, alors il fait face à un mur
		} else {
			return ag.env.station[ag.position[0]-1][ag.position[1]]
		}
	case ag.direction == 1: // vers la droite
		if (ag.position[1] + 1) > 49 {
			return "X" // si le controleur est au bord de la station, alors il fait face à un mur
		} else {
			return ag.env.station[ag.position[0]][ag.position[1]+1]
		}
	case ag.direction == 2: // vers le bas
		if (ag.position[0] + 1) > 49 {
			return "X" // si le controleur est au bord de la station, alors il fait face à un mur
		} else {
			return ag.env.station[ag.position[0]+1][ag.position[1]]
		}

	case ag.direction == 3: // vers la gauche
		if (ag.position[1] - 1) < 0 {
			return "X" // si le controleur est au bord de la station, alors il fait face à un mur
		} else {
			return ag.env.station[ag.position[0]][ag.position[1]-1]
		}
	}
	return "X"
}

func initDirection(depart alg.Coord, dimensionCarte int) int {
	n := rand.Intn(4) // direction aléatoire
	for !verifyDirection(n, depart, dimensionCarte) {
		n = rand.Intn(4) // direction aléatoire
	}
	return n
}

func verifyDirection(n int, depart alg.Coord, dimensionCarte int) bool {
	switch n {
	case 0: // vers le haut
		if (depart[0] - 1) < 0 {
			return false
		} else {
			return true
		}
	case 1: // vers la droite
		if (depart[1] + 1) > dimensionCarte {
			return false
		} else {
			return true
		}
	case 2: // vers le bas
		if (depart[0] + 1) > dimensionCarte {
			return false
		} else {
			return true
		}
	case 3: // vers la gauche
		if (depart[1] - 1) < 0 {
			return false
		} else {
			return true
		}
	}
	return false
}

/*============================FONCTION POUR TROUVER UN DESTINATION ============================*/

// Structure pour associer une Coord et sa distance par rapport au position d'un agent
type Gate struct {
	Position alg.Coord // Coordonnées de la porte
	Distance float64
	NbAgents float64
}

func (ag *Agent) findNearestExit_v0() alg.Coord {
	// Recherche de la sortie la plus proche
	nearest := alg.Coord{0, 0}
	min := 1000000
	//n := len(ag.env.station[0])
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			if ag.env.station[i][j] == "S" || ag.env.station[i][j] == "W" {
				dist := alg.Abs(ag.position[0]-i) + alg.Abs(ag.position[1]-j)
				if dist < min {
					min = dist
					nearest = alg.Coord{i, j}
				}
			}
		}
	}
	return nearest
}

func (ag *Agent) findNearestExit() alg.Coord {
	// Recherche de la sortie la plus proche
	sorties := ag.env.exits
	nearest := sorties[0]
	min := 1000000
	for _, sortie := range sorties {
		dist := alg.Abs(ag.position[0]-sortie[0]) + alg.Abs(ag.position[1]-sortie[1])
		if dist < min {
			min = dist
			nearest = sortie
		}
	}
	return nearest
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

func (env *Environment) GetAgentByChannel(channel chan req.Request) *Agent {
	env.RLock()
	defer env.RUnlock()

	for agentID, agentChannel := range env.agentsChan {
		if agentChannel == channel {
			return env.FindAgentByID(agentID)
		}
	}
	return nil
}

/* func (ag *Agent) AgentTakesNextPosition() bool {

	// ================ Récupère la force de l'agent =============
	ImpoliteAgent := ag.env.GetAgentByChannel(ag.request.demandeur)
	if ImpoliteAgent == nil {
		fmt.Printf("Channel not available")
	}
	ImpoliteAgentForce := ImpoliteAgent.force

	// ========== Vérifier la sécurité d'une direction ==========
	for or := 0; or < 4; or++ {
		safe, _ := IsMovementSafe(ag.path, ag, ag.env)
		if safe { //Si c'est safe dans une direction, on regarde jusqu'ou ça l'est encore avec la force
			}
		}
	}
	// Aucune direction n'est sûre, il ne bouge pas
	return false
} */
