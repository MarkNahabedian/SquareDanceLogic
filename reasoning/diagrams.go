// Generate an HTML file showing diagrams for the Formations we can
// make samples for.

package reasoning

import "fmt"
import "os"
import "reflect"
import "sort"
import "html/template"
import "squaredance/dancer"
import "squaredance/geometry"


var DancerTemplateFunctions = template.FuncMap{
	"JSDirection": func(d geometry.Direction) float32 {
		return float32(d) * 4
	},
	"JSGender": func (gender dancer.Gender) string {
		switch gender {
		case dancer.Guy: return "guy"
		case dancer.Gal: return "gal"
		default: return "unspecified"
		}
	},
}


// DancersSVGTemplateArg is the type of parameter for the
// DancersSVGTemplate HTML template.
type DancersSVGTemplateArg interface {
	SVGId() string          //defimpl:"read svg_id"
	Sample() Formation     //defimpl:"read sample"
	Name() string           // The name of the formation
	HasSample() bool
	DancerCount() int
}

// dancersSVGTemplate is a template for generating the javascript code
// for drawing a set of dancers.  It shoud be called with a
// DancersSVGTemplateArg as argument.
var dancersSVGTemplate = template.Must(template.New("DancersSVGTemplate").
	Funcs(DancerTemplateFunctions).Parse(`
  {{- with $dsta := .}}
  new Floor([
      {{- range $dsta.Sample.Dancers -}}
	new Dancer({{.Position.Left}}, {{.Position.Down}}, {{JSDirection .Direction}}, "{{.Ordinal}}",
		   {{JSGender .Gender}}, "white", "{{.Ordinal}}"),
      {{end -}}
    ]).draw("{{$dsta.SVGId}}");
  {{- end -}}
`))

// DancersSVGTemplate returns an HTML Template that includes a
// Template named DancersSVGTemplate which will generate javascript
// code to draw the specified dancers in an SVG element when given a
// DancersSVGTemplateArg as argument.
func DancersSVGTemplate() *template.Template {
	t, err := dancersSVGTemplate.Clone()
	if err != nil {
		panic(err)
	}
	return t
}

// NewDancersSVGTemplateArg returns a minimal implementation of the
// DancersSVGTemplateArg.
func NewDancersSVGTemplateArg(formation_type reflect.Type) DancersSVGTemplateArg {
	return &DancersSVGTemplateArgImpl{
		svg_id: formation_type.Name(),
		sample: MakeSampleFormation(formation_type),
	}
}

func (dsta *DancersSVGTemplateArgImpl) Name() string {
	return dsta.svg_id
}

func (dsta *DancersSVGTemplateArgImpl) HasSample() bool {
	return dsta.sample != nil
}

func (dsta *DancersSVGTemplateArgImpl) DancerCount() int {
	if dsta.sample == nil {
		return -1
	}
	return len(dsta.sample.Dancers())
}

// templateFormationType implements the DancersSVGTemplateArg
// interface.
func (dsta *DancersSVGTemplateArgImpl) Dancers() dancer.Dancers {
	return dsta.sample.Dancers()
}


type FormationTypeSort []DancersSVGTemplateArg

func (fts FormationTypeSort) Len() int {
	return len(fts)
}

func (fts FormationTypeSort) Swap(i, j int) {
	fts[i], fts[j] = fts[j], fts[i]
}

func (fts FormationTypeSort) Less (i, j int) bool {
	// Sort first by number of dancers
	if fts[i].DancerCount() < fts[j].DancerCount() {
		return true
	}
	if fts[i].DancerCount() > fts[j].DancerCount() {
		return false
	}
	// then by Formation name
	return fts[i].Name() < fts[j].Name()
}


func WriteFormationDiagrams() {
	filename := "formation_types.html"
	// Sort
	fts := FormationTypeSort{}
	for _, ft := range AllFormationTypes {
		fts = append(fts, NewDancersSVGTemplateArg(ft))
	}
	sort.Sort(fts)
	// Create file
	out, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Can't open %s: %s", filename, err))
	}
	defer out.Close()
	// Generate HTML
	err = html_page.Execute(out, fts)
	if err != nil {
		fmt.Println(err)
		return
	}
}


func init() {
	child := DancersSVGTemplate()
	_, err := html_page.AddParseTree(child.Name(), child.Tree)
	if err != nil {
		panic(err)
	}
}

var html_page = template.Must(template.New("html_page").Funcs(
	DancerTemplateFunctions).Parse(`<html>
  <head>
    <title>
      Supported Formation Types
    </title>
    <style>
td {
  text-align: center;
  vertical-align: middle;
}
    </style>
    <script type="text/javascript"
            src="https://marknahabedian.github.io/SquareDanceFormationDiagrams/dancers.js">
    </script>
    <script type="text/javascript">
function contentLoaded() {
      {{range . -}}
        {{if .HasSample -}}
          {{- template "DancersSVGTemplate" . -}}
        {{end -}}
      {{- end -}}
}

document.addEventListener("DOMContentLoaded", contentLoaded, false);
    </script>
  </head>
  <body>
    <h1>
      Supported Square Dance Formation Types
    </h1>
    <p>
      This file is written by the function
      reasoning.WriteFormationDiagrams in
      <tt>squaredance/reasoning/diagrams.go</tt>.
      WriteFormationDiagrams is called by reasoning.TestWriteDiagrams
      defined in 
      <tt>squaredance/reasoning/reasoning_test.go</tt>.
    </p>
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Diagram</th>
        </tr>
      </thead>
      <tbody>
        {{range .}}
          <tr>
            <td>{{.Name}}</td>
            <td>
              {{- if .HasSample -}}
                <svg id="{{.SVGId}}"></svg>
              {{- end -}}
            </td>
          </tr>
        {{end}}
      </tbody>
    </table>
  </body>
</html>
`))


