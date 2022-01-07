package schema

import (
	"aws-blended-instances-advisor/utils"
	"errors"
	"fmt"
)

type Service struct {
	Name         string  `json:"name"`
	MinMemory    float64 `json:"minMemory"`
	MaxVcpu      int     `json:"maxVcpu"`
	MinInstances int     `json:"minInstances"`
	MaxInstances int     `json:"maxInstances"`
}

// Validate checks that a Service is well-formed
// and is true to the API specification.
func (s *Service) Validate() error {
	if s.MinMemory < 0 {
		return errors.New("minMemory is not positive")
	}
	if s.MaxVcpu <= 0 {
		return errors.New("maxVcpu is not postive")
	}
	if s.MinInstances < 0 {
		return errors.New("minInstances is not positive")
	}
	if s.MaxInstances <= 0 {
		return errors.New("totalInstances is not positive")
	}
	if s.MinInstances > s.MaxInstances {
		return errors.New("minInstances is greater than totalInstances")
	}
	return nil
}

// ValidateServices validates multiple services, ensuring that
// they are well-formed and true to the API specification.
func ValidateServices(services []Service) error {
	if !namesAreUnique(services) {
		return errors.New("service names are not unique")
	}
	for _, s := range services {
		err := s.Validate()
		if err != nil {
			return utils.PrependToError(
				err,
				fmt.Sprintf("service %s invalid", s.Name),
			)
		}
	}
	return nil
}

func namesAreUnique(services []Service) bool {
	namesSet := make(map[string]bool)
	for _, s := range services {
		if _, exists := namesSet[s.Name]; exists {
			return false
		}
		namesSet[s.Name] = true
	}
	return true
}
