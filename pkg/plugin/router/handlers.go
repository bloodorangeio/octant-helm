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

package router // import "github.com/bloodorangeio/octant-helm/pkg/plugin/router"

import (
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"github.com/bloodorangeio/octant-helm/pkg/plugin/views"
)

func rootHandler(request service.Request) (component.ContentResponse, error) {
	rootView, err := views.BuildRootViewForRequest(request)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	response := component.NewContentResponse(nil)
	response.Add(rootView)
	return *response, nil
}

func helmReleaseHandler(request service.Request) (component.ContentResponse, error) {
	helmReleaseView, title, err := views.BuildHelmReleaseViewForRequest(request)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	response := component.NewContentResponse(title)
	response.Add(helmReleaseView)
	return *response, nil
}
