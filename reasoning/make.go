// This file contains no code.  It is just a repository of go:generate directives.
//
//  go generate squaredance/reasoning
//
package reasoning

// Support for square dance Roles
//go:generate go build squaredance/reasoning/generate_roles
//go:generate generate_roles
// OUTPUTS: generated_roles.go

// Expanding formation boilerplate code
//go:generate defimpl
// OUTPUTS impl_*.go
//go:generate go build squaredance/reasoning/formation_expander
//go:generate formation_expander two_dancers_rules.go four_dancers_rules.go
// OUTPUTS: feout_*.go

// Compiling rules
//go:generate rule_compiler
// INPUTS: *_rules.go
// OUTPUTS: *_rules_out.go

// Running unit tests and generating some documentation fioles:
//go:generate go test squaredance/reasoning
// OUTPUTS: all_rules.txt, formation_types.txt, formations_rete.dot

