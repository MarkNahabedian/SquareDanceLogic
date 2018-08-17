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

If the rule compiler is not already built you'll need to
<pre>
go install goshua/rete/rule_compiler
</pre>

After the rules are compiled, catalog_rule_type_info should be run to
write a rule_emits.go file which describes what types of objects are
emitted by each rule.
<pre>
go install goshua/rete/rule_compiler/catalog_rule_type_info
</pre>


Much of the code that defines each formation is automatically
generated from that formation's interface definition.

<pre>
go build squaredance/reasoning/formation_expander
formation_expander.exe two_dancers.go
</pre>

<b>All of these steps</b> (except for that from goshua/rete)
are automated in make.go, so all you should need to do is

<pre>
go generate make.go
</pre>
