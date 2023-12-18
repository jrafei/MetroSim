package simulation

//ajouter liste des agents déjà controllés

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
	//"time"
)

/*
	Je suppose que l'id du controleur est de format "Cont + un chiffre"
	Exemple : "Cont1"
	et l'id du l'agent est de format "Agent + un chiffre"
	Exemple : "Agent1"
*/

type Controleur struct {
	faceCase string // chaine de caractère qui contient l'id de l'agent qui se trouve devant le controleur, exemple : "Agent1", "Fraudeur1", "X" ,etc.
}

func (c *Controleur) Percept(ag *Agent) {
	env := ag.env

	if ag.direction == 0 { // vers le haut
		c.faceCase = env.station[ag.position[0]-1][ag.position[1]]
	} else if ag.direction == 1 { // vers la droite
		c.faceCase = env.station[ag.position[0]][ag.position[1]+1]
	} else if ag.direction == 2 { // vers le bas
		c.faceCase = env.station[ag.position[0]+1][ag.position[1]]
	} else { // vers la gauche
		c.faceCase = env.station[ag.position[0]][ag.position[1]-1]
	}

}

func (c *Controleur) Deliberate(ag *Agent) {
	// Verifier si la  devant lui contient un agent ou un fraudeur
	// Créer l'expression régulière
	regexAgent := `^Agent\d+$` // \d+ correspond à un ou plusieurs chiffres
	regexFraudeur := `^Fraudeur\d+$`

	// Vérifier si la valeur de faceCase ne correspond pas au motif
	matchedAgt, err1 := regexp.MatchString(regexAgent, c.faceCase)
	matchedFraud, err2 := regexp.MatchString(regexFraudeur, c.faceCase)
	//fmt.Println("faceCase : ", c.faceCase)
	//fmt.Println("matchedAgt : ", matchedAgt)

	if err1 != nil || err2 != nil {
		//fmt.Println("Erreur lors de l'analyse de la regex :", err1, err2)
		return
	} else {
		if matchedAgt && ag.env.controlledAgents[AgentID(c.faceCase)] == false { // si l'agent devant le controleur est un agent et qu'il n'a pas encore été controlé
			//fmt.Println("L'agent ", c.face, " a été détecté par le controleur")
			ag.decision = Stop // arreter l'agent devant lui
		} else if matchedFraud && !ag.env.controlledAgents[AgentID(c.faceCase)] {
			ag.decision = Expel // virer l'agent devant lui
		} else {
			// Comportement de l'usager lambda (comportement par defaut)
			if ag.stuck {
				ag.decision = Wait // attendre
			} else {
				ag.decision = Move // avancer
			}
		}
	}
}

func (c *Controleur) Act(ag *Agent) {
	if ag.decision == Move {
		ag.MoveAgent()
	} else if ag.decision == Wait {
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	} else {
		agt_face_id := AgentID(c.faceCase) //id de l'agent qui se trouve devant le controleur
		fmt.Print("L'agent ", agt_face_id, " a été expulsé\n")
		ag.env.agentsChan[agt_face_id] <- *NewRequest(ag.env.agentsChan[ag.id], ag.decision) // envoie la decision du controleur à l'agent qui se trouve devant lui
	}
}
