/*
formation_expander looks for Formation interface definitions in a
go source file and writes a new source file with additional
definitions that implement those Formation interfaces.

For example, from the definition

type TwoFacedLine interface {
	Formation
	TwoFacedLine()
	Couple1() Couple
	Couple2() Couple
	MiniWave() MiniWave   // redundant
	Handedness() ...      // no-slot
}

we would generate the definition of the implementing struct

type implTwoFacedLine struct {
	couple1 Couple
	couple2 Couple
	miniwave Miniwave	// redundant
}

func (f *implTwoFacedLine) TwoFacedLine() {}

and the field accessor methods

func (f *implTwoFacedLine) Couple1() Couple { return f.couple1 }
func (f *implTwoFacedLine) Couple2() Couple { return f.couple2 }
func (f *implTwoFacedLine) MiniWave() Couple { return f.miniwave }

We also define the methods of the Formation interface itself:

func (f *implTwoFacedLine) NumberOfDancers() int { ... }
func (f *implTwoFacedLine) Dancers() dancer.Dancers { ... }
func (f *implTwoFacedLine) HasDancer(d dancer.Dancer) bool { ... }

*/

package main

import "bytes"
import "flag"
import "fmt"
import "os"
import "path"
import "strings"
import "go/ast"
import "go/parser"
import "go/token"
import "go/format"
import "goshua/go_tools"


// Any method of a Formation interface definition that has a comment
// containing this keyword is not used to compute the Dancers of a formation.
const Redundant string = "redundant"

// This keyword in a Formation interface method comment indicates that
// the struct that implements the interface should not define a slot
// corresponding to this method.
const NoSlot string = "no-slot"

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s formation_defining_go_source_file...\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "For each interface definition that defines a Formation type, generates the implementation of that formation.\n")
	flag.PrintDefaults()
}


func main() {
	flag.Parse()
	input_fileset := token.NewFileSet()
	for _, f := range flag.Args() {
		processFile(input_fileset, f)
	}
}


func output_path(input_path string) string {
	cleaned := path.Clean(input_path)
	return path.Join(path.Dir(cleaned), "impl_" + path.Base(cleaned))
}


const output_file_template string = `
// This file was automatically generated by %s from %s.
package foo
import "squaredance/dancer"
` // end


