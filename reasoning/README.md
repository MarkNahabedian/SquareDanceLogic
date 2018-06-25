The reasoning package is an attempt to reason about square dance calls
and formations.  It uses my "goshua" expert system tools, particularly
the rete implementation and rule compiler.

This package depends on several files of automatically generated code:

<b>generated_roles.go</b> is the output of the generate_roles command:

<pre>
go build squaredance/reasoning/generate_roles
./generate_roles
<pre>

Each of the files with a ".rules" extension must be processed by the
rule compiler at goshua/rete/rule_compiler.  It will output one ".go"
file for each ".rules" file.
