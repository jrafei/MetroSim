package simulation

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
	alg "metrosim/internal/algorithms"
)


type MobiliteReduite struct {
	req *Request
	once sync.Once
}


/*
* Fonction qui permet de définir la destination d'un agent à mobilité réduite
*/
func (mr *MobiliteReduite)setUpDestination(ag *Agent){
	choix_voie := rand.Intn(2) // choix de la voie de départ aléatoire
	//fmt.Println("[MobiliteReduite, setUpPath] choix_voie = ",choix_voie)
	dest_porte := (ag.findNearestGates(ag.env.metros[choix_voie].way.gates))[0].Gate
	fmt.Println("[MobiliteReduite, setUpPath] dest_porte = ",dest_porte)
	ag.destination = dest_porte
	mr.setUpPath(ag)
}
/*
	Calcule le chemin de l'agent à mobilité réduite vers la porte la plus proche
*/
func (mr * MobiliteReduite) setUpPath(ag *Agent) {
	start, end := ag.generatePathExtremities()
	// Recherche d'un chemin si inexistant
	path := alg.FindPath(ag.env.station, start, end, *alg.NewNode(-1, -1, 0, 0, 0, 0), false, 2*time.Second)
	ag.path = path
	fmt.Println("[MobiliteReduite, setUpPath] path = ",ag.path)
	ag.direction = calculDirection(ag.position, Coord{ag.path[0].Row(), ag.path[0].Col()})
}

func (mr *MobiliteReduite) Percept(ag *Agent) {
	mr.once.Do(func(){mr.setUpDestination(ag)}) // la fonction setUp est executé à la premiere appel de la fonction Percept()
	switch {
	case ag.request != nil: //verifier si l'agent est communiqué par un autre agent, par exemple un controleur lui a demandé de s'arreter
		//print("Requete recue par l'agent lambda : ", ag.request.decision, "\n")
		mr.req = ag.request
	default:
		ag.stuck = ag.isStuck()
		if ag.stuck {
			return
		}
	}
}

func (mr *MobiliteReduite) Deliberate(ag *Agent) {
	//fmt.Println("[AgentLambda Deliberate] decision :", ul.req.decision)
	if (mr.req != nil ) {
		if mr.req.decision == Stop{
			ag.decision = Wait
			mr.req = nil //demande traitée
		} else if mr.req.decision == Expel { // sinon alors la requete est de type "Viré" cette condition est inutile car l'usager lambda ne peut pas etre expulsé , elle est nécessaire pour les agents fraudeurs
				//fmt.Println("[AgentLambda, Deliberate] Expel")
				ag.decision = Expel
				mr.req = nil //demande traitée
			}else if mr.req.decision == Disappear {
				ag.decision = Disappear
			}
	}else if ag.position == ag.destination && (ag.isOn[ag.position] == "W" || ag.isOn[ag.position] == "S") { // si l'agent est arrivé à sa destination et qu'il est sur une sortie
			//fmt.Println(ag.id, "disappear")
			ag.decision = Disappear
		} else if ag.stuck{ // si l'agent est bloqué
			ag.decision = Wait
			}else {
			ag.decision = Move
			}	
}

func (mr *MobiliteReduite) Act(ag *Agent) {
	//fmt.Println("[AgentLambda Act] decision :",ag.decision)
	switch ag.decision {
	case Move:
		mr.MoveMR(ag)
	case Wait:
		n := rand.Intn(2) // temps d'attente aléatoire
		time.Sleep(time.Duration(n) * time.Second)
	case Disappear:
		RemoveAgent(&ag.env.station, ag)
	case Expel : 
		//fmt.Println("[AgentLambda, Act] Expel")
		ag.destination = ag.findNearestExit()
		//fmt.Println("[AgentLambda, Act] destination = ",ag.destination)
		ag.env.controlledAgents[ag.id] = true
		ag.path = make([]alg.Node,0)
		ag.MoveAgent()
	}
}

/*
* Fonction qui permet de déplacer un agent à mobilité réduite
*/
func (mr *MobiliteReduite) MoveMR(ag *Agent) {
	// ================== Déplacement si aucun problème =======================
	if(len(ag.path) <= 0){
		mr.setUpPath(ag)
	}
	safe, or := IsMovementSafe(ag.path, ag, ag.env)
	if safe {
		if len(ag.isOn) > 0 {
			RemoveAgent(&ag.env.station, ag)
		}
		rotateAgent(ag, or) // mise à jour de l'orientation
		//ag.env.station[ag.coordBasOccupation[0]][ag.coordBasOccupation[1]] = ag.isOn
		
		//MODIFICATION DE DIRECTION 
		ag.direction = calculDirection(ag.position, Coord{ag.path[0].Row(), ag.path[0].Col()})
		//fmt.Println("[MoveAgent]Direction : ", ag.direction)
		ag.position[0] = ag.path[0].Row()
		ag.position[1] = ag.path[0].Col()
		if len(ag.path) > 1 {
			ag.path = ag.path[1:]
		} else {
			ag.path = nil
		}
		saveCells(&ag.env.station, ag.isOn, ag.position, ag.width, ag.height, ag.orientation)
		writeAgent(&ag.env.station, ag)
		// ============ Prise en compte de la vitesse de déplacement ======================
		time.Sleep(ag.vitesse * time.Millisecond)
	}
}
