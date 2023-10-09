package reconciler

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/kyma-project/application-connector-manager/pkg/reconciler/mocks"
	"github.com/onsi/gomega/types"
	v1 "k8s.io/api/core/v1"
)

type testDetectDomainOptions struct {
	MatchExpectedErr      types.GomegaMatcher
	MatchNextFnState      types.GomegaMatcher
	MatchExpectedInstance types.GomegaMatcher
}

// isGardenerConfigMap - checks if object key matches gardener configmap key
func isGardenerConfigMap(n client.ObjectKey) bool {
	return gardenerCM == n
}

func mockClient(data map[string]string) client.Client {
	cmType := v1.ConfigMap{}

	k8sMock := &mocks.Client{}
	k8sMock.On(
		"Get",
		mock.Anything,
		mock.MatchedBy(isGardenerConfigMap),
		mock.IsType(&cmType),
	).Return(nil).Run(func(args mock.Arguments) {
		cm := args.Get(2).(*v1.ConfigMap)
		cm.Data = data
	})
	return k8sMock
}

var _ = Describe("ACM sFnDetectDomain", func() {
	var testDomainName = "test.domain.name"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	DescribeTable(
		"detect domain function",
		testDetectDomain,
		Entry("domain name set on instance",
			ctx,
			&fsm{},
			&systemState{
				instance: v1alpha1.ApplicationConnector{
					Spec: v1alpha1.ApplicationConnectorSpec{
						DomainName: testDomainName,
					},
				},
			},
			testDetectDomainOptions{
				MatchExpectedErr:      BeNil(),
				MatchNextFnState:      equalStateFunction(sFnUpdate),
				MatchExpectedInstance: BeEmpty(),
			},
		),
		Entry("domain fetched from CM not found",
			ctx,
			&fsm{
				K8s: K8s{Client: mockClient(map[string]string{})},
			},
			&systemState{
				instance: v1alpha1.ApplicationConnector{
					Spec: v1alpha1.ApplicationConnectorSpec{},
				},
			},
			testDetectDomainOptions{
				MatchExpectedErr:      BeNil(),
				MatchNextFnState:      equalStateFunction(sFnUpdateStatus(nil, nil)),
				MatchExpectedInstance: BeEmpty(),
			},
		),
		Entry("domain fetched from CM",
			ctx,
			&fsm{
				K8s: K8s{Client: mockClient(map[string]string{"domain": testDomainName})},
			},
			&systemState{
				instance: v1alpha1.ApplicationConnector{
					Spec: v1alpha1.ApplicationConnectorSpec{},
				},
			},
			testDetectDomainOptions{
				MatchExpectedErr:      BeNil(),
				MatchNextFnState:      equalStateFunction(sFnUpdate),
				MatchExpectedInstance: Equal(testDomainName),
			},
		),
		Entry("domain previously fetched",
			ctx,
			&fsm{},
			&systemState{
				domainName: testDomainName,
			},
			testDetectDomainOptions{
				MatchExpectedErr:      BeNil(),
				MatchNextFnState:      equalStateFunction(sFnUpdate),
				MatchExpectedInstance: Equal(testDomainName),
			},
		),
	)
})

func testDetectDomain(ctx context.Context, r *fsm, s *systemState, ops testDetectDomainOptions) {
	next, _, err := sFnDetectDomain(ctx, r, s)
	Expect(err).To(ops.MatchExpectedErr)
	Expect(next).To(ops.MatchNextFnState)
	Expect(s.domainName).To(ops.MatchExpectedInstance)
}
