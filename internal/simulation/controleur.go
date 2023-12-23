package simulation

//ajouter liste des agents déjà controllés

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
	alg "metrosim/internal/algorithms"
)

/*
	Je suppose que l'id du controleur est de format "Cont + un chiffre"
	Exemple : "Cont1"
	et l'id du l'agent est de format "Agent + un chiffre"
	Exemple : "Agent1"
*/

type Controleur struct {
	req *Request // requete reçue par le controleur
	faceCase string // chaine de caractère qui contient l'id de l'agent qui se trouve devant le controleur, exemple : "Agent1", "Fraudeur1", "X" ,etc.
	timer *time.Timer // timer qui permet de définir la durée de vie du controleur
	isExpired bool // true si le controleur est expiré, false sinon
}



func (c *Controleur) Percept(ag *Agent) {
	//initialiser le faceCase en fonction de la direction de l'agent
	c.faceCase = ag.getFaceCase()
	switch {
		// comportement par défaut (comportement agent Lambda)
		case ag.request != nil: //verifier si l'agent est communiqué par un autre agent (A VOIR SI IL EXISTE DEJA UN AGENT QUI COMMUNIQUE AVEC LE CONTROLEUR)
			//print("Requete recue par l'agent lambda : ", ag.request.decision, "\n")
			c.req = ag.request
		default:
			ag.stuck = ag.isStuck()
			if ag.stuck {
				return
	
			}
	}
}

func (c *Controleur) Deliberate(ag *Agent) {
	// Verifier si la case devant lui contient un agent ou un fraudeur
	// Créer l'expression régulière
	regexAgent := `^Agent\d+$` // \d+ correspond à un ou plusieurs chiffres
	regexFraudeur := `^Fraudeur\d+$`

	// Vérifier si la valeur de faceCase ne correspond pas au motif
	matchedAgt, err1 := regexp.MatchString(regexAgent, c.faceCase)
	matchedFraud, err2 := regexp.MatchString(regexFraudeur, c.faceCase)
	//fmt.Println("faceCase : ", c.faceCase)
	//fmt.Println("matchedAgt : ", matchedAgt)

	if err1 != nil || err2 != nil {
		fmt.Println("Erreur lors de l'analyse de la regex :", err1, err2)
		return
	} else {
		if matchedAgt && ag.env.controlledAgents[AgentID(c.faceCase)] == false { // si l'agent devant le controleur est un agent et qu'il n'a pas encore été controlé
			//fmt.Println("L'agent ", c.face, " a été détecté par le controleur")
			ag.decision = Stop // arreter l'agent devant lui
		} else if matchedFraud && !ag.env.controlledAgents[AgentID(c.faceCase)] {
			ag.decision = Expel // virer l'agent devant lui
			//sinon comportement par défaut (comportement de l'usager lambda)
			}else if ag.position == ag.destination && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si l'agent est arrivé à sa destination et qu'il est sur une sortie
				//fmt.Println(ag.id, "disappear")
				ag.decision = Disappear
				} else if ag.stuck{ // si l'agent est bloqué
					ag.decision = Wait
					}else {
					ag.decision = Move
					}
	}
}

func (c *Controleur) Act(ag *Agent) {
	switch ag.decision {
	case Move:
		if !c.isExpired {
			//fmt.Printf("[Controleur, Act, non expiré] Le controleur %s est en mouvement \n", ag.id)
			ag.destination = c.randomDestination(ag)
			//fmt.Printf("[Controleur, Act] destination s = %s : %d \n",ag.id,ag.destination)
		}else {
			//fmt.Printf("[Controleur, Act] Le controleur %s est expiré \n",ag.id)
			ag.destination = ag.findNearestExit()
			//fmt.Printf("[Controleur, Act, Expire] destination de %s = %d \n",ag.id,ag.destination)
		}
		ag.MoveAgent()

	case Wait:
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)

	case Disappear:
		RemoveAgent(&ag.env.station, ag)

	default : //Expel ou Wait
		agt_face_id := AgentID(c.faceCase) //id de l'agent qui se trouve devant le controleur
		//fmt.Print("L'agent ", agt_face_id, " a été expulsé\n")
		ag.env.agentsChan[agt_face_id] <- *NewRequest(ag.env.agentsChan[ag.id], ag.decision) // envoie la decision du controleur à l'agent qui se trouve devant lui
	}
}

func (c *Controleur) randomDestination(ag *Agent) alg.Coord {
	rand.Seed(time.Now().UnixNano()) // le générateur de nombres aléatoires
	randomRow := rand.Intn(20) // Génère un entier aléatoire entre 0 et 19
	randomCol := rand.Intn(20) // Génère un entier aléatoire entre 0 et 19
	for ag.env.station[randomRow][randomCol] != "_" {
		randomRow = rand.Intn(20) // Génère un entier aléatoire entre 0 et 19
		randomCol = rand.Intn(20) // Génère un entier aléatoire entre 0 et 19
	}
	return alg.Coord{randomRow, randomCol}
}

func (c *Controleur) startTimer() {
	rand.Seed(time.Now().UnixNano()) // le générateur de nombres aléatoires
	//randomSeconds := rand.Intn(9) + 2 // Génère un entier aléatoire entre 2 et 10
	randomSeconds := 2
	lifetime := time.Duration(randomSeconds) * time.Second
    c.timer = time.NewTimer(lifetime)
	//fmt.Println("[Controleur , startTimer] Le controleur est créé avec une durée de vie de ", lifetime)
    go func() {
        <-c.timer.C // attend que le timer expire
        c.isExpired = true
    }()
}

