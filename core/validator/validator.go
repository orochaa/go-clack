package validator

import "fmt"

type Validator struct {
	Name string
}

// NewValidator creates and returns a new instance of Validator with the specified name.
//
// Parameters:
//   - name: The name of the prompt being validated.
//
// Returns:
//   - *Validator: A pointer to the newly created Validator instance.
func NewValidator(name string) *Validator {
	return &Validator{
		Name: name,
	}
}

// Panic triggers a panic with a formatted message that includes the validator's name and the provided cause.
//
// Parameters:
//   - cause: The reason for the panic.
func (v *Validator) Panic(cause string) {
	panic(fmt.Sprintf("%s: %s", v.Name, cause))
}

// PanicMissingParam triggers a panic indicating that a specific parameter is missing.
// The panic message is formatted to include the validator's name and the missing parameter.
//
// Parameters:
//   - param: The name of the missing parameter.
func (v *Validator) PanicMissingParam(param string) {
	v.Panic(fmt.Sprintf("missing %sParams.%s", v.Name, param))
}

// ValidateRender checks if a render function is provided. If the render function is nil,
// it triggers a panic indicating that the "Render" parameter is missing.
//
// Parameters:
//   - render: The render function to validate.
func (v *Validator) ValidateRender(render any) {
	if render == nil {
		v.PanicMissingParam("Render")
	}
}

// ValidateOptions checks if at least one option is provided. If no options are provided,
// it triggers a panic indicating that the "Options" parameter is missing.
//
// Parameters:
//   - length: The number of options provided.
func (v *Validator) ValidateOptions(length int) {
	if length == 0 {
		v.PanicMissingParam("Options")
	}
}
