package reconciler

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	"golang.org/x/exp/slices"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type update = func() error

func sFnUpdate(_ context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	for _, f := range []func(v1alpha1.ApplicationConnectorSpec, []unstructured.Unstructured) error{
		updateCentralApplicationGateway,
	} {
		if err := f(s.instance.Spec, r.Objs); err != nil {
			return stopWithErrorAndNoRequeue(err)
		}
	}
	return switchState(sFnApply)
}

func envVarUpdate(envs []corev1.EnvVar, newEnv corev1.EnvVar) update {
	return func() error {
		envIndex := slices.IndexFunc(
			envs,
			func(env corev1.EnvVar) bool { return newEnv.Name == env.Name })
		// return error if env variable was not found
		if envIndex == -1 {
			return fmt.Errorf(`'%s' env variable: %w`, newEnv.Name, unstructured.ErrNotFound)
		}

		envs[envIndex] = newEnv
		return nil
	}
}

func updateCRA(d *appv1.Deployment, v v1alpha1.RuntimeAgentSpec) error {
	// find compass-runtime-agent container
	index := slices.IndexFunc(
		d.Spec.Template.Spec.Containers,
		func(c corev1.Container) bool { return c.Name == "compass-runtime-agent" })
	// return error if compass-runtime-agent container was not found
	if index == -1 {
		return fmt.Errorf("compass-runtime-agent container: %w", unstructured.ErrNotFound)
	}
	compassRtAgentEnvs := d.Spec.Template.Spec.Containers[index].Env
	// define all update functions
	fns := []update{
		envVarUpdate(
			compassRtAgentEnvs,
			corev1.EnvVar{
				Name:  v1alpha1.EnvRuntimeAgentControllerSyncPeriod,
				Value: v.ControllerSyncPeriod.Duration.String(),
			}),
		envVarUpdate(
			compassRtAgentEnvs,
			corev1.EnvVar{
				Name:  v1alpha1.EnvRuntimeAgentCertValidityRenevalThreshold,
				Value: v.CertValidityRenewalThreshold,
			}),
		envVarUpdate(
			compassRtAgentEnvs,
			corev1.EnvVar{
				Name:  v1alpha1.EnvRuntimeAgentMinimalCompassSyncTime,
				Value: v.MinConfigSyncTime.Duration.String(),
			}),
	}
	// perform update
	for _, f := range fns {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func updateCentralApplicationGateway(i v1alpha1.ApplicationConnectorSpec, objs []unstructured.Unstructured) error {
	u, err := unstructured.IsDeployment("central-application-gateway").First(objs)
	if err != nil {
		return err
	}

	if err := unstructured.Update(u, i.ApplicationGatewaySpec, updateCAG); err != nil {
		return err
	}
	return nil
}

func updateCAG(d *appv1.Deployment, v v1alpha1.AppGatewaySpec) error {
	// find compass-runtime-agent container
	index := slices.IndexFunc(
		d.Spec.Template.Spec.Containers,
		func(c corev1.Container) bool { return c.Name == "central-application-gateway" })
	// return error if compass-runtime-agent container was not found
	if index == -1 {
		return fmt.Errorf("central-application-gateway container: %w", unstructured.ErrNotFound)
	}
	cAppG8wayArgs := d.Spec.Template.Spec.Containers[index].Args

	// define all update functions
	fns := []update{
		argValueUpdate(cAppG8wayArgs, v1alpha1.ArgCentralAppGatewayRequestTimeout, fmt.Sprintf("%.0f", v.RequestTimeout.Seconds())),
		argValueUpdate(cAppG8wayArgs, v1alpha1.ArgCentralAppGatewayProxyTimeout, fmt.Sprintf("%.0f", v.ProxyTimeout.Seconds())),
	}
	// perform update
	for _, f := range fns {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func argValueUpdate(args []string, key string, newValue any) update {
	return func() error {
		// assume argument is distinct
		argIndex := slices.IndexFunc(args, func(s string) bool {
			return strings.HasPrefix(s, key)
		})
		// return error if env variable was not found
		if argIndex == -1 {
			return fmt.Errorf(`argument with key: '%s': %w `, key, unstructured.ErrNotFound)
		}

		args[argIndex] = fmt.Sprintf("%s=%v", key, newValue)
		return nil
	}
}
