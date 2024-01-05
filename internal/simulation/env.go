package simulation

import (
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"sync"
)

//TODO:rajouter les entrées et sorties

type Environment struct {
	sync.RWMutex
	ags              []Agent
	agentCount       int
	station          [50][50]string
	agentsChan       map[AgentID]chan req.Request
	controlledAgents map[AgentID]bool
	newAgentChan     chan Agent
	metros           []Metro
	exits            []alg.Coord
	entries          []alg.Coord
	gates            []alg.Coord
}

func NewEnvironment(ags []Agent, carte [50][50]string, metros []Metro, newAgtCh chan Agent, agtCount int, sim *Simulation) (env *Environment) {

	agentsCh := make(map[AgentID]chan req.Request)
	mapControlled := make(map[AgentID]bool)

	// Récupération des entrées et sorties
	entries := make([]alg.Coord, 0)
	exits := make([]alg.Coord, 0)
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			switch carte[i][j] {
			case "E":
				entries = append(entries, alg.Coord{i, j})
			case "S":
				exits = append(exits, alg.Coord{i, j})
			case "W":
				entries = append(entries, alg.Coord{i, j})
				exits = append(exits, alg.Coord{i, j})
			}
		}
	}

	// Récupération des portes

	gates := make([]alg.Coord, 0)
	for _, metro := range metros {
		//TODELETEfmt.Println(metro.way.gates)
		for _, gate := range metro.way.gates {
			gates = append(gates, gate)
		}
	}

	// Initialisation de la vérification du contrôle agents (faux au départ)
	for _, ag := range ags {
		mapControlled[ag.id] = false
	}

	return &Environment{
		ags:              ags,
		agentCount:       agtCount,
		station:          carte,
		agentsChan:       agentsCh,
		controlledAgents: mapControlled,
		newAgentChan:     newAgtCh,
		exits:            exits,
		entries:          entries,
		gates:            gates,
		metros:           metros,
	}
}

func (env *Environment) AddAgent(agt Agent) {
	// Ajout d'un agent à l'environnement
	env.Lock()
	defer env.Unlock()
	env.ags = append(env.ags, agt)
	env.controlledAgents[agt.id] = false
	// ajout du channel de l'agent à l'environnement
	env.agentsChan[agt.id] = make(chan req.Request, 5)
	env.agentCount++
	env.newAgentChan <- agt

}

func (env *Environment) DeleteAgent(agt Agent) {
	// Suppression d'un agent de l'environnement
	env.Lock()
	defer env.Unlock()
	for i := 0; i < len(env.station); i++ {
		if env.ags[i].id == agt.id {
			// Utiliser la syntaxe de découpage pour supprimer l'élément
			env.ags = append(env.ags[:i], env.ags[i+1:]...)
			delete(env.agentsChan, agt.id)
			// Sortir de la boucle après avoir trouvé et supprimé l'élément
			break
		}
	}

}

func (env *Environment) GetAgentChan(agt_id AgentID) chan req.Request {
	// Récupère le channel de communication d'un agent
	return env.agentsChan[agt_id]
}

func (env *Environment) FindAgentByID(agtId AgentID) *Agent {
	for i := range env.ags {
		if env.ags[i].id == agtId {
			return &env.ags[i]
		}
	}
	return nil
}

func (env *Environment) verifyEmptyCase(c alg.Coord) bool {
	return env.station[c[0]][c[1]] == "_" //|| env.station[c[0]][c[1]] == "E" || env.station[c[0]][c[1]] == "S" || env.station[c[0]][c[1]] == "W"
}

func existAgent(c string) bool {
	// Vérifie si c'est un agent
	return c != "X" && c != "E" && c != "S" && c != "W" && c != "Q" && c != "_" && c != "B"
}

func calculDirection(depart alg.Coord, arrive alg.Coord) int {
	// Calcul de la direction d'un agent
	if depart[0] == arrive[0] {
		if depart[1] > arrive[1] {
			return 3 //Gauche
		} else {
			return 1 //droite
		}
	} else {
		if depart[0] > arrive[0] {
			return 0 //haut
		} else {
			return 2 //bas
		}
	}
}

func (env *Environment) RemoveAgent(agt *Agent) {
	// Supprime l'agent de la matrice

	// Calcul des bornes de position de l'agent
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := alg.CalculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			if len(env.station[i][j]) > 1 {
				env.station[i][j] = agt.isOn[alg.Coord{i, j}]
			}
			alg.RemoveCoord(alg.Coord{i, j}, agt.isOn)

		}
	}
}

func (env *Environment) writeAgent(agt *Agent) {
	// Ecris l'agent dans la matrice

	env.Lock()
	defer env.Unlock()

	// Calcul des bornes de position de l'agent
	borneInfRow, borneSupRow, borneInfCol, borneSupCol := alg.CalculateBounds(agt.position, agt.width, agt.height, agt.orientation)

	for i := borneInfRow; i < borneSupRow; i++ {
		for j := borneInfCol; j < borneSupCol; j++ {
			env.station[i][j] = string(agt.id)
		}
	}

}

func (env *Environment) getNbAgentsAround(pos alg.Coord) int {
	//pos est la position de la porte
	way := env.getWay(pos)
	upleft := way.upLeftCoord
	//downright := way.downRightCoord

	upmetro := false
	if pos[0] == upleft[0]-1 { //si la porte est en haut du métro
		upmetro = true
	}

	nb := 0
	if upmetro {
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				a := pos[0] - i
				b := pos[1] + j
				if a >= 0 && b >= 0 && a < len(env.station[0]) && b < len(env.station[1]) {
					c := env.station[a][b]
					if existAgent(c) {
						nb++
					}
				}
				b = pos[1] - j
				if a >= 0 && b >= 0 && a < len(env.station[0]) && b < len(env.station[1]) {
					c := env.station[a][b]
					if existAgent(c) {
						nb++
					}
				}
			}
		}
	} else { //si la porte est en bas du métro
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				a := pos[0] + i
				b := pos[1] + j
				if a >= 0 && b >= 0 && a < len(env.station[0]) && b < len(env.station[1]) {
					c := env.station[a][b]
					if existAgent(c) {
						nb++
					}
				}
				b = pos[1] - j
				if a >= 0 && b >= 0 && a < len(env.station[0]) && b < len(env.station[1]) {
					c := env.station[a][b]
					if existAgent(c) {
						nb++
					}
				}
			}
		}
	}
	return nb
}

func (env *Environment) getWay(pos alg.Coord) *Way {
	for _, metro := range env.metros {
		for _, gate := range metro.way.gates {
			if gate[0] == pos[0] && gate[1] == pos[1] {
				return metro.way
			}
		}
	}
	return nil
}

func (env *Environment) Station() [50][50]string {
	// Récupère la matrice de l'environnement
	env.RLock()
	defer env.RUnlock()
	return env.station
}

func (env *Environment) getNearGateFromGate(gate Gate) []alg.Coord {
	/*
	 * Renvoie les portes proches de la porte passée en paramètre
	 */
	gates := env.gates
	nearGates := make([]alg.Coord, 0)
	for _, g := range gates {
		if alg.Abs(g[0]-gate.Position[0]) <= 1 && alg.Abs(g[1]-gate.Position[1]) <= 1 {
			nearGates = append(nearGates, g)
		}
	}
	nearGates = append(nearGates, gate.Position)
	return nearGates
}
