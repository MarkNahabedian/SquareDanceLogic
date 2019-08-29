// This file contains no code.  It is just a repository of go:generate directives.
//
//  go generate squaredance/reasoning
//
package reasoning

// Support for square dance Roles
//go:generate go build squaredance/reasoning/generate_roles
//go:generate generate_roles

// Expanding formation boilerplate code
//go:generate defimpl
//go:generate go build squaredance/reasoning/formation_expander
//go:generate formation_expander two_dancers_rules.go four_dancers_rules.go

// Compiling rules
//go:generate rule_compiler

