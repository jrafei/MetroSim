package algorithms

type Coord [2]int

func EqualCoord(coord1, coord2 *Coord) bool {
	// Vérifie l'égalité de 2 objets Coord
	return coord1[0] == coord2[0] && coord1[1] == coord2[1]
}

func RemoveCoord(to_remove Coord, mapping map[Coord]string) {
	// Suppression d'une clé dans une map
	for coord, _ := range mapping {
		if EqualCoord(&coord, &to_remove) {
			delete(mapping, coord)
		}
	}
}


func CalculateBounds(position Coord, width, height, orientation int) (infRow, supRow, infCol, supCol int) {
	// Fonction de génération des frontières d'un objet ayant une largeur et une hauteur, en focntion de son orientation
	borneInfRow := 0
	borneSupRow := 0
	borneInfCol := 0
	borneSupCol := 0

	// Calcul des bornes de position de l'agent après mouvement
	switch orientation {
	case 0:
		// Orienté vers le haut
		borneInfRow = position[0] - width + 1
		borneSupRow = position[0] + 1
		borneInfCol = position[1]
		borneSupCol = position[1] + height
	case 1:
		// Orienté vers la droite
		borneInfRow = position[0]
		borneSupRow = position[0] + height
		borneInfCol = position[1]
		borneSupCol = position[1] + width
	case 2:
		// Orienté vers le bas
		borneInfRow = position[0]
		borneSupRow = position[0] + width
		borneInfCol = position[1]
		borneSupCol = position[1] + height
	case 3:
		// Orienté vers la gauche
		borneInfRow = position[0]
		borneSupRow = position[0] + height
		borneInfCol = position[1] - width + 1
		borneSupCol = position[1] + 1

	}
	return borneInfRow, borneSupRow, borneInfCol, borneSupCol
}