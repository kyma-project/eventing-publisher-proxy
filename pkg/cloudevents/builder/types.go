package builder

import (
	ceeventv2 "github.com/cloudevents/sdk-go/v2/event"
	"github.com/kyma-project/eventing-manager/pkg/backend/cleaner"
	"github.com/kyma-project/eventing-manager/pkg/logger"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/application"
)

const (
	OriginalTypeHeaderName = "originaltype"
)

type CloudEventBuilder interface {
	Build(event ceeventv2.Event) (*ceeventv2.Event, error)
}

type GenericBuilder struct {
	typePrefix        string
	applicationLister *application.Lister // applicationLister will be nil when disabled.
	cleaner           cleaner.Cleaner
	logger            *logger.Logger
}

type EventMeshBuilder struct {
	genericBuilder     *GenericBuilder
	eventMeshNamespace string
}
