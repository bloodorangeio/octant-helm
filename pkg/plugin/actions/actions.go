package actions

import (
	"fmt"
	"github.com/bloodorangeio/octant-helm/pkg/config"
	"github.com/bloodorangeio/octant-helm/pkg/helm"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	UpdateHelmReleaseValues = "octant-helm.dev/update"
	UninstallHelmReleaseAction = "octant-helm.dev/uninstall"
)

func ActionHandler(request *service.ActionRequest) error {
	actionConfig, err := config.NewActionConfig(request.ClientState.Namespace)
	if err != nil {
		return err
	}

	actionName, err := request.Payload.String("action")
	if err != nil {
		return err
	}

	switch actionName {
	case UninstallHelmReleaseAction:
		releaseName, err := request.Payload.String("releaseName")
		if err != nil {
			return err
		}
		return uninstallRelease(request, actionConfig, releaseName)
	case UpdateHelmReleaseValues:
		releaseValues, err := request.Payload.String("update")
		if err != nil {
			return err
		}
		releaseName, err := request.Payload.String("releaseName")
		if err != nil {
			return err
		}
		return updateReleaseValues(request, actionConfig, releaseValues, releaseName)
	default:
		return fmt.Errorf("unable to find handler for plugin: %s", "octant-helm")
	}
}

func uninstallRelease(request *service.ActionRequest, config *helmAction.Configuration, releaseName string) error {
	uninstallClient := helmAction.NewUninstall(config)
	release, err := uninstallClient.Run(releaseName)
	if err != nil {
		return err
	}

	alert := action.CreateAlert(action.AlertTypeInfo, "Uninstalled helm release: "+release.Release.Name, action.DefaultAlertExpiration)
	request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, alert)
	return nil
}

func updateReleaseValues(request *service.ActionRequest, config *helmAction.Configuration, values string, releaseName string) error {
	upgradeClient := helmAction.NewUpgrade(config)

	client := request.DashboardClient
	ul, err := client.List(request.Context(), store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Selector: &labels.Set{
			"owner": "helm",
			"name":  releaseName,
		},
	})
	if err != nil {
		return err
	}
	r := helm.UnstructuredListToHelmReleaseByName(ul, releaseName)
	if r == nil {
		return fmt.Errorf("cannot find release name: %s", releaseName)
	}
	v, err := chartutil.ReadValues([]byte(values))
	if err != nil {
		return err
	}
	// TODO: There are upgrades which require secrets. How would this be shown to a user?
	release, err := upgradeClient.Run(releaseName, r.Chart, v)
	if err != nil {
		message := fmt.Sprintf("Unable to upgrade release: %s", err)
		request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, action.CreateAlert(action.AlertTypeError, message, action.DefaultAlertExpiration))
	} else {
		alert := action.CreateAlert(action.AlertTypeInfo, "Upgrade helm release: "+release.Name, action.DefaultAlertExpiration)
		request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, alert)
	}
	return nil
}
