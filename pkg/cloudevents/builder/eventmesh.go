package builder

import (
	ceevent "github.com/cloudevents/sdk-go/v2/event"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/application"

	"github.com/kyma-project/eventing-manager/pkg/backend/cleaner"
	"github.com/kyma-project/eventing-manager/pkg/logger"
)

// Perform a compile-time check.
var _ CloudEventBuilder = &EventMeshBuilder{}

func NewEventMeshBuilder(prefix string, eventMeshNamespace string, cleaner cleaner.Cleaner,
	applicationLister *application.Lister, logger *logger.Logger,
) CloudEventBuilder {
	genericBuilder := GenericBuilder{
		typePrefix:        prefix,
		applicationLister: applicationLister,
		logger:            logger,
		cleaner:           cleaner,
	}

	return &EventMeshBuilder{
		genericBuilder:     &genericBuilder,
		eventMeshNamespace: eventMeshNamespace,
	}
}

func (emb *EventMeshBuilder) Build(event ceevent.Event) (*ceevent.Event, error) {
	ceEvent, err := emb.genericBuilder.Build(event)
	if err != nil {
		return nil, err
	}

	// set eventMesh namespace as event source (required by EventMesh)
	ceEvent.SetSource(emb.eventMeshNamespace)

	return ceEvent, err
}
