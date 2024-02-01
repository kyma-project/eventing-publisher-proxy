package application

import (
	"context"
	"errors"
	"time"

	"github.com/kyma-project/eventing-publisher-proxy/pkg/informers"
	kymaappconnv1alpha1 "github.com/kyma-project/kyma/components/central-application-gateway/pkg/apis/applicationconnector/v1alpha1"
	kcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"

	emlogger "github.com/kyma-project/eventing-manager/pkg/logger"
)

var ErrFailedToConvertObjectToUnstructured = errors.New("failed to convert runtime object to unstructured")

type Lister struct {
	lister cache.GenericLister
}

func NewLister(ctx context.Context, client dynamic.Interface) *Lister {
	const defaultResync = 10 * time.Second
	gvr := GroupVersionResource()
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, defaultResync, kcorev1.NamespaceAll, nil)
	factory.ForResource(gvr)
	lister := factory.ForResource(gvr).Lister()
	logger, _ := emlogger.New("json", "error")
	informers.WaitForCacheSyncOrDie(ctx, factory, logger)
	return &Lister{lister: lister}
}

func (l Lister) Get(name string) (*kymaappconnv1alpha1.Application, error) {
	object, err := l.lister.Get(name)
	if err != nil {
		return nil, err
	}

	u, ok := object.(*unstructured.Unstructured)
	if !ok {
		return nil, ErrFailedToConvertObjectToUnstructured
	}

	a := &kymaappconnv1alpha1.Application{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, a); err != nil {
		return nil, err
	}

	return a, nil
}

func GroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    kymaappconnv1alpha1.SchemeGroupVersion.Group,
		Version:  kymaappconnv1alpha1.SchemeGroupVersion.Version,
		Resource: "applications",
	}
}
