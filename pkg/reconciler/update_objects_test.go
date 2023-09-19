package reconciler

import (
	"context"
	"reflect"
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	updateTimeout = time.Second * 5
)

func Test_sFnUpdate(t *testing.T) {
	type args struct {
		r *fsm
		s *systemState
	}
	tests := []struct {
		name    string
		args    args
		want    stateFn
		want1   *ctrl.Result
		wantErr bool
	}{}

	ctx, cancel := context.WithTimeout(context.Background(), updateTimeout)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := sFnUpdate(ctx, tt.args.r, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("sFnUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sFnUpdate() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("sFnUpdate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
