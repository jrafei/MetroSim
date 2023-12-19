package simulation

/*
 * Classe et méthodes principales de la structure Way (porte de métro)
 */

type Way struct {
	id             WayID
	upLeftCoord    Coord   // inclus
	downRightCoord Coord   // inclus
	goToLeft       bool    // si vrai, le métro se déplace de droite à gauche, si faux de gauche à droite
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
	return &Way{
		id:             wayId,
		upLeftCoord:    upLeftCoord,
		downRightCoord: downRightCoord,
		goToLeft:       goToLeft,
		gates:          gates,
		env:            env}
}
