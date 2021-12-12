package api

type Service struct {
	Name           string  `json:"name"`
	MinMemory      float64 `json:"minMemory"`
	MaxVcpu        int     `json:"maxVcpu"`
	MinInstances   int     `json:"minInstances"`
	TotalInstances int     `json:"totalInstances"`
}

// TODO: Doc, test & do for other API types
func (svc *Service) Validate() error {
	// TODO
	return nil
}