func processFile(input_fileset *token.FileSet, filepath string) {
	fmt.Printf("Processing file %s\n", filepath)
	astFile, err := parser.ParseFile(input_fileset, filepath, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	any := false
	package_name := astFile.Name.Name
	output_fileset := token.NewFileSet()
	newAstFile, err := parser.ParseFile(output_fileset, "",
		fmt.Sprintf(output_file_template,
			os.Args[0], filepath),
		parser.ParseComments)
	newAstFile.Name = ast.NewIdent(package_name)    // package name
	for _, decl := range astFile.Decls {
		ok, fdef := IsFormationDefinition(input_fileset, decl)
		if !ok {
			continue
		}
		any = true
		newAstFile.Decls = append(newAstFile.Decls,
			fdef.generate(output_fileset)...)
	}
	if !any {
		fmt.Fprintf(os.Stderr, "  No new definitions.\n")
		return
	}
	output_file := output_path(filepath)
	out, err := os.Create(output_file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Can't create %s: %s", output_file, err)
		return
	}
	format.Node(out, output_fileset, newAstFile)
	out.Close()
	fmt.Printf("  Wrote %s\n", output_file)
}


// Our internal represerntation of a parsed formation definition -- something
// we can hang convenience methods off of.
type formationDef struct {
	fset *token.FileSet
	ts *ast.TypeSpec
	fields []*ast.Field
}


const reader_method_template = `
package foo
func (f *STRUCT_TYPE) READER_NAME() FIELD_TYPE { return f.FIELD_NAME }
` // end


const formation_methods_template = `
package foo

// Code below expects this definition to be first.
func ( f * STRUCT_TYPE)nonRedundant() []Formation {
	return []Formation { }
}

func (f *STRUCT_TYPE) FORMATION_NAME() {}

func (f * STRUCT_TYPE) NumberOfDancers() int {
	count := 0
	for _, f := range f.nonRedundant() {
		count += f.NumberOfDancers()
	}
	return count
}

func (f * STRUCT_TYPE) Dancers() dancer.Dancers {
	dancers := dancer.Dancers {}
	for _, f := range f.nonRedundant() {
		dancers = append(dancers, f.Dancers()...) 
	}
	return dancers
}

func (f * STRUCT_TYPE) HasDancer(d dancer.Dancer) bool {
	for _, f := range f.nonRedundant() {
		if f.HasDancer(d) {
			return true
		}
	}
	return false
}

` // end


func (fdef *formationDef) generate(output_fileset *token.FileSet) (decls []ast.Decl) {
	implName := "impl" + fdef.ts.Name.Name
	struct_type := &ast.StructType { Fields: &ast.FieldList{ List: []*ast.Field{} } }
	decls = append(decls, &ast.GenDecl {
		Tok: token.TYPE,
		Specs: []ast.Spec {
			&ast.TypeSpec {
				Name: ast.NewIdent(implName),
				Type: struct_type,
			},
		},
	})
	// field reader methods
	nonRedundant := []string{}
	for _, field := range fdef.fields {
		var field_type ast.Expr
		if ft, ok := field.Type.(*ast.FuncType); ok {
			if ft.Results == nil || len(ft.Results.List) == 0 {
				continue
			}
			field_type = ft.Results.List[0].Type
		} else {
			continue
		}
		new_field := &ast.Field {
			Type: field_type,
		}
		noslot := commentHasKeyword(field.Comment, NoSlot)
		for _, name := range field.Names {
			field_name := strings.ToLower(name.Name)
			new_field.Names = append(new_field.Names, ast.NewIdent(field_name))
			if noslot {
				continue
			}
			reader := parseDefinition(reader_method_template)
			v := go_tools.NewSubstitutingVisitor()
			v.Substitutions["STRUCT_TYPE"] = implName
			v.Substitutions["READER_NAME"] = name.Name
			v.Substitutions["FIELD_TYPE"] = NodeString(field_type)
			v.Substitutions["FIELD_NAME"] = field_name
			ast.Walk(v, reader)
			decls = append(decls, reader.Decls...)
			if !commentHasKeyword(field.Comment, Redundant) {
				nonRedundant = append(nonRedundant, name.Name)
			}
		}
		if !noslot {
			struct_type.Fields.List = append(struct_type.Fields.List, new_field)
		}
	}
	// Formation methods
	{
		code := parseDefinition(formation_methods_template)
		v := go_tools.NewSubstitutingVisitor()
		v.Substitutions["STRUCT_TYPE"] = implName
		v.Substitutions["FORMATION_NAME"] = fdef.ts.Name.Name
		ast.Walk(v, code)
		result0slice := code.Decls[0].(*ast.FuncDecl).Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit)
		for _, reader := range nonRedundant {
			result0slice.Elts = append(result0slice.Elts,
				parseExpression(fmt.Sprintf("f.%s()", reader)))
		}
		decls = append(decls, code.Decls...)
	}
	return decls
}


func commentHasKeyword(comment_group *ast.CommentGroup, keyword string) bool {
	if comment_group == nil {
		return false
	}
	for _, c := range comment_group.List {
		if strings.Contains(c.Text, keyword) {
			return true
		}
	}
	return false
}

func IsFormationDefinition(fset *token.FileSet, decl ast.Decl) (bool, *formationDef) {
	gd, ok := decl.(*ast.GenDecl)
	if !ok {
		return false, nil
	}
	spec, ok := gd.Specs[0].(*ast.TypeSpec)
	if !ok {
		return false, nil
	}
	it, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return false, nil
	}
	if len(gd.Specs) > 1 {
		fmt.Fprintf(os.Stderr, "  Type definition of an interface type has more than one Spec: %s\n",
			fset.Position(gd.TokPos).String())
	}
	foundFormation := false
	fd := &formationDef {
		fset: fset,
		ts: spec,
	}
	fieldIsFormation := func(field *ast.Field) bool {
		return NodeString(field.Type) == "Formation"
	}
	for _, field := range it.Methods.List {
		if fieldIsFormation(field) {
			foundFormation = true
			continue
		}
		// Collect the fields that aren't the Formation interface.
		fd.fields = append(fd.fields, field)
	} 
	if !foundFormation {
		return false, nil
	}
	return true, fd
}


func parseExpression(exp string) ast.Expr {
	e, err := parser.ParseExpr(exp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Offending code:\n%s\n", exp)
		panic(err)
	}
	return e
}

func parseDefinition(def string) *ast.File {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", def, 0)
	if err != nil {
		panic(fmt.Sprintf("Errors:\n%s", err))
	}
	return astFile
}

func NodeString(n ast.Node) string {
	w := bytes.NewBufferString("")
	format.Node(w, token.NewFileSet(), n)
	return w.String()
}
