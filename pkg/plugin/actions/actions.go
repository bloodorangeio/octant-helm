package actions

import (
	"fmt"
	"github.com/bloodorangeio/octant-helm/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	helmAction "helm.sh/helm/v3/pkg/action"
)

const (
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
