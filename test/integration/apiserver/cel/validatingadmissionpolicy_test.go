/*
Copyright 2022 The Kubernetes Authors.

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

package cel

import (
	"context"
	"testing"
	"time"

	genericfeatures "k8s.io/apiserver/pkg/features"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
	apiservertesting "k8s.io/kubernetes/cmd/kube-apiserver/app/testing"
	"k8s.io/kubernetes/test/integration/framework"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

	v1 "k8s.io/api/core/v1"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// Test_ValidateNamespace_NoParams tests a ValidatingAdmissionPolicy that validates creation of a Namespace with no params.
func Test_ValidateNamespace_NoParams(t *testing.T) {
	failurePolicy := admissionregistrationv1alpha1.Fail
	ignorePolicy := admissionregistrationv1alpha1.Ignore
	forbiddenReason := metav1.StatusReasonForbidden

	testcases := []struct {
		name          string
		policy        *admissionregistrationv1alpha1.ValidatingAdmissionPolicy
		policyBinding *admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding
		namespace     *v1.Namespace
		err           string
		failureReason metav1.StatusReason
	}{
		{
			name: "namespace name contains suffix enforced by validating admission policy, using object metadata fields",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "object.metadata.name.endsWith('k8s')",
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					PolicyName: "validate-namespace-suffix",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-k8s",
				},
			},
			err: "",
		},
		{
			name: "namespace name does NOT contain suffix enforced by validating admission policyusing, object metadata fields",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "object.metadata.name.endsWith('k8s')",
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					PolicyName: "validate-namespace-suffix",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-foobar",
				},
			},
			err:           "ValidatingAdmissionPolicy 'validate-namespace-suffix' with binding 'validate-namespace-suffix-binding' denied request: failed expression: object.metadata.name.endsWith('k8s')",
			failureReason: metav1.StatusReasonInvalid,
		},
		{
			name: "namespace name does NOT contain suffix enforced by validating admission policy using object metadata fields, AND validating expression returns StatusReasonForbidden",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "object.metadata.name.endsWith('k8s')",
							Reason:     &forbiddenReason,
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					PolicyName: "validate-namespace-suffix",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "forbidden-test-foobar",
				},
			},
			err:           "ValidatingAdmissionPolicy 'validate-namespace-suffix' with binding 'validate-namespace-suffix-binding' denied request: failed expression: object.metadata.name.endsWith('k8s')",
			failureReason: metav1.StatusReasonForbidden,
		},
		{
			name: "namespace name contains suffix enforced by validating admission policy, using request field",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "request.name.endsWith('k8s')",
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					PolicyName: "validate-namespace-suffix",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-k8s",
				},
			},
			err: "",
		},
		{
			name: "namespace name does NOT contains suffix enforced by validating admission policy, using request field",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "request.name.endsWith('k8s')",
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					PolicyName: "validate-namespace-suffix",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-k8s",
				},
			},
			err: "",
		},
		{
			name: "runtime error when validating namespace, but failurePolicy=Ignore",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "object.nonExistentProperty == 'someval'",
						},
					},
					FailurePolicy: &ignorePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					PolicyName: "validate-namespace-suffix",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-k8s",
				},
			},
			err: "",
		},
		{
			name: "runtime error when validating namespace, but failurePolicy=Fail",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "object.nonExistentProperty == 'someval'",
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					PolicyName: "validate-namespace-suffix",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-k8s",
				},
			},
			err:           "ValidatingAdmissionPolicy 'validate-namespace-suffix' with binding 'validate-namespace-suffix-binding' denied request: expression 'object.nonExistentProperty == 'someval'' resulted in error: no such key: nonExistentProperty",
			failureReason: metav1.StatusReasonInvalid,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, genericfeatures.CELValidatingAdmission, true)()
			server, err := apiservertesting.StartTestServer(t, nil, []string{
				"--enable-admission-plugins", "ValidatingAdmissionPolicy",
			}, framework.SharedEtcd())
			if err != nil {
				t.Fatal(err)
			}
			defer server.TearDownFn()

			config := server.ClientConfig

			client, err := clientset.NewForConfig(config)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicies().Create(context.TODO(), testcase.policy, metav1.CreateOptions{}); err != nil {
				t.Fatal(err)
			}

			if _, err := client.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicyBindings().Create(context.TODO(), testcase.policyBinding, metav1.CreateOptions{}); err != nil {
				t.Fatal(err)
			}

			// TODO: add retry logic instead
			time.Sleep(time.Second)

			_, err = client.CoreV1().Namespaces().Create(context.TODO(), testcase.namespace, metav1.CreateOptions{})
			if err == nil && testcase.err == "" {
				return
			}

			if err == nil && testcase.err != "" {
				t.Logf("actual error: %v", err)
				t.Logf("expected error: %v", testcase.err)
				t.Fatal("got nil error but expected an error")
			}

			if err != nil && testcase.err == "" {
				t.Logf("actual error: %v", err)
				t.Logf("expected error: %v", testcase.err)
				t.Fatal("got error but expected none")
			}

			if err.Error() != testcase.err {
				t.Logf("actual validation error: %v", err)
				t.Logf("expected validation error: %v", testcase.err)
				t.Error("unexpected validation error")
			}

			checkFailureReason(t, err, testcase.failureReason)
		})
	}
}

// Test_ValidateNamespace_WithConfigMapParams tests a ValidatingAdmissionPolicy that validates creation of a Namespace,
// using ConfigMap as a param reference.
func Test_ValidateNamespace_WithConfigMapParams(t *testing.T) {
	failurePolicy := admissionregistrationv1alpha1.Fail
	// ignorePolicy := admissionregistrationv1alpha1.Ignore

	testcases := []struct {
		name          string
		policy        *admissionregistrationv1alpha1.ValidatingAdmissionPolicy
		policyBinding *admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding
		configMap     *v1.ConfigMap
		namespace     *v1.Namespace
		err           string
		failureReason metav1.StatusReason
	}{
		{
			name: "namespace name contains suffix enforced by validating admission policy",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					ParamKind: &admissionregistrationv1alpha1.ParamKind{
						APIVersion: "v1",
						Kind:       "ConfigMap",
					},
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "object.metadata.name.endsWith(params.data.namespaceSuffix)",
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					ParamRef: &admissionregistrationv1alpha1.ParamRef{
						Name:      "validate-namespace-suffix-param",
						Namespace: "default",
					},
					PolicyName: "validate-namespace-suffix",
				},
			},
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "validate-namespace-suffix-param",
					Namespace: "default",
				},
				Data: map[string]string{
					"namespaceSuffix": "k8s",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-k8s",
				},
			},
			err: "",
		},
		{
			name: "namespace name does NOT contain suffix enforced by validating admission policy",
			policy: &admissionregistrationv1alpha1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicySpec{
					ParamKind: &admissionregistrationv1alpha1.ParamKind{
						APIVersion: "v1",
						Kind:       "ConfigMap",
					},
					MatchConstraints: &admissionregistrationv1alpha1.MatchResources{
						ResourceRules: []admissionregistrationv1alpha1.NamedRuleWithOperations{
							{
								RuleWithOperations: admissionregistrationv1alpha1.RuleWithOperations{
									Operations: []admissionregistrationv1.OperationType{
										"CREATE",
									},
									Rule: admissionregistrationv1.Rule{
										APIGroups: []string{
											"",
										},
										APIVersions: []string{
											"*",
										},
										Resources: []string{
											"namespaces",
										},
									},
								},
							},
						},
					},
					Validations: []admissionregistrationv1alpha1.Validation{
						{
							Expression: "object.metadata.name.endsWith(params.data.namespaceSuffix)",
						},
					},
					FailurePolicy: &failurePolicy,
				},
			},
			policyBinding: &admissionregistrationv1alpha1.ValidatingAdmissionPolicyBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "validate-namespace-suffix-binding",
				},
				Spec: admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingSpec{
					ParamRef: &admissionregistrationv1alpha1.ParamRef{
						Name:      "validate-namespace-suffix-param",
						Namespace: "default",
					},
					PolicyName: "validate-namespace-suffix",
				},
			},
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "validate-namespace-suffix-param",
					Namespace: "default",
				},
				Data: map[string]string{
					"namespaceSuffix": "k8s",
				},
			},
			namespace: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-foo",
				},
			},
			err:           "ValidatingAdmissionPolicy 'validate-namespace-suffix' with binding 'validate-namespace-suffix-binding' denied request: failed expression: object.metadata.name.endsWith(params.data.namespaceSuffix)",
			failureReason: metav1.StatusReasonInvalid,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, genericfeatures.CELValidatingAdmission, true)()
			server, err := apiservertesting.StartTestServer(t, nil, []string{
				"--enable-admission-plugins", "ValidatingAdmissionPolicy",
			}, framework.SharedEtcd())
			if err != nil {
				t.Fatal(err)
			}
			defer server.TearDownFn()

			config := server.ClientConfig

			client, err := clientset.NewForConfig(config)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.CoreV1().ConfigMaps("default").Create(context.TODO(), testcase.configMap, metav1.CreateOptions{}); err != nil {
				t.Fatal(err)
			}

			if _, err := client.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicies().Create(context.TODO(), testcase.policy, metav1.CreateOptions{}); err != nil {
				t.Fatal(err)
			}

			if _, err := client.AdmissionregistrationV1alpha1().ValidatingAdmissionPolicyBindings().Create(context.TODO(), testcase.policyBinding, metav1.CreateOptions{}); err != nil {
				t.Fatal(err)
			}

			// TODO: add retry logic instead
			time.Sleep(time.Second)

			_, err = client.CoreV1().Namespaces().Create(context.TODO(), testcase.namespace, metav1.CreateOptions{})
			if err == nil && testcase.err == "" {
				return
			}

			if err == nil && testcase.err != "" {
				t.Logf("actual error: %v", err)
				t.Logf("expected error: %v", testcase.err)
				t.Fatal("got nil error but expected an error")
			}

			if err != nil && testcase.err == "" {
				t.Logf("actual error: %v", err)
				t.Logf("expected error: %v", testcase.err)
				t.Fatal("got error but expected none")
			}

			if err.Error() != testcase.err {
				t.Logf("actual validation error: %v", err)
				t.Logf("expected validation error: %v", testcase.err)
				t.Error("unexpected validation error")
			}

			checkFailureReason(t, err, testcase.failureReason)
		})
	}
}

func checkFailureReason(t *testing.T, err error, expectedReason metav1.StatusReason) {
	reason := err.(apierrors.APIStatus).Status().Reason
	if reason != expectedReason {
		t.Logf("actual error reason: %v", reason)
		t.Logf("expected failure reason: %v", expectedReason)
		t.Error("unexpected error reason")
	}
}
