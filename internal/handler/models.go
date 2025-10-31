package handler

type AddDomainReq struct {
	Domain string `json:"domain"`
	Target string `json:"target"`
}
