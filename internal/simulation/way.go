package simulation

/*
 * Classe et méthodes principales de la structure Way (porte de métro)
 */

import (
	alg "metrosim/internal/algorithms"
	"time"
)

type Way struct {
	id             WayID
	upLeftCoord    alg.Coord // inclus
	downRightCoord alg.Coord // inclus
	goToLeft       bool      // si vrai, le métro se déplace de droite à gauche, si faux de gauche à droite
	horizontal     bool
	gates          []alg.Coord //listes des portes associée à la voie
	nearestExit    []alg.Coord // Chemin vers la sortie la plus proche pour chaque porte (index vers pathsToExit)
	pathsToExit    [][]alg.Node
	gatesClosed    bool
	env            *Environment
}

type WayID int

func NewWay(wayId WayID, upLeftCoord, downRightCoord alg.Coord, goToLeft bool, gates []alg.Coord, env *Environment) *Way {
	/* Affichage des portes */
	for _, gate := range gates {
		if !(gate[0] < 0 || gate[1] > 49) && env.station[gate[0]][gate[1]] != "X" && env.station[gate[0]][gate[1]] != "Q" {
			env.station[gate[0]][gate[1]] = "G"
		}

	}
	/* Sens de la voie */
	horizontal := true
	if alg.Abs(upLeftCoord[0]-downRightCoord[0]) > alg.Abs(upLeftCoord[1]-downRightCoord[1]) {
		horizontal = false
	}
	nearestExit := make([]alg.Coord, len(gates))
	pathsToExit := make([][]alg.Node, len(gates))
	for index, gate := range gates {
		row, col := alg.FindNearestExit(&env.exits, gate[0], gate[1])
		nearestExit[index] = alg.Coord{row, col}
		pathsToExit[index] = alg.FindPath(env.station, *alg.NewNode(gate[0], gate[1], 0, 0, 1, 1), *alg.NewNode(row, col, 0, 0, 0, 0), *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 5*time.Second)
		index++
	}

	return &Way{
		id:             wayId,
		upLeftCoord:    upLeftCoord,
		downRightCoord: downRightCoord,
		goToLeft:       goToLeft,
		horizontal:     horizontal,
		gates:          gates,
		nearestExit:    nearestExit,
		pathsToExit:    pathsToExit,
		env:            env}
}
