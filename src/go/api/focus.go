package api

type ServiceFocus string

const (
	Availability ServiceFocus = "availability"
	Balanced     ServiceFocus = "balanced"
	Cost         ServiceFocus = "cost"
	Performance  ServiceFocus = "performance"
)
