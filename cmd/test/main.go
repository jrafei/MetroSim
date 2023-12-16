package main

import (
	"fmt"
	"regexp"
)

func main() {
	faceCase := "Agent1000" // Remplacez ceci par votre variable

	// Créer l'expression régulière
	regexPattern := `^Agent\d+$` // \d+ correspond à un ou plusieurs chiffres
	matched, err := regexp.MatchString(regexPattern, faceCase)

	if err != nil {
		fmt.Println("Erreur lors de l'analyse de la regex :", err)
		return
	}

	// Vérifiez si la chaîne ne correspond pas au motif
	if !matched {
		fmt.Println("La chaîne ne correspond pas au motif 'Agentx'")
	} else {
		fmt.Println("La chaîne correspond au motif 'Agentx'")
	}
}
