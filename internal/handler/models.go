package handler

type AddDomainReq struct {
	Domain string `json:"domain"`
	Target string `json:"target"`
}

type RemoveDomainReq struct {
	Domain string `json:"domain"`
}
