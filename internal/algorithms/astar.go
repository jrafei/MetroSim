package algorithms

import (
	"container/heap"
)

/*
 * Utilisation de l'algorithme A* pour les déplacements
 * //TODO: Peut-être gérer un passage par référence et non par copie
 * //TODO: faire des points de repère
 */
type Node struct {
	row, col, cost, heuristic, width, height, orientation int
}

func NewNode(row, col, cost, heuristic, width, height int) *Node {
	//fmt.Println()
	return &Node{row, col, cost, heuristic, width, height, 0}
}

func (nd *Node) Row() int {
	return nd.row
}

func (nd *Node) Col() int {
	return nd.col
}

func (nd *Node) Or() int {
	return nd.orientation
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
type  ZoneID int
type Coord [2]int

func FindPath(matrix [20][20]string, start, end Node, forbidenCell Node) []Node {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &start)
	visited := make(map[Node]bool)
	parents := make(map[Node]Node)

	closestPoint := start // Initialisation avec le point de départ
	closestDistance := Heuristic(start.row, start.col, end)

	foundPath := false

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Node)

		// Mise à jour du point le plus proche si le point actuel est plus proche
		currentDistance := Heuristic(current.row, current.col, end)
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

func getNeighbors(matrix [20][20]string, current, end Node, forbiddenCell Node) []*Node {
	//fmt.Println("okk")
	neighbors := make([]*Node, 0)

	// Possible moves: up, down, left, right
	possibleMoves := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, move := range possibleMoves {
		newRow, newCol := current.row+move[0], current.col+move[1]
		for orientation := 0; orientation < 4; orientation++ {
			current.orientation = orientation
			// fmt.Println(orientation)
			// Check if the new position is valid, considering agent dimensions and rotation
			if isValidMove(matrix, current, forbiddenCell, newRow, newCol) {
				neighbors = append(neighbors, &Node{
					row:         newRow,
					col:         newCol,
					cost:        current.cost + 1,
					heuristic:   Heuristic(newRow, newCol, end),
					width:       current.width,
					height:      current.height,
					orientation: current.orientation,
				})
			}
		}

	}

	return neighbors
}

func Heuristic(row, col int, end Node) int {
	// Heuristique simple : distance de Manhattan
	return abs(row-end.row) + abs(col-end.col)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isValidMove(matrix [20][20]string, current Node, forbiddenCell Node, newRow, newCol int) bool {
	// Check if the new position is within the bounds of the matrix
	if newRow < 0 || newRow >= len(matrix) || newCol < 0 || newCol >= len(matrix[0]) {
		return false
	}

	// Check if the new position overlaps with forbidden cells or obstacles
	if forbiddenCell.row == newRow && forbiddenCell.col == newCol {
		return false
	}

	// Check if the agent fits in the new position, considering its dimensions and rotation
	lRowBound, uRowBound, lColBound, uColBound := calculateBounds(newRow, newCol, current.width, current.height, current.orientation)

	for i := lRowBound; i < uRowBound; i++ {
		for j := lColBound; j < uColBound; j++ {

			// Calculate the absolute coordinates in the matrix
			absRow, absCol := i, j

			// Check if the absolute coordinates are within the bounds of the matrix
			if absRow < 0 || absRow >= len(matrix) || absCol < 0 || absCol >= len(matrix[0]) {
				return false
			}

			// Check if the absolute coordinates overlap with forbidden cells or obstacles
			if forbiddenCell.row == absRow && forbiddenCell.col == absCol {
				return false
			}

			// Check if the absolute coordinates overlap with obstacles in the matrix
			if matrix[absRow][absCol] == "Q" || matrix[absRow][absCol] == "X" {
				return false
			}
		}
	}

	return true
}

func rotateCoordinates(i, j, orientation int) (rotatedI, rotatedJ int) {

	switch orientation {
	case 0: // No rotation
		rotatedI, rotatedJ = i, j
	case 1: // 90-degree rotation
		rotatedI, rotatedJ = j, -i
	case 2: // 180-degree rotation
		rotatedI, rotatedJ = -i, -j
	case 3: // 270-degree rotation
		rotatedI, rotatedJ = -j, i
	}

	return rotatedI, rotatedJ
}

func calculateBounds(row, col, width, height, orientation int) (infRow, supRow, infCol, supCol int) {
	borneInfRow := 0
	borneSupRow := 0
	borneInfCol := 0
	borneSupCol := 0

	// Calcul des bornes de position de l'agent après mouvement
	switch orientation {
	case 0:
		borneInfRow = row - width + 1
		borneSupRow = row + 1
		borneInfCol = col
		borneSupCol = col + height
	case 1:
		borneInfRow = row
		borneSupRow = row + height
		borneInfCol = col
		borneSupCol = col + width
	case 2:
		borneInfRow = row
		borneSupRow = row + width
		borneInfCol = col
		borneSupCol = col + height
	case 3:
		borneInfRow = row
		borneSupRow = row + height
		borneInfCol = col - width + 1
		borneSupCol = col + 1

	}
	return borneInfRow, borneSupRow, borneInfCol, borneSupCol
}
