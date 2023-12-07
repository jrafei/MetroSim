package algorithms

import (
	"container/heap"
)

/*
 * Utilisation de l'algorithme A* pour les déplacements
 *
 *
 */
type Node struct {
	row, col, cost, heuristic int
}

func NewNode(row, col, cost, heuristic int) *Node{
	return &Node{row, col ,cost , heuristic}
}

func (nd *Node) Row() int{
	return nd.row
}

func (nd *Node) Col() int{
	return nd.col
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

func FindPath(matrix [20][20]string, start, end Node, forbidenCell Node) []Node {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &start)
	visited := make(map[Node]bool)
	parents := make(map[Node]Node)

	closestPoint := start // Initialisation avec le point de départ
	closestDistance := heuristic(start.row, start.col, end)

	foundPath := false

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Node)

		// Mise à jour du point le plus proche si le point actuel est plus proche
		currentDistance := heuristic(current.row, current.col, end)
		if currentDistance < closestDistance {
			closestPoint = *current
			closestDistance = currentDistance
		}

		if current.row == end.row && current.col == end.col {
			// Construire le chemin à partir des parents
			path := []Node{closestPoint}
			for parent, ok := parents[closestPoint]; ok; parent, ok = parents[parent] {
				path = append([]Node{parent}, path...)
			}

			return path[1:]
		}

		visited[*current] = true

		neighbors := getNeighbors(matrix, *current, end, forbidenCell)
		for _, neighbor := range neighbors {
			if !visited[*neighbor] {
				parents[*neighbor] = *current
				heap.Push(&pq, neighbor)
			}
		}
		foundPath = true
	}

	if foundPath {
		// Retourner le chemin le plus proche même si la destination n'a pas été atteinte
		path := []Node{closestPoint}
		for parent, ok := parents[closestPoint]; ok; parent, ok = parents[parent] {
			path = append([]Node{parent}, path...)
		}
		return path[1:]
	}

	return nil // Aucun chemin trouvé
}

func getNeighbors(matrix [20][20]string, current, end Node, forbidenCell Node) []*Node {
	neighbors := make([]*Node, 0)

	// Déplacements possibles : haut, bas, gauche, droite
	possibleMoves := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, move := range possibleMoves {
		newRow, newCol := current.row+move[0], current.col+move[1]

		// Vérifier si la nouvelle position est valide
		if (forbidenCell.row != newRow || forbidenCell.col != newCol) && newRow >= 0 && newRow < len(matrix) && newCol >= 0 && newCol < len(matrix[0]) && (matrix[newRow][newCol] != "Q" && matrix[newRow][newCol] != "X") {
			neighbors = append(neighbors, &Node{newRow, newCol, current.cost + 1, heuristic(newRow, newCol, end)})
		}
	}

	return neighbors
}

func heuristic(row, col int, end Node) int {
	// Heuristique simple : distance de Manhattan
	return abs(row-end.row) + abs(col-end.col)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
