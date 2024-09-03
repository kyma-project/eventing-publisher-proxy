// Package applicationtest provides utilities for Application testing.
package applicationtest

import (
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kymaappconnv1alpha1 "github.com/kyma-project/kyma/components/central-application-gateway/pkg/apis/applicationconnector/v1alpha1"
)

func NewApplication(name string, labels map[string]string) *kymaappconnv1alpha1.Application {
	return &kymaappconnv1alpha1.Application{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
}
