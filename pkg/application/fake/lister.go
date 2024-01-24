package fake

import (
	"context"
	"log"

	kcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kdynamicfake "k8s.io/client-go/dynamic/fake"

	kymaappconnv1alpha1 "github.com/kyma-project/kyma/components/central-application-gateway/pkg/apis/applicationconnector/v1alpha1"

	"github.com/kyma-project/eventing-publisher-proxy/pkg/application"
)

func NewApplicationListerOrDie(ctx context.Context, app *kymaappconnv1alpha1.Application) *application.Lister {
	scheme := setupSchemeOrDie()
	dynamicClient := kdynamicfake.NewSimpleDynamicClient(scheme, app)
	return application.NewLister(ctx, dynamicClient)
}

func setupSchemeOrDie() *runtime.Scheme {
	scheme := runtime.NewScheme()
	if err := kcorev1.AddToScheme(scheme); err != nil {
		log.Fatalf("Failed to setup scheme with error: %v", err)
	}
	if err := kymaappconnv1alpha1.AddToScheme(scheme); err != nil {
		log.Fatalf("Failed to setup scheme with error: %v", err)
	}
	return scheme
}
