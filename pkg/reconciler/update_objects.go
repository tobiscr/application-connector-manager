package reconciler

import (
	"context"
	"fmt"

	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	"golang.org/x/exp/slices"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type compassRtAgentDTO struct {
	syncPeriod metav1.Duration
}

func sFnUpdate(_ context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	u, err := unstructured.IsDeployment("compass-runtime-agent").First(r.Objs)
	if err != nil {
		return stopWithErrorAndNoRequeue(err)
	}

	d := compassRtAgentDTO{syncPeriod: s.instance.Spec.RuntimeAgentSpec.ControllerSyncPeriod}
	if err := unstructured.Update(u, d, updateSyncPeriod); err != nil {
		return stopWithErrorAndNoRequeue(err)
	}

	return switchState(sFnApply)
}

type update = func() error

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

func updateSyncPeriod(d *appv1.Deployment, dto compassRtAgentDTO) error {
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
				Name:  "APP_CONTROLLER_SYNC_PERIOD",
				Value: dto.syncPeriod.Duration.String(),
			}),
		envVarUpdate(
			compassRtAgentEnvs,
			corev1.EnvVar{
				Name:  "APP_MINIMAL_COMPASS_SYNC_TIME",
				Value: dto.syncPeriod.Duration.String(),
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
