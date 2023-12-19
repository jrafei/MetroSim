package simulation

/*
 * Classe et méthodes principales de la structure Way (porte de métro)
 */

import (
	alg "metrosim/internal/algorithms"
)

type Way struct {
	id             WayID
	upLeftCoord    Coord // inclus
	downRightCoord Coord // inclus
	goToLeft       bool  // si vrai, le métro se déplace de droite à gauche, si faux de gauche à droite
	horizontal     bool
	gates          []Coord //listes des portes associée à la voie
	env            *Environment
}

type WayID int

func NewWay(wayId WayID, upLeftCoord, downRightCoord Coord, goToLeft bool, gates []Coord, env *Environment) *Way {
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
	return &Way{
		id:             wayId,
		upLeftCoord:    upLeftCoord,
		downRightCoord: downRightCoord,
		goToLeft:       goToLeft,
		horizontal:     horizontal,
		gates:          gates,
		env:            env}
}
