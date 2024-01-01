package algorithms

import (
	"container/heap"
	"context"
	"math/rand"
	"time"
)

/*
 * Utilisation de l'algorithme A* pour les déplacements
 * //TODO: Peut-être gérer un passage par référence et non par copie
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

func (nd *Node) Heuristic() int {
	return nd.heuristic
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

func FindPath(matrix [50][50]string, start, end Node, forbidenCell Node, orientation bool, timeout time.Duration) []Node {
	// Création d'un context avec timeout, pour limiter le calcul
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &start)
	visited := make(map[Node]bool)
	parents := make(map[Node]Node)

	closestPoint := start // Initialisation avec le point de départ
	closestDistance := Heuristic(matrix, start.row, start.col, end)

	foundPath := false

	for pq.Len() > 0 {
		select {
		case <-ctx.Done():
			// Timeout reached, return an error or handle accordingly
			return nil
		default:
			// Continue with the algorithm
		}
		current := heap.Pop(&pq).(*Node)

		// Mise à jour du point le plus proche si le point actuel est plus proche
		currentDistance := Heuristic(matrix, current.row, current.col, end)
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

		neighbors := getNeighbors(matrix, *current, end, forbidenCell, orientation)
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

func getNeighbors(matrix [50][50]string, current, end Node, forbiddenCell Node, orientation bool) []*Node {
	neighbors := make([]*Node, 0)

	// Déplacements possibles: up, down, left, right
	possibleMoves := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, move := range possibleMoves {
		newRow, newCol := current.row+move[0], current.col+move[1]
		if orientation {
			for or := 0; or < 4; or++ {
				current.orientation = or
				//Vérifie que le déplacement soit valide
				if isValidMove(matrix, current, forbiddenCell, newRow, newCol, orientation) {
					neighbors = append(neighbors, &Node{
						row:         newRow,
						col:         newCol,
						cost:        current.cost + 1,
						heuristic:   Heuristic(matrix, newRow, newCol, end),
						width:       current.width,
						height:      current.height,
						orientation: current.orientation,
					})
				}
			}

		} else {
			if isValidMove(matrix, current, forbiddenCell, newRow, newCol, orientation) {
				neighbors = append(neighbors, &Node{
					row:         newRow,
					col:         newCol,
					cost:        current.cost + 1,
					heuristic:   Heuristic(matrix, newRow, newCol, end),
					width:       current.width,
					height:      current.height,
					orientation: current.orientation,
				})
			}
		}

	}

	return neighbors
}

func Heuristic(matrix [50][50]string, row, col int, end Node) int {
	// Heuristique simple : distance de Manhattan
	// On introduit de l'aléatoire pour ajouter de la diversité dans la construction des chemins
	// On évite d'avoir tout le temps le même chemin pour un même point de départ et d'arrivée
	//return abs(row-end.row) + abs(col-end.col) + rand.Intn(3)
	malus := 0
	if len(matrix[row][col]) > 1 {
		malus += 10
	}
	return Abs(row-end.row) + Abs(col-end.col) + rand.Intn(10) + malus
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isValidMove(matrix [50][50]string, current Node, forbiddenCell Node, newRow, newCol int, orientation bool) bool {
	// Check if the new position is within the bounds of the matrix
	if newRow < 0 || newRow >= len(matrix) || newCol < 0 || newCol >= len(matrix[0]) {
		return false
	}

	// Check if the new position overlaps with forbidden cells or obstacles
	if forbiddenCell.row == newRow && forbiddenCell.col == newCol {
		//return false
		current.heuristic = current.heuristic + 100
	}
	// Check if the absolute coordinates overlap with obstacles in the matrix
	if matrix[newRow][newCol] == "Q" || matrix[newRow][newCol] == "X" || matrix[newRow][newCol] == "M" {
		return false
	}

	// Check if the agent fits in the new position, considering its dimensions and rotation
	if orientation {
		lRowBound, uRowBound, lColBound, uColBound := CalculateBounds(Coord{newRow, newCol}, current.width, current.height, current.orientation)

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
					//return false
					current.heuristic = current.heuristic + 100
				}

				// Check if the absolute coordinates overlap with obstacles in the matrix
				if matrix[absRow][absCol] == "Q" || matrix[absRow][absCol] == "X" || matrix[absRow][absCol] == "M" {
					return false
				}
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

func FindNearestExit(exits *[]Coord, row, col int) (dest_row, dest_col int) {
	// Recherche de la sortie la plus proche
	min := 1000000
	for _, exit := range *exits {
		dist := Abs(row-exit[0]) + Abs(col-exit[1])
		if dist < min {
			min = dist
			dest_row = exit[0]
			dest_col = exit[1]
		}
	}

	return dest_row, dest_col
}
