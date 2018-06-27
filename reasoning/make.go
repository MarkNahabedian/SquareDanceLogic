// This file contains no code.  It is just a repository og go:generate directives.
package reasoning

// Support for square dance Roles
//go:generate go build squaredance/reasoning/generate_roles
//go:generate generate_roles

// Compiling rules
//go:generate rule_compiler set.rules two_dancers.rules 

// Expanding formation boilerplate code
//go:generate go build squaredance/reasoning/formation_expander
//go:generate formation_expander two_dancers.go
