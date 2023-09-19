package reconciler

import (
	"context"
	"testing"
	"time"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	modtest "github.com/kyma-project/application-connector-manager/pkg/reconciler/testing"
)

var (
	updateTimeout = time.Second * 5
)

func Test_sFnUpdate(t *testing.T) {
	// load all test data from testdata/update
	objs, err := modtest.LoadTestData(modtest.SfnUpdate)
	if err != nil {
		t.Fatalf("unable to extract test data: %s", err)
	}

	defaultState := &systemState{
		instance: v1alpha1.ApplicationConnector{
			Spec: v1alpha1.ApplicationConnectorSpec{
				SyncPeriod: "10s",
			},
		},
	}

	type args struct {
		r *fsm
		s *systemState
	}

	tests := []struct {
		name           string
		args           args
		wantNextFnName string
		wantErr        bool
	}{
		{
			name: "happy path",
			args: args{
				s: defaultState,
				r: &fsm{
					Cfg: Cfg{
						Objs: objs[modtest.TdUpdateAcmValid],
					},
				},
			},
			wantNextFnName: modtest.NamesFnApply,
		},
		{
			name: "missing deployment",
			args: args{
				s: defaultState,
				r: &fsm{},
			},
			wantNextFnName: modtest.NamesFnUpdateStatus,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), updateTimeout)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFn, _, err := sFnUpdate(ctx, tt.args.r, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("sFnUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNextFnName != sFn.name() {
				t.Errorf("sFnUpdate() sFn = %s, want %s", sFn.name(), tt.wantNextFnName)
			}
		})
	}
}
