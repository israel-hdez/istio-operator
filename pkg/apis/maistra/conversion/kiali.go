package conversion

import (
	v2 "github.com/maistra/istio-operator/pkg/apis/maistra/v2"
)

func populateKialiAddonValues(kiali *v2.KialiAddonConfig, values map[string]interface{}) error {
	if kiali == nil {
		return nil
	}
	if kiali.Install == nil {
		// we don't want to process the charts
		if err := setHelmBoolValue(values, "kiali.enabled", false); err != nil {
			return err
		}
		return nil
	}

	kialiValues := make(map[string]interface{})
	if kiali.Enabled != nil {
		if err := setHelmBoolValue(kialiValues, "enabled", *kiali.Enabled); err != nil {
			return err
		}
	}

	dashboardConfig := kiali.Install.Config.Dashboard
	if dashboardConfig.ViewOnly != nil {
		if err := setHelmBoolValue(kialiValues, "dashboard.viewOnlyMode", *dashboardConfig.ViewOnly); err != nil {
			return err
		}
	}
	if dashboardConfig.EnableGrafana != nil {
		if err := setHelmBoolValue(kialiValues, "dashboard.enableGrafana", *dashboardConfig.EnableGrafana); err != nil {
			return err
		}
	}
	if dashboardConfig.EnablePrometheus != nil {
		if err := setHelmBoolValue(kialiValues, "dashboard.enablePrometheus", *dashboardConfig.EnablePrometheus); err != nil {
			return err
		}
	}
	if dashboardConfig.EnableTracing != nil {
		if err := setHelmBoolValue(kialiValues, "dashboard.enableTracing", *dashboardConfig.EnableTracing); err != nil {
			return err
		}
	}
	ingressValues := make(map[string]interface{})
	if err := populateAddonIngressValues(kiali.Install.Service.Ingress, ingressValues); err == nil {
		if len(ingressValues) > 0 {
			if err := setHelmValue(kialiValues, "ingress", ingressValues); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	if kiali.Install.Runtime != nil {
		runtime := kiali.Install.Runtime
		if err := populateRuntimeValues(runtime, kialiValues); err != nil {
			return err
		}

		if kialiContainer, ok := runtime.Pod.Containers["kiali"]; ok {
			if kialiContainer.Image != "" {
				if err := setHelmStringValue(kialiValues, "image", kialiContainer.Image); err != nil {
					return err
				}
			}
			if kialiContainer.ImageRegistry != "" {
				if err := setHelmStringValue(kialiValues, "hub", kialiContainer.ImageRegistry); err != nil {
					return err
				}
			}
			if kialiContainer.ImageTag != "" {
				if err := setHelmStringValue(kialiValues, "tag", kialiContainer.ImageTag); err != nil {
					return err
				}
			}
			if kialiContainer.ImagePullPolicy != "" {
				if err := setHelmStringValue(kialiValues, "imagePullPolicy", string(kialiContainer.ImagePullPolicy)); err != nil {
					return err
				}
			}
		}
	}

	if len(kialiValues) > 0 {
		if err := setHelmValue(values, "kiali", kialiValues); err != nil {
			return err
		}
	}
	return nil
}
