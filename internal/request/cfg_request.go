package request

import (
	alg "metrosim/internal/algorithms"
)

type Cfg struct {
	// Section des paramètres généraux
	Station     [50][50]string `json:"map"`
	MaxDuration int            `json:"maxDuration"` // durée maximale de la simulation (en secondes)
	Flow        int            `json:"flow"`        // débit de création d'agents/seconde, 1 agent créé toutes les flow MILLIsecondes
	Controleurs bool           `json:"controleurs"` // présence de controleurs
	Fraudeurs   bool           `json:"fraudeurs"`   // présence de fraudeurs
	Impolis     bool           `json:"impolis"`     // présence d'impolis
	Mob_reduite bool           `json:"mob_reduite"` // présence de personnes à mobilité réduite
	Patients    bool           `json:"patients"`    // présence d'usagers patients

	// Section des paramètres des voies
	// (un indice est associé à une voie, ex: goToLeft[0]
	// concerne la voie 0, goToLeft[1] concerne la voie 1, etc...)
	LeftTopCorners     []alg.Coord   `json:"leftTopCorners"`   // coins haut-gauche des voies
	RightBottomCorners []alg.Coord   `json:"rightDownCorners"` // coins bas-droit des voies
	GoToLeft           []bool        `json:"goToLeft"`         // sens de circulation des voies
	Gates              [][]alg.Coord `json:"gates"`            // coordonnées des portes (gates[i] correspond aux portes de la voie i)

	// Section des paramètres des trains (un indice est associé à un train, et un train à le même indice que sa voie ex : capacity[0] correspond à la capacité du train de la voie 0])
	Frequency []int `json:"frequency"` // fréquence de passage des trains en secondes
	StopTime  []int `json:"stopTime"`  // temps d'arrêt des trains en secondes
	Capacity  []int `json:"capacity"`  // capacité des trains (nombre de cases occupées par un train)

}
