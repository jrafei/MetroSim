package simulation

import (
	//"container/heap"
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
	vitesse             int
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
}

type Behavior interface {
	// TODO: Comment faire pour ne pas passer l'agent en param ? C'est possible ?
	Percept(*Agent, *Environment)
	Deliberate(*Agent)
	Act(*Agent, *Environment)
}

type UsagerLambda struct{}

func (ul *UsagerLambda) Percept(ag *Agent, env *Environment) {
	// TODO: Essayer un nouveau chemin quand l'agent est bloqué
}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	if ag.env.station[ag.departure[0]][ag.departure[1]] == "A" {
		ag.decision = Wait
	} else {
		ag.decision = Move
	}
}

func (ul *UsagerLambda) Act(ag *Agent, env *Environment) {
	// TODO: Je crois que la construction d'un chemin s'arrête s'il y a déjà un agent sur destination. Il faudrait donc faire en sorte de s'approcher le plus possible
	if ag.decision == Move {
		start := ag.coordBasOccupation
		end := ag.destination
		path := findPathBFS(ag.env.station, start, end)
		if len(path) > 0 && env.station[path[0][0]][path[0][1]] != "A" { // TODO: Pas ouf les conditions je trouve
			ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = ag.isOn
			ag.isOn = ag.env.station[path[0][0]][path[0][1]]
			ag.coordBasOccupation[0] = path[0][0]
			ag.coordBasOccupation[1] = path[0][1]
			ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = "A"
		}
		//vitesseInSeconds := int(ag.vitesse)
		// Multiply the vitesse by time.Second
		//sleepDuration := time.Duration(vitesseInSeconds) * time.Second
		time.Sleep(200 * time.Millisecond)
	}
	if ag.decision == Wait {
		n := rand.Intn(1) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
	}

}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse int, force int, politesse bool, UpCoord Coord, DownCoord Coord, behavior Behavior, departure, destination Coord) *Agent {
	return &Agent{AgentID(id), vitesse, force, politesse, UpCoord, DownCoord, departure, destination, behavior, env, syncChan, Noop, env.station[UpCoord[0]][UpCoord[1]]}
}

func (ag *Agent) ID() AgentID {
	return ag.id
}

func (ag *Agent) Start() {
	log.Printf("%s starting...\n", ag.id)

	go func() {
		env := ag.env
		var step int
		for {
			step = <-ag.syncChan
			ag.behavior.Percept(ag, env)
			ag.behavior.Deliberate(ag)
			ag.behavior.Act(ag, env)
			ag.syncChan <- step
		}
	}()
}

func (ag *Agent) Percept(env *Environment) {
	//ag.rect = env.Rect()
}

func (ag *Agent) Deliberate() {
	if rand.Float64() < 0.1 {
		ag.decision = Noop
	} else {
		ag.decision = Mark
	}
}

func (ag *Agent) Act(env *Environment) {
	if ag.decision == Noop {
		env.Do(Noop, Coord{})
	}
}

// ================================================================================

/*
 * Utilisation de l'algorithme BFS pour les déplacements. Beaucoup plus rapide que A*
 *
 *
 */
/*
func findPathBFS(matrix [20][20]string, start, end Coord) []Coord {
	queue := []Coord{start}
	visited := make(map[Coord]bool)
	parents := make(map[Coord]Coord)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			// Construire le chemin à partir des parents
			path := []Coord{current}
			for parent, ok := parents[current]; ok; parent, ok = parents[parent] {
				path = append([]Coord{parent}, path...)
			}
			return path[1:]
		}

		visited[current] = true

		neighbors := getNeighborsBFS(matrix, current)
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				parents[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	return nil // Aucun chemin trouvé
}
*/

func findPathBFS(matrix [20][20]string, start, end Coord) []Coord {
	queue := []Coord{start}
	visited := make(map[Coord]bool)
	parents := make(map[Coord]Coord)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			// Construire le chemin à partir des parents
			path := []Coord{current}
			for parent, ok := parents[current]; ok; parent, ok = parents[parent] {
				path = append([]Coord{parent}, path...)
			}
			return path[1:]
		}

		visited[current] = true

		neighbors := getNeighborsBFS(matrix, current)
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				parents[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	return nil // Aucun chemin trouvé
}

func getNeighborsBFS(matrix [20][20]string, current Coord) []Coord {
	neighbors := make([]Coord, 0)

	// Déplacements possibles : haut, bas, gauche, droite
	possibleMoves := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, move := range possibleMoves {
		newRow, newCol := current[0]+move[0], current[1]+move[1]

		// Vérifier si la nouvelle position est valide et non visitée
		if newRow >= 0 && newRow < len(matrix) && newCol >= 0 && newCol < len(matrix[0]) && (matrix[newRow][newCol] != "Q" && matrix[newRow][newCol] != "X") {
			neighbors = append(neighbors, Coord{newRow, newCol})
		}
	}

	return neighbors
}

/*
type Node struct {
	row, col, cost, heuristic int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return (pq[i].cost + pq[i].heuristic) < (pq[j].cost + pq[j].heuristic)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func findPath(matrix [20][20]string, start, end Node) []Node {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &start)
	visited := make(map[Node]bool)
	parents := make(map[Node]Node)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Node)

		if current.row == end.row && current.col == end.col {
			// Construire le chemin à partir des parents
			path := []Node{*current}
			for parent, ok := parents[*current]; ok; parent, ok = parents[parent] {
				path = append([]Node{parent}, path...)
			}

			return path[1:]
		}

		visited[*current] = true

		neighbors := getNeighbors(matrix, *current)
		for _, neighbor := range neighbors {
			if !visited[*neighbor] {
				parents[*neighbor] = *current
				heap.Push(&pq, neighbor)
			}
		}
	}

	return nil // Aucun chemin trouvé
}

func getNeighbors(matrix [20][20]string, current Node) []*Node {
	neighbors := make([]*Node, 0)

	// Déplacements possibles : haut, bas, gauche, droite
	possibleMoves := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, move := range possibleMoves {
		newRow, newCol := current.row+move[0], current.col+move[1]

		// Vérifier si la nouvelle position est valide et non visitée
		if newRow >= 0 && newRow < len(matrix) && newCol >= 0 && newCol < len(matrix[0]) && (matrix[newRow][newCol] == "_" || matrix[newRow][newCol] == "B") {
			neighbors = append(neighbors, &Node{newRow, newCol, current.cost + 1, heuristic(newRow, newCol, current, current)})
		}
	}

	return neighbors
}

func heuristic(row, col int, current, end Node) int {
	// Heuristique simple : distance de Manhattan
	return abs(row-end.row) + abs(col-end.col)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
*/
