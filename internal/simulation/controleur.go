package simulation

/*
	Le controleur se déplace aléatoirement dans la station pendant un certain temps
	et controle les agents qui se trouvent devant lui si ils ne sont pas déjà controllés par un autre controleur,
	si l'agent est un fraudeur alors il est expulsé, sinon il est arreté pendant un certain temps
*/

import (
	"fmt"
	"math/rand"
	alg "metrosim/internal/algorithms"
	req "metrosim/internal/request"
	"regexp"
	"time"
)

type Controleur struct {
	req       *req.Request // requete reçue par le controleur
	faceCase  string       // chaine de caractère qui contient l'id de l'agent qui se trouve devant le controleur, exemple : "Agent1", "Fraudeur1", "X" ,etc.
	timer     *time.Timer  // timer qui permet de définir la durée de vie du controleur
	isExpired bool         // true si le controleur est expiré, false sinon
}

func (c *Controleur) Percept(ag *Agent) {
	//initialiser le faceCase en fonction de la direction de l'agent
	c.faceCase = ag.getFaceCase()
	switch {
	// comportement par défaut (comportement agent Lambda)
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent
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
	regexCont := `^Cont\d+$` // \d+ correspond à un ou plusieurs chiffres
	regexFraudeur := `^Fraudeur\d+$`

	// Vérifier si la valeur de faceCase ne correspond pas au motif
	matchedCont, err1 := regexp.MatchString(regexCont, c.faceCase)
	matchedFraud, err := regexp.MatchString(regexFraudeur, c.faceCase)
	//
	//.0fmt.Println("faceCase : ", c.faceCase)
	//fmt.Println("matchedAgt : ", matchedAgt)

	if err != nil {
		fmt.Println("Erreur lors de l'analyse de la regex :", err)
		return
	}
	if err1 != nil {
		fmt.Println("Erreur lors de l'analyse de la regex :", err1)
		return
	} else {
		if matchedCont {
			// si l'agent devant le controleur est un autre controleur alors il faut attendre qu'il se déplace
			ag.decision = Move
			return
		}
		if matchedFraud && !ag.env.controlledAgents[AgentID(c.faceCase)] {
			ag.decision = Expel // virer l'agent devant lui
			return
		} else if ag.position == ag.destination && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si le controleur est arrivé à sa destination et qu'il est sur une sortie
			//fmt.Println(ag.id, "disappear")
			ag.decision = Disappear
			return
		} else if ag.stuck { // si le controleur est bloqué
			ag.decision = Wait
			return
		} else {
			ag.decision = Move
			return
		}
	}
}

func (c *Controleur) Act(ag *Agent) {
	switch ag.decision {
	case Move:
		if !c.isExpired {
			//fmt.Printf("[Controleur, Act, non expiré] Le controleur %s est en mouvement \n", ag.id)
			c.SetUpDestination(ag)
			//fmt.Printf("[Controleur, Act] destination s = %s : %d \n",ag.id,ag.destination)
		} else {
			//fmt.Printf("[Controleur, Act] Le controleur %s est expiré \n",ag.id)
			ag.destination = ag.findNearestExit()
			//fmt.Printf("[Controleur, Act, Expire] destination de %s = %d \n",ag.id,ag.destination)
		}
		ag.MoveAgent()

	case Wait:
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)

	case Disappear:
		ag.env.RemoveAgent(ag)

	case Expel: //Expel
		agt_face_id := AgentID(c.faceCase) //id de l'agent qui se trouve devant le controleur
		//TODELETEfmt.Print("L'agent ", agt_face_id, " a été expulsé \n")
		ag.env.controlledAgents[agt_face_id] = true                                              // l'agent qui se trouve devant le controleur est controlé
		ag.env.agentsChan[agt_face_id] <- *req.NewRequest(ag.env.agentsChan[ag.id], ag.decision) // envoie la decision du controleur à l'agent qui se trouve devant lui
		//time.Sleep(500 * time.Millisecond)
	}
	ag.request = nil
}

func (c *Controleur) SetUpDestination(ag *Agent) {
	rand.Seed(time.Now().UnixNano())               // le générateur de nombres aléatoires
	randomRow := rand.Intn(len(ag.env.station[0])) // Génère un entier aléatoire entre 0 et 49
	randomCol := rand.Intn(len(ag.env.station[1])) // Génère un entier aléatoire entre 0 et 49
	for ag.env.station[randomRow][randomCol] != "_" {
		randomRow = rand.Intn(len(ag.env.station[0])) // Génère un entier aléatoire entre 0 et 49
		randomCol = rand.Intn(len(ag.env.station[1])) // Génère un entier aléatoire entre 0 et 49
	}
	ag.destination = alg.Coord{randomRow, randomCol}
}

func (c *Controleur) startTimer() {
	//rand.Seed(time.Now().UnixNano()) // le générateur de nombres aléatoires
	//randomSeconds := rand.Intn(9) + 2 // Génère un entier aléatoire entre 2 et 10
	randomSeconds := 500
	lifetime := time.Duration(randomSeconds) * time.Second
	c.timer = time.NewTimer(lifetime)
	//fmt.Println("[Controleur , startTimer] Le controleur est créé avec une durée de vie de ", lifetime)
	go func() {
		<-c.timer.C // attend que le timer expire
		c.isExpired = true
	}()
}
