package simulation

import (
	"fmt"
	"log"
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"sync"
	"time"
)

/*
 * //TODO: Mettre en place un débit d'apparition des agents
 */

// Déclaration de la matrice
/*
  - X : Mur, zone inatteignable
  - E : Entrée
  - S : Sortie
  - W : Entrée et Sortie
  - Q : Voie
  - _ : Couloir, case libre
  - B: Bridge/Pont, zone accessible
  - G: gate/porte de métro
  - O : Porte ouverte
  - M : rame de métro
  - valeur de AgentID : Agent
*/
var carte [50][50]string = [50][50]string{
	{"X", "X", "X", "X", "X", "X", "X", "X", "W", "W", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "W", "W", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "_", "X", "X", "X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "_", "X", "X", "X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "_", "_", "X", "X", "_", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "B", "B", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "_", "_", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
	{"X", "X", "X", "X", "S", "S", "X", "X", "X", "X", "X", "X", "E", "E", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "E", "E", "X", "X", "X", "X", "X", "X", "S", "S", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
}

var playground [50][50]string = [50][50]string{
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "X", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
	{"_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_"},
}

type Simulation struct {
	env          Environment
	maxStep      int
	maxDuration  time.Duration
	step         int // Stats
	start        time.Time
	syncChans    sync.Map
	newAgentChan chan Agent // permet la création d'agent en cours de simulation (canal de com avec env)
}

// TODO:voir si agents est mis à jour lors de suppression d'agent

func (sim *Simulation) Env() *Environment {
	return &sim.env
}

func NewSimulation(agentCount int, maxStep int, maxDuration time.Duration) (simu *Simulation) {
	simu = &Simulation{}
	simu.maxStep = maxStep
	simu.maxDuration = maxDuration

	simu.newAgentChan = make(chan Agent, 200) // channel avec buffer pour gérer les sorties simultanées

	// Création de l'environement
	// Création du métro
	metro1 := *NewMetro(10*time.Second, 5*time.Second, 10, 5, NewWay(1, alg.Coord{9, 0}, alg.Coord{10, 39}, true, []alg.Coord{{8, 5}, {8, 6}, {8, 34}}, &simu.env))
	metro2 := *NewMetro(10*time.Second, 5*time.Second, 10, 5, NewWay(2, alg.Coord{11, 0}, alg.Coord{12, 39}, false, []alg.Coord{{13, 5}, {13, 6}, {13, 34}}, &simu.env))
	simu.env = *NewEnvironment([]Agent{}, carte, []Metro{metro1, metro2}, simu.newAgentChan, agentCount)
	//simu.env = *NewEnvironment([]Agent{}, playground, mapChan)

	// création des agents et des channels
	for i := 0; i < agentCount; i++ {
		// création de l'agent

		syncChan := make(chan int)
		//ag := NewAgent(id, &simu.env, syncChan, time.Duration(time.Second), 0, true, Coord{0, 8 + i%2}, Coord{0, 8 + i%2}, &UsagerLambda{}, Coord{0, 8 + i%2}, Coord{12 - 4*(i%2), 18 - 15*(i%2)})
		//ag := NewAgent(id, &simu.env, syncChan, 1000, 0, true, &UsagerLambda{},  Coord{18, 4}, Coord{0, 8}, 2, 1)

		ag := &Agent{}

		if i%2 == 0 { //Type Agent
			id := fmt.Sprintf("Agent%d", i)
			//NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, force int, politesse bool, behavior Behavior, departure, destination Coord, width, height int)
			ag = NewAgent(id, &simu.env, syncChan, 200, 0, true, &UsagerLambda{}, alg.Coord{49, 32}, alg.Coord{0, 9}, 2, 1)
		} else { // Type Controleur
			//id := fmt.Sprintf("Controleur%d", i)
			id := fmt.Sprintf("Agent%d", i)
			ag = NewAgent(id, &simu.env, syncChan, 200, 0, true, &UsagerLambda{}, alg.Coord{0, 8}, alg.Coord{8, 5}, 1, 1)
			//ag = NewAgent(id, &simu.env, syncChan, 1000, 0, true, &Controleur{}, Coord{18, 12}, Coord{18, 4}, 1, 1)
		}

		//ag := NewAgent(id, &simu.env, syncChan, 1000, 0, true, &UsagerLambda{}, Coord{19, 12}, Coord{0, 8}, 2, 1)

		// ajout de l'agent à la simulation
		simu.env.ags = append(simu.env.ags, *ag)

		simu.env.agentsChan[ag.id] = make(chan req.Request, 5)

		// ajout du channel de synchro
		simu.syncChans.Store(ag.ID(), syncChan)

	}

	return simu
}

func (simu *Simulation) Run() {
	log.Printf("Démarrage de la simulation ")
	// Démarrage du micro-service de Log
	go simu.Log()
	// Démarrage du micro-service d'affichage
	go simu.Print()

	// Démarrage des agents

	var wg sync.WaitGroup
	for _, agt := range simu.env.ags {
		wg.Add(1)
		go func(agent Agent) {
			defer wg.Done()
			agent.Start()
		}(agt)
	}

	// On sauvegarde la date du début de la simulation
	simu.start = time.Now()

	// Lancement des métros
	for _, metro := range simu.env.metros {
		wg.Add(1)
		go func(metro Metro) {
			defer wg.Done()
			metro.Start()
		}(metro)
	}

	// Lancement de l'orchestration de tous les agents
	for _, agt := range simu.env.ags {
		go func(agt Agent) {
			step := 0
			for {
				step++
				c, _ := simu.syncChans.Load(agt.ID()) // communiquer les steps aux agents
				c.(chan int) <- step                  // /!\ utilisation d'un "Type Assertion"
				time.Sleep(1 * time.Millisecond)      // "cool down"
				<-c.(chan int)
			}
		}(agt)
	}

	// Lancement du flow d'agents
	go simu.ActivateFlow(1000)

	// Activation de l'ouie des agents
	go simu.listenNewAgentChan()

	time.Sleep(simu.maxDuration)

}
func (simu *Simulation) listenNewAgentChan() {
	for {
		select {
		case newAgent := <-simu.newAgentChan:
			//simu.env.ags = append(simu.env.ags, newAgent)
			simu.syncChans.Store(newAgent.ID(), newAgent.syncChan)
			go func(agent Agent) {
				agent.Start()
			}(newAgent)
			go func(agent Agent) {
				step := 0
				for {
					step++
					c, _ := simu.syncChans.Load(agent.ID()) // communiquer les steps aux agents
					c.(chan int) <- step                    // /!\ utilisation d'un "Type Assertion"
					time.Sleep(1 * time.Millisecond)        // "cool down"
					<-c.(chan int)
				}
			}(newAgent)

			// Add the new agent to simu.agents

		}
	}
}

func (simu *Simulation) Print_v0() {
	for {
		for i := 0; i < 20; i++ {
			fmt.Println(simu.env.station[i])
		}
		//fmt.Println("=================================================================================")
		fmt.Println()
		fmt.Println()
		//time.Sleep(time.Second / 4) // 60 fps !
		time.Sleep(500 * time.Millisecond) // 1 fps !
		//fmt.Print("\033[H\033[2J") // effacement du terminal
	}
}
func (simu *Simulation) Print() [][]string {
	result := make([][]string, 50)
	for i := 0; i < 50; i++ {
		result[i] = make([]string, 50)
		for j := 0; j < 50; j++ {
			element := simu.env.station[i][j]
			if len(element) > 1 {
				result[i][j] = element[:1] // Stocker le premier caractère si la longueur est supérieure à 1
				fmt.Print(result[i][j] + " ")
			} else {
				result[i][j] = element
				fmt.Print(result[i][j] + " ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	time.Sleep(200 * time.Millisecond) // 1 fps !
	return result
}

func (simu *Simulation) Log() {
	// Not implemented
}

func (simu *Simulation) ActivateFlow(flow int) {
	for {
		time.Sleep(time.Duration(flow) * time.Millisecond)
		simu.env.AddAgent(*NewAgent("Agent"+fmt.Sprint(simu.env.agentCount), &simu.env, make(chan int), 200, 0, true, &UsagerLambda{}, simu.env.entries[rand.Intn(len(simu.env.entries))], simu.env.gates[rand.Intn(len(simu.env.gates))], 1, 1))
	}
}
