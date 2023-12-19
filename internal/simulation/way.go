package simulation

/*
 * Classe et méthodes principales de la structure Way (porte de métro)
 */

type Way struct {
	id    WayID
	gates []Coord //listes des portes associée à la voie
}

type WayID int

func NewWay(wayId WayID, gates []Coord) *Way {
	return &Way{
		id:    wayId,
		gates: gates}
}
