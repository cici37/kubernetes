/*
Copyright 2024 The Kubernetes Authors.

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

package mutating

import (
	celgo "github.com/google/cel-go/cel"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission/plugin/cel"
)

var _ cel.ExpressionAccessor = &MutationCondition{}

// MutationCondition contains the inputs needed to compile, evaluate and validate a cel expression
type MutationCondition struct {
	Expression string
	Message    string
	Reason     *metav1.StatusReason
}

func (v *MutationCondition) GetExpression() string {
	return v.Expression
}

func (v *MutationCondition) ReturnTypes() []*celgo.Type {
	return []*celgo.Type{celgo.BoolType}
}

// Variable is a named expression for composition.
type Variable struct {
	Name       string
	Expression string
}

func (v *Variable) GetExpression() string {
	return v.Expression
}

func (v *Variable) ReturnTypes() []*celgo.Type {
	return []*celgo.Type{celgo.AnyType, celgo.DynType}
}

func (v *Variable) GetName() string {
	return v.Name
}
