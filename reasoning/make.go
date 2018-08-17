// This file contains no code.  It is just a repository of go:generate directives.
//
//  go generate squaredance/reasoning
//
package reasoning

// Support for square dance Roles
//go:generate go build squaredance/reasoning/generate_roles
//go:generate generate_roles

// Expanding formation boilerplate code
//go:generate go build squaredance/reasoning/formation_expander
//go:generate formation_expander two_dancers.go four_dancers.go

// Compiling rules
//go:generate rule_compiler set.rules two_dancers.rules four_dancers.rules

// Noting what's emitted by each rule
//go:generate catalog_rule_type_info

