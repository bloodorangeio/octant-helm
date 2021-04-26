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
	"github.com/bloodorangeio/octant-helm/pkg/config"
	"github.com/bloodorangeio/octant-helm/pkg/plugin/actions"
	helmAction "helm.sh/helm/v3/pkg/action"
	//"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"

	//"helm.sh/helm/v3/pkg/cli"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/bloodorangeio/octant-helm/pkg/helm"
)

func BuildHelmReleaseViewForRequest(request service.Request) (component.Component, []component.TitleComponent, error) {
	releaseName := strings.TrimPrefix(request.Path(), "/")

	ctx := request.Context()
	client := request.DashboardClient()

	ul, err := client.List(ctx, store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Selector: &labels.Set{
			"owner": "helm",
			"name":  releaseName,
		},
	})

	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	r := helm.UnstructuredListToHelmReleaseByName(ul, releaseName)
	if r == nil {
		return component.NewText("Error: release not found"), nil, nil
	}

	title := component.Title(component.NewLink("", "Helm", "/helm"))
	title = append(title, component.NewText(releaseName))

	statusSummarySections := []component.SummarySection{
		{"Name", component.NewText(r.Name)},
		{"Last Deployed", component.NewText(r.Info.LastDeployed.Format(time.ANSIC))},
		{"Namespace", component.NewText(r.Namespace)},
		{"Status", component.NewText(r.Info.Status.String())},
		{"Revision", component.NewText(fmt.Sprintf("%d", r.Version))},
	}

	statusSummary := component.NewSummary("Status", statusSummarySections...)

	actionConfig, err := config.NewActionConfig(request.ClientState().Namespace)
	if err != nil {
		return nil, nil, err
	}
	historyClient := helmAction.NewHistory(actionConfig)
	history, err := historyClient.Run(r.Name)
	if err != nil {
		return nil, nil, err
	}
	historyColumns := component.NewTableCols("Revision", "Updated", "Status", "Chart", "App Version", "Description")
	historyTable := component.NewTable("History", "There is no history!", historyColumns)
	for i := len(history)-1; i >= 0; i-- {
		var appVersion string
		h := history[i]
		if h.Chart.Metadata != nil {
			appVersion = h.Chart.Metadata.Version
		}
		historyTable.Add(component.TableRow{
			"Revision":    component.NewText(strconv.Itoa(h.Version)),
			"Updated":     component.NewTimestamp(h.Info.LastDeployed.Time),
			"Status":      component.NewText(h.Info.Status.String()),
			"Chart":       component.NewText(fmt.Sprintf("%s-%s", h.Chart.Metadata.Name, h.Chart.Metadata.Version)),
			"App Version": component.NewText(appVersion),
			"Description": component.NewText(h.Info.Description),
		})
	}

	notesCard := component.NewCard(component.TitleFromString("Notes"))
	notesBody := component.NewMarkdownText(fmt.Sprintf("```\n%s\n```", strings.TrimSpace(r.Info.Notes)))
	notesCard.SetBody(notesBody)

	flexLayout := component.NewFlexLayout("Summary")
	flexLayout.SetAccessor("summary")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthHalf, View: statusSummary},
		{Width: component.WidthFull, View: historyTable},
		{Width: component.WidthFull, View: notesCard},
	})

	return flexLayout, title, nil
}

func BuildHelmReleaseViewForValues(request service.Request) (component.Component, error) {
	releaseName := strings.TrimPrefix(request.Path(), "/")
	client := request.DashboardClient()
	ul, err := client.List(request.Context(), store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Selector: &labels.Set{
			"owner": "helm",
			"name":  releaseName,
		},
	})
	if err != nil {
		return component.NewError(component.TitleFromString("Cannot list Helm releases: "), err), nil
	}

	r := helm.UnstructuredListToHelmReleaseByName(ul, releaseName)
	if r == nil {
		return component.NewText("Error: release not found"), nil
	}


	actionConfig, err := config.NewActionConfig(request.ClientState().Namespace)
	if err != nil {
		return component.NewError(component.TitleFromString("Create Helm config: "), err), nil
	}
	valuesClient := helmAction.NewGetValues(actionConfig)
	values, err := valuesClient.Run(r.Name)
	if err != nil {
		return component.NewError(component.TitleFromString("Get values: "), err), nil
	}

	// No user supplied values so show all supplied defaults instead
	if len(values) == 0 {
		values = r.Chart.Values
	}
	v, err := yaml.Marshal(values)
	if err != nil {
		return component.NewError(component.TitleFromString("Error"), err), nil
	}

	editor := component.NewEditor(component.TitleFromString("Values"), string(v), false)
	editor.Config.Language = "yaml"
	editor.Config.Metadata = map[string]string{
		"releaseName": releaseName,
	}
	editor.SetAccessor("chart")
	editor.Config.SubmitAction = actions.UpdateHelmReleaseValues
	return editor, nil
}
