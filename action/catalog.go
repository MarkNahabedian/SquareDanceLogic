// Produce a catalog of all of the actions at a given level.

package action

import "fmt"
import "os"
import "sort"
import "html/template"
import "squaredance/reasoning"


type catalogSort []*dancersTemplateArg

func (cs catalogSort) Len() int {
	return len(cs)
}

func (cs catalogSort) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

func (cs catalogSort) Less (i, j int) bool {
	// Sort first by Action Name
	fa1 := cs[i].FormationAction
	fa2 := cs[j].FormationAction
	if fa1.Action().Name() < fa2.Action().Name() {
		return true
	}
	if fa1.Action().Name() > fa2.Action().Name() {
		return false
	}
	// then by Formation name
	return fa1.FormationType().Name() < fa2.FormationType().Name()
}


// WriteCatalog writes an HTML file listing all FormationActions for
// the specified level.
func WriteCatalog(level Level) {
	fas := []*dancersTemplateArg{}
	// Filter by level:
	for _, action := range AllActions {
		action.DoFormationActions(func(fa FormationAction) bool {
			if fa.Level() == level {
				fas = append(fas, newDancersTemplateArg(fa))
			}
			return true
		})
	}
	sort.Sort(catalogSort(fas))
	f, err := os.Create(fmt.Sprintf("catalog-%v.html", level))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	err = html_page.Execute(f, html_page_arg {
		Level: level,
		DancersTemplateArgs: fas,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

type html_page_arg struct {
	Level Level
	DancersTemplateArgs []*dancersTemplateArg
}


func init() {
	child := reasoning.DancersSVGTemplate()
	_, err := html_page.AddParseTree(child.Name(), child.Tree)
	if err != nil {
		panic(err)
	}
	html_page.Funcs(reasoning.DancerTemplateFunctions)
}


type dancersTemplateArg struct {
	FormationAction FormationAction
	svg_id string
	sample reasoning.Formation
}

func newDancersTemplateArg(fa FormationAction) *dancersTemplateArg {
	return &dancersTemplateArg {
		FormationAction: fa,
		svg_id: fmt.Sprintf("%s-%d-%s-start",
			fa.Action().Name(),
			fa.Level(),
			fa.FormationType().Name()),
		sample: reasoning.MakeSampleFormation(fa.FormationType()),
	}
}

func (dta *dancersTemplateArg) SVGId() string {
	return dta.svg_id
}

func (dta *dancersTemplateArg) Sample() reasoning.Formation {
	return dta.sample
}

func (dta *dancersTemplateArg) HasSample() bool {
	return dta.sample != nil
}

func (dta *dancersTemplateArg) DancerCount() int {
	if dta.sample == nil {
		return -1
	}
	return len(dta.sample.Dancers())
}

func (dta *dancersTemplateArg) Name() string {
	return dta.FormationAction.FormationType().Name()
}


// Parameters are level and sorted slice of dancersTemplateArg.
var html_page = template.Must(template.New("html_page").Funcs(
	template.FuncMap{
		"NewDancersSVGTemplateArg":
		func (fa FormationAction) reasoning.DancersSVGTemplateArg {
			return newDancersTemplateArg(fa)
		},
	}).Parse(`<html>
  <head>
    <title>
      Catalog of {{.Level}} level Formation Actions
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
      {{range .DancersTemplateArgs -}}
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
      Catalog of {{.Level}} level Formation Actions
    </h1>
    <table>
      <thead>
        <tr>
          <th>Action</th>
          <th>Formation</th>
          <th>Before</th>
          <th>After</th>
        </tr>
      </thead>
      {{range .DancersTemplateArgs -}}
        <tr>
          <td>{{.FormationAction.Action.Name}}</td>
          <td>{{.FormationAction.FormationType.Name}}</td>
          <td>
            {{if .Sample -}}
              <svg id="{{.SVGId}}"></svg>
            {{- else -}}
              <span>
                {{printf "%s" .FormationAction.FormationType.Name}}
              </span>
            {{- end}}
          </td>
          <td></td>
        </tr>
      {{- end}}
    </table>
  </body>
</html>
`))

