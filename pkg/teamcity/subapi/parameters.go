package subapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

type paramType = string

const (
	configParamType = "configuration"
	systemParamType = "system"
	envVarParamType = "env"
)

// ParameterTypes represent the possible parameter types
var (
	ParameterTypes = struct {
		Configuration       paramType
		System              paramType
		EnvironmentVariable paramType
	}{
		Configuration:       configParamType,
		System:              systemParamType,
		EnvironmentVariable: envVarParamType,
	}

	paramPrefixByType = map[string]string{
		ParameterTypes.Configuration:       "",
		ParameterTypes.System:              "system.",
		ParameterTypes.EnvironmentVariable: "env.",
	}
)

// Parameters is a collection of Parameters.
type Parameters struct {
	Count int32        `json:"count,omitempty" xml:"count"`
	Href  string       `json:"href,omitempty" xml:"href"`
	Items []*Parameter `json:"property,omitempty"`
}

// Parameter represents a project or build configuration parameter
type Parameter struct {
	Name      string `json:"name,omitempty" xml:"name"`
	Value     string `json:"value" xml:"value"`
	Type      string `json:"-"`
	Inherited bool   `json:"inherited,omitempty" xml:"inherited"`
}

// MarshalJSON implements JSON serialization for Parameter
func (p *Parameter) MarshalJSON() ([]byte, error) {
	out := p.Property()

	return json.Marshal(out)
}

// UnmarshalJSON implements JSON deserialization for Parameter
func (p *Parameter) UnmarshalJSON(data []byte) error {
	var property Property
	if err := json.Unmarshal(data, &property); err != nil {
		return err
	}

	var name, paramType string

	if strings.HasPrefix(property.Name, "system.") {
		name = strings.TrimPrefix(property.Name, "system.")
		paramType = ParameterTypes.System
	} else if strings.HasPrefix(property.Name, "env.") {
		name = strings.TrimPrefix(property.Name, "env.")
		paramType = ParameterTypes.EnvironmentVariable
	} else {
		name = property.Name
		paramType = ParameterTypes.Configuration
	}

	p.Name = name

	if property.Inherited != nil {
		p.Inherited = *property.Inherited
	}

	p.Value = property.Value
	p.Type = paramType
	return nil
}

// Property converts a Parameter instance to a Property
func (p *Parameter) Property() *Property {
	out := &Property{
		Name:  fmt.Sprintf("%s%s", paramPrefixByType[p.Type], p.Name),
		Value: p.Value,
	}
	return out
}
