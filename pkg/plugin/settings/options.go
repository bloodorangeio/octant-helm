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

package settings // import "github.com/bloodorangeio/octant-helm/pkg/plugin/settings"

import (
	"strings"

	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	"github.com/bloodorangeio/octant-helm/pkg/plugin/actions"
	"github.com/bloodorangeio/octant-helm/pkg/plugin/router"
)

func GetOptions() []service.PluginOption {
	return []service.PluginOption{
		service.WithActionHandler(actions.ActionHandler),
		service.WithNavigation(
			func(request *service.NavigationRequest) (navigation.Navigation, error) {
				return navigation.Navigation{
					Title:    strings.Title(name),
					Path:     name,
					IconName: rootNavIcon,
					Children: []navigation.Navigation{
						{
							Title:    "Repositories",
							Path:     request.GeneratePath("repositories"),
							IconName: "folder",
						},
						{
							Title:    "Environment",
							Path:     request.GeneratePath("environment"),
							IconName: "cog",
						},
					},
				}, nil
			},
			router.InitRoutes,
		),
	}
}
