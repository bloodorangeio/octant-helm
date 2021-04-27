/*
Copyright 2021 Blood Orange

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package views // import "github.com/bloodorangeio/octant-helm/pkg/plugin/views"

import (
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

func BuildRepoViewForRequest(_ service.Request) (component.Component, error) {
	f, err := repo.LoadFile(cli.New().RepositoryConfig)
	if err != nil {
		return nil, err
	}

	table := component.NewTable("Chart Repositories", "There are no repositories!", component.NewTableCols("Name", "URL"))
	for _, repo := range f.Repositories {
		table.Add(component.TableRow{
			"Name": component.NewText(repo.Name),
			"URL":  component.NewLink("", repo.URL, repo.URL),
		})
	}
	return table, nil
}
