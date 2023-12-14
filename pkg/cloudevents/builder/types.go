package builder

import (
	cev2event "github.com/cloudevents/sdk-go/v2/event"
	"github.com/kyma-project/eventing-manager/pkg/backend/cleaner"
	"github.com/kyma-project/eventing-manager/pkg/logger"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/application"
)

const (
	OriginalTypeHeaderName = "originaltype"
)

type CloudEventBuilder interface {
	Build(event cev2event.Event) (*cev2event.Event, error)
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
