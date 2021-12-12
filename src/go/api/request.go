package api

// TODO: Response type in new file or this one renamed schema.go
// TODO: Move all schema to schema sub package & move service to this package and remove all functions for schema and put them in this package, such that schema ONLY contains schema def

type Request struct {
	Regions  []Region  `json:"regions"`
	Services []Service `json:"services"`
	Advisor  Advisor   `json:"advisor"`
}
