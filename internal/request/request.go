package request

type Request struct {
	demandeur chan Request //channel de l'émetteur de la demande
	decision  int
}

func NewRequest(demandeur chan Request, decision int) (req *Request) {
	return &Request{demandeur, decision}
}

func (req *Request) Demandeur() (demandeur chan Request) {
	return req.demandeur
}

func (req *Request) Decision() (decision int) {
	return req.decision
}

func (req *Request) SetDemandeur(demandeur chan Request) {
	req.demandeur = demandeur
}

func (req *Request) SetDecision(decision int) {
	req.decision = decision
}
