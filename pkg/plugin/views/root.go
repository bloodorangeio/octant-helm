/*
Copyright 2019 Blood Orange

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
	"fmt"
	"strconv"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/bloodorangeio/octant-helm/pkg/helm"
)

func BuildRootViewForRequest(request service.Request) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	ul, err := client.List(ctx, store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Selector: &labels.Set{
			"owner": "helm",
		},
	})

	if err != nil {
		return nil, err
	}

	helmReleases := helm.UnstructuredListToHelmReleaseList(ul)

	header := component.NewMarkdownText(fmt.Sprintf("## Helm"))

	table := component.NewTableWithRows(
		"Releases", "There are no Helm releases!",
		component.NewTableCols("Name", "Namespace", "Revision", "Updated", "Status", "Chart", "App Version"),
		[]component.TableRow{})

	for _, r := range helmReleases {
		tr := component.TableRow{
			"Name":      component.NewLink(r.Name, r.Name, r.Name),
			"Namespace": component.NewText(r.Namespace),
			"Revision":  component.NewText(strconv.Itoa(r.Version)),
			"Status":    component.NewText(r.Info.Status.String()),
			"Chart": component.NewText(
				fmt.Sprintf("%s-%s", r.Chart.Metadata.Name, r.Chart.Metadata.Version)),
			"App Version": component.NewText(r.Chart.Metadata.AppVersion),
		}
		tr["Updated"] = component.NewTimestamp(r.Info.LastDeployed.Time)
		table.Add(tr)
	}

	table.Sort("Name")

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: header},
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}
