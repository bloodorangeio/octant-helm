package config

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

var settings = cli.New()

// NewActionConfig creates a config client for Helm actions
func NewActionConfig(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	getter := settings.RESTClientGetter()

	if err := actionConfig.Init(getter, namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		return nil, err
	}

	return actionConfig, nil
}
