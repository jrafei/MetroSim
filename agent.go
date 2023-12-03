package simulation

import (
	//"container/heap"
	"fmt"
	"log"
	"math"
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
	stuck               bool
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

	// Perception des éléments autour de l'agent pour déterminer si bloqué
	s := 0 // nombre de cases indisponibles autour de l'agent
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			ord := (ag.coordBasOccupation[0] - 1) + i
			abs := (ag.coordBasOccupation[1] - 1) + j
	
			if(ord!=ag.coordBasOccupation[0] && abs !=ag.coordBasOccupation[1]){
				if ord < 0 || abs < 0 || ord > 19 || abs > 19 {
					s++
				} else if env.station[ord][abs] == "X" || env.station[ord][abs] == "Q" || env.station[ord][abs] == "A" {
					s++
				}
			}
		}
	}
	// Si pas de case disponible autour de lui, il est bloqué
	// (ag.env.station[ag.departure[0]][ag.departure[1]] == "A" && ag.coordBasOccupation[0] == ag.departure[0] && ag.coordBasOccupation[1] == ag.departure[1])
	if s == 8  {
		ag.stuck = true
		fmt.Println(ag.id, ag.stuck, ag.coordBasOccupation, s)
	} else {
		ag.stuck = false
	}
}

func (ul *UsagerLambda) Deliberate(ag *Agent) {
	// Si l'agent est bloqué, il doit attendre qu'une case se libère autour de lui
	if ag.stuck {
		ag.decision = Wait
	} else {
		ag.decision = Move
	}
}

func (ul *UsagerLambda) Act(ag *Agent, env *Environment) {
	if ag.decision == Move {

		start := ag.coordBasOccupation
		end := ag.destination
		_, path := findClosestPointBFS(ag.env.station, start, end)

		if len(path) > 0 && env.station[path[0][0]][path[0][1]] != "A" { // TODO: Pas ouf les conditions je trouve
			ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = ag.isOn
			ag.isOn = ag.env.station[path[0][0]][path[0][1]]
			ag.coordBasOccupation[0] = path[0][0]
			ag.coordBasOccupation[1] = path[0][1]
			ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = "A"
		}
		//vitesseInSeconds := int(ag.vitesse)
		//sleepDuration := time.Duration(vitesseInSeconds) * time.Second
		time.Sleep(200 * time.Millisecond)
	}
	if ag.decision == Wait {
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	}

}

func NewAgent(id string, env *Environment, syncChan chan int, vitesse int, force int, politesse bool, UpCoord Coord, DownCoord Coord, behavior Behavior, departure, destination Coord) *Agent {
	return &Agent{AgentID(id), vitesse, force, politesse, UpCoord, DownCoord, departure, destination, behavior, env, syncChan, Noop, env.station[UpCoord[0]][UpCoord[1]], false}
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

func distance(coord1, coord2 Coord) float64 {
	dx := float64(coord1[0] - coord2[0])
	dy := float64(coord1[1] - coord2[1])
	return math.Sqrt(dx*dx + dy*dy)
}

func findClosestPointBFS(matrix [20][20]string, start, end Coord) (Coord, []Coord) {
	queue := []Coord{start}
	visited := make(map[Coord]bool)
	parents := make(map[Coord]Coord)
	closestPoint := start // Initialisez avec le point de départ
	closestDistance := distance(start, end)
	foundPath := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Mettez à jour le point le plus proche si le point actuel est plus proche
		currentDistance := distance(current, end)
		if currentDistance < closestDistance {
			closestPoint = current
			closestDistance = currentDistance
		}

		if current == end {
			// Construire le chemin du point le plus proche à la destination
			path := []Coord{closestPoint}
			for parent, ok := parents[closestPoint]; ok; parent, ok = parents[parent] {
				path = append([]Coord{parent}, path...)
			}
			return closestPoint, path[1:]
		}

		visited[current] = true

		neighbors := getNeighborsBFS(matrix, current, end)
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				parents[neighbor] = current
				queue = append(queue, neighbor)
			}
		}

		foundPath = true
	}

	if foundPath {
		// Retourner le chemin le plus proche même si la destination n'a pas été atteinte
		path := []Coord{closestPoint}
		for parent, ok := parents[closestPoint]; ok; parent, ok = parents[parent] {
			path = append([]Coord{parent}, path...)
		}
		return closestPoint, path[1:]
	}

	return closestPoint, nil // Aucun chemin trouvé
}

func getNeighborsBFS(matrix [20][20]string, current Coord, end Coord) []Coord {
	neighbors := make([]Coord, 0)

	// Déplacements possibles : haut, bas, gauche, droite
	possibleMoves := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, move := range possibleMoves {
		newRow, newCol := current[0]+move[0], current[1]+move[1]

		// Vérifier si la nouvelle position est valide et non visitée
		if newRow >= 0 && newRow < len(matrix) && newCol >= 0 && newCol < len(matrix[0]) && (matrix[newRow][newCol] != "Q" && matrix[newRow][newCol] != "X") {
			if !(matrix[newRow][newCol] == "A" && newRow != end[0] && newCol != end[1]) {
				// Si la case du chemin ne comporte pas d'agent et que ce n'est pas la case d'arrivée, on peut l'ajouter
				neighbors = append(neighbors, Coord{newRow, newCol})
			}

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
