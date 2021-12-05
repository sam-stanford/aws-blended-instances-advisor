package api

// TODO: Doc in README

type Advice map[string]RegionAdvice

type RegionAdvice struct {
	Score       float64     `json:"score"`
	Instances   []Instance  `json:"instances"`
	Assignments Assignments `json:"assignments"`
}

type Assignments struct {
	ServicesToInstances map[string][]string `json:"servicesToInstances"`
	InstancesToServices map[string][]string `json:"instancesToServices"`
}
