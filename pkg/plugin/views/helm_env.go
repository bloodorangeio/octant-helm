package views

import (
	"sort"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"helm.sh/helm/v3/pkg/cli"
)

func BuildHelmEnvViewForRequest(_ service.Request) (component.Component, error) {
	table := component.NewTable("Environment Variables", "There are no env vars!", component.NewTableCols("Name", "Value"))
	// Loop vars and add to table as rows
	envVars := cli.New().EnvVars()
	sortedEnvs := getSortedEnvVarKeys(envVars)
	for _, envName := range sortedEnvs {
		row := component.TableRow{}
		row["Name"] = component.NewText(envName)
		row["Value"] = component.NewText(envVars[envName])
		table.Add(row)
	}
	return table, nil
}

func getSortedEnvVarKeys(envVars map[string]string) []string {
	var keys []string
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
