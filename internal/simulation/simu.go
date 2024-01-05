package simulation

import (
	"fmt"
	"log"
	"math/rand"
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

type Simulation struct {
	env          Environment
	maxStep      int
	maxDuration  time.Duration
	step         int // Stats
	start        time.Time
	syncChans    sync.Map
	newAgentChan chan Agent // permet la création d'agent en cours de simulation (canal de com avec env)
	active       bool       // simulation active (en cours)
	controleurs  bool       // présence de controleurs
	fraudeurs    bool       // présence de fraudeurs
	impolis      bool       // présence d'impolis
	mob_reduite  bool       // présence de personnes à mobilité réduite
	patients     bool       // présence d'usagers patients
	flow         int        //débit de création d'agents/seconde
}

func (sim *Simulation) Env() *Environment {
	return &sim.env
}

func NewSimulation(conf req.Cfg) (simu *Simulation) {
	simu = &Simulation{}
	simu.maxStep = -1
	simu.maxDuration = time.Duration(conf.MaxDuration) * time.Second
	simu.flow = conf.Flow

	simu.controleurs = conf.Controleurs
	simu.patients = conf.Patients
	simu.fraudeurs = conf.Fraudeurs
	simu.impolis = conf.Impolis
	simu.mob_reduite = conf.Mob_reduite

	simu.newAgentChan = make(chan Agent, 200) // channel avec buffer pour gérer les sorties simultanées

	// Création du métro
	metros := []Metro{}
	for i, _ := range conf.LeftTopCorners {
		metro := *NewMetro(time.Duration(conf.Frequency[i])*time.Second, time.Duration(conf.StopTime[i])*time.Second, conf.Capacity[i], rand.Intn(conf.Capacity[i]), NewWay(WayID(i), conf.LeftTopCorners[i], conf.RightBottomCorners[i], conf.GoToLeft[i], conf.Gates[i], &simu.env))
		metros = append(metros, metro)
	}

	// Création de l'environnement
	simu.env = *NewEnvironment([]Agent{}, conf.Station, metros, simu.newAgentChan, 0, simu)

	// Simulation pas encore démarrée
	simu.active = false

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
	go simu.ActivateFlow()

	// Activation de l'ouie des agents
	go simu.listenNewAgentChan()

	simu.active = true

	time.Sleep(simu.maxDuration)

}
func (simu *Simulation) listenNewAgentChan() {
	// Ecoute du channel de création d'agents
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

func (simu *Simulation) Print() [][]string {
	// Affichage de la station
	result := make([][]string, 50)
	for i := 0; i < 50; i++ {
		result[i] = make([]string, 50)
		for j := 0; j < 50; j++ {
			element := simu.env.station[i][j]
			if len(element) > 1 {
				result[i][j] = element[:1] // Stocker le premier caractère si la longueur est supérieure à 1
				//fmt.Print(result[i][j] + " ")
			} else {
				result[i][j] = element
				//fmt.Print(result[i][j] + " ")
			}
		}
		//TODELETEfmt.Println()
	}
	//TODELETEfmt.Println()
	time.Sleep(200 * time.Millisecond) // 1 fps !
	return result
}

func (simu *Simulation) Log() {
	// Not implemented
}

func (simu *Simulation) ActivateFlow() {
	rand.Seed(time.Now().UnixNano())
	probability := rand.Float64()
	for {
		ag := &Agent{}
		ag = nil
		//TODELETEfmt.Println(probability)
		switch {
		case probability < 0.1: // 10% de probabilité
			ag = NewAgent("Normal"+fmt.Sprint(simu.env.agentCount), &simu.env, make(chan int), 200, true, &UsagerNormal{}, simu.env.entries[rand.Intn(len(simu.env.entries))], simu.env.gates[rand.Intn(len(simu.env.gates))], 1, 1)
		case probability < 0.5: // 40% de probabilité (cumulative 50%)
			if simu.patients {
				ag = NewAgent("Patient"+fmt.Sprint(simu.env.agentCount), &simu.env, make(chan int), 200, true, &UsagerLambda{}, simu.env.entries[rand.Intn(len(simu.env.entries))], simu.env.gates[rand.Intn(len(simu.env.gates))], 1, 1)
			}

		case probability < 0.6: // 10% de probabilité (cumulative 60%)
			if simu.controleurs {
				ag = NewAgent("Cont"+fmt.Sprint(simu.env.agentCount), &simu.env, make(chan int), 200, true, &Controleur{}, simu.env.entries[rand.Intn(len(simu.env.entries))], simu.env.gates[rand.Intn(len(simu.env.gates))], 1, 1)
			}

		case probability < 0.8: // 20% de probabilité (cumulative 80%)
			if simu.fraudeurs {
				ag = NewAgent("Fraudeur"+fmt.Sprint(simu.env.agentCount), &simu.env, make(chan int), 200, true, &UsagerLambda{}, simu.env.entries[rand.Intn(len(simu.env.entries))], simu.env.gates[rand.Intn(len(simu.env.gates))], 1, 1)
			}

		case probability < 0.9: // 10% de probabilité (cumulative 90%)
			if simu.impolis {
				// ag = *NewAgent("Impoli"+fmt.Sprint(simu.env.agentCount), &simu.env, make(chan int), 200, false, &UsagerLambda{}, simu.env.entries[rand.Intn(len(simu.env.entries))], simu.env.gates[rand.Intn(len(simu.env.gates))], 1, 1)
			}

		case probability < 1.0: // 10% de probabilité (cumulative 100%)
			if simu.mob_reduite {
				ag = NewAgent("H"+fmt.Sprint(simu.env.agentCount), &simu.env, make(chan int), 200, true, &MobiliteReduite{}, simu.env.entries[rand.Intn(len(simu.env.entries))], simu.env.gates[rand.Intn(len(simu.env.gates))], 1, 1)
			}
		}
		if ag != nil {
			ag.behavior.SetUpDestination(ag)
			simu.env.AddAgent(*ag)
		}

		// NewAgent(id string, env *Environment, syncChan chan int, vitesse time.Duration, politesse bool, behavior Behavior, departure, destination alg.Coord, width, height int) *Agent {

		time.Sleep(time.Duration(simu.flow) * time.Millisecond)
		//log.Println(simu.env.ags[len(simu.env.ags)-1].path)
		probability = rand.Float64()
	}
}

func (simu *Simulation) IsRunning() bool {
	// Détermine si la simulation est en cours
	return simu.active
}
