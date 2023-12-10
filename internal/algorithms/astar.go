package algorithms

import (
	"container/heap"
	"fmt"
)

/*
 * Utilisation de l'algorithme A* pour les déplacements
 * //TODO: Peut-être gérer un passage par référence et non par copie
 * 
 */
type Node struct {
	row, col, cost, heuristic, width, height, orientation int
}

func NewNode(row, col, cost, heuristic, width, height int) *Node {
	//fmt.Println()
	return &Node{row, col, cost, heuristic, width , height, 0}
}

func (nd *Node) Row() int {
	return nd.row
}

func (nd *Node) Col() int {
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

func getNeighbors(matrix [20][20]string, current, end Node, forbiddenCell Node) []*Node {
	fmt.Println("okkk")
	neighbors := make([]*Node, 0)

	// Possible moves: up, down, left, right, rotate (clockwise)
	possibleMoves := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, move := range possibleMoves {
		newRow, newCol := current.row+move[0], current.col+move[1]

		// Check if the new position is valid, considering agent dimensions and rotation
		if isValidMove(matrix, current, forbiddenCell, newRow, newCol) {
			neighbors = append(neighbors, &Node{
				row:         newRow,
				col:         newCol,
				cost:        current.cost + 1,
				heuristic:   heuristic(newRow, newCol, end),
				width:       current.width,
				height:      current.height,
				orientation: current.orientation,
			})
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
	for i := 0; i < current.height; i++ {
		for j := 0; j < current.width; j++ {
			// Calculate the rotated coordinates based on the agent's orientation
			rotatedI, rotatedJ := rotateCoordinates(i, j, current.orientation)

			// Calculate the absolute coordinates in the matrix
			absRow, absCol := newRow+rotatedI, newCol+rotatedJ

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
	// Rotate the coordinates based on the agent's orientation
	// You need to implement the logic for rotation based on your specific rules
	// This is a simple example that assumes the agent can rotate in all directions
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
