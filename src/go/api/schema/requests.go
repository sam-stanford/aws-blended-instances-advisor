package schema

type AdviseRequest struct {
	Services []Service `json:"services"`
	Advisor  Advisor   `json:"advisor"`
	Options  Options   `json:"options"`
}

func (r *AdviseRequest) Validate() error {
	err := ValidateServices(r.Services)
	if err != nil {
		return err
	}
	err = r.Advisor.Validate()
	if err != nil {
		return err
	}
	err = r.Options.Validate()
	return err
}
