package simulation

type Way struct {
	id    WayID
	gates []Coord //listes des portes du métro
}

type WayID int
