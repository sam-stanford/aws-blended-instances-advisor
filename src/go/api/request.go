package api

type Request struct {
	Services []Service `json:"services"`
	Advisor  Advisor   `json:"advisor"`
}
