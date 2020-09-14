// Produce a catalog of all of the actions at a given level.

package action

import "fmt"
import "os"
import "sort"
import "html/template"


type catalogSort []FormationAction

func (cs catalogSort) Len() int {
	return len(cs)
}

func (cs catalogSort) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

func (cs catalogSort) Less (i, j int) bool {
	// Sort first by Action Name
	if cs[i].Action().Name() < cs[j].Action().Name() {
		return true
	}
	if cs[i].Action().Name() > cs[j].Action().Name() {
		return false
	}
	// then by Formation name
	return cs[i].FormationType().Name() < cs[j].FormationType().Name()
}


func WriteCatalog(level Level) {
	fas := []FormationAction{}
	for _, action := range AllActions {
		action.DoFormationActions(func(fa FormationAction) bool {
			if fa.Level() == level {
				fas = append(fas, fa)
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
		FormationActions: fas,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

type html_page_arg struct {
	Level Level
	FormationActions []FormationAction
}

// Parameters are level and sortes slice of FormationAction.
var html_page = template.Must(template.New("html_page").Parse(`<html>
  <head>
    <title>
      Catalog of {{.Level}} level Formation Actions
    </title>
  </head>
  <body>
    <h1>
      Catalog of {{.Level}} level Formation Actions
    </h1>
    <table>
      <tr>
        <th>Action</th>
        <th>Formation</th>
        <th>Before</th>
        <th>After</th>
      </tr>
      {{range .FormationActions}}
        <tr>
          <td>{{.Action.Name}}</td>
          <td>{{.FormationType.Name}}</td>
          <td></td>
          <td></td>
        </tr>
      {{end}}
    </table>
  </body>
</html>
`))

