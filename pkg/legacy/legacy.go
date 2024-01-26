package legacy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	ceevent "github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/kyma-project/eventing-publisher-proxy/internal"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/application"
	eppapi "github.com/kyma-project/eventing-publisher-proxy/pkg/legacy/api"
)

var (
	validEventTypeVersion = regexp.MustCompile(AllowedEventTypeVersionChars)
	validEventID          = regexp.MustCompile(AllowedEventIDChars)
)

const (
	requestBodyTooLargeErrorMessage = "http: request body too large"
	eventTypeVersionExtensionKey    = "eventtypeversion"
)

type RequestToCETransformer interface {
	ExtractPublishRequestData(*http.Request) (*eppapi.PublishRequestData, *eppapi.PublishEventResponses, error)
	TransformPublishRequestToCloudEvent(*eppapi.PublishRequestData) (*ceevent.Event, error)
	WriteLegacyRequestsToCE(http.ResponseWriter, *eppapi.PublishRequestData) (*ceevent.Event, string)
	WriteCEResponseAsLegacyResponse(http.ResponseWriter, int, *ceevent.Event, string)
}

type Transformer struct {
	eventMeshNamespace string
	eventTypePrefix    string
	applicationLister  *application.Lister // applicationLister will be nil when disabled.
}

func NewTransformer(bebNamespace string, eventTypePrefix string, applicationLister *application.Lister) *Transformer {
	return &Transformer{
		eventMeshNamespace: bebNamespace,
		eventTypePrefix:    eventTypePrefix,
		applicationLister:  applicationLister,
	}
}

func (t *Transformer) isApplicationListerEnabled() bool {
	return t.applicationLister != nil
}

// CheckParameters validates the parameters in the request and sends error responses if found invalid.
func (t *Transformer) checkParameters(parameters *eppapi.PublishEventParametersV1) *eppapi.PublishEventResponses {
	if parameters == nil {
		return ErrorResponseBadRequest(ErrorMessageBadPayload)
	}
	if len(parameters.PublishrequestV1.EventType) == 0 {
		return ErrorResponseMissingFieldEventType()
	}
	if len(parameters.PublishrequestV1.EventTypeVersion) == 0 {
		return ErrorResponseMissingFieldEventTypeVersion()
	}
	if !validEventTypeVersion.MatchString(parameters.PublishrequestV1.EventTypeVersion) {
		return ErrorResponseWrongEventTypeVersion()
	}
	if len(parameters.PublishrequestV1.EventTime) == 0 {
		return ErrorResponseMissingFieldEventTime()
	}
	if _, err := time.Parse(time.RFC3339, parameters.PublishrequestV1.EventTime); err != nil {
		return ErrorResponseWrongEventTime()
	}
	if len(parameters.PublishrequestV1.EventID) > 0 && !validEventID.MatchString(parameters.PublishrequestV1.EventID) {
		return ErrorResponseWrongEventID()
	}
	if parameters.PublishrequestV1.Data == nil {
		return ErrorResponseMissingFieldData()
	}
	if d, ok := (parameters.PublishrequestV1.Data).(string); ok && len(d) == 0 {
		return ErrorResponseMissingFieldData()
	}
	// OK
	return &eppapi.PublishEventResponses{}
}

// ExtractPublishRequestData extracts the data for publishing event from the given legacy event request.
func (t *Transformer) ExtractPublishRequestData(request *http.Request) (*eppapi.PublishRequestData,
	*eppapi.PublishEventResponses, error) {
	// parse request body to PublishRequestV1
	if request.Body == nil || request.ContentLength == 0 {
		resp := ErrorResponseBadRequest(ErrorMessageBadPayload)
		return nil, resp, errors.New(resp.Error.Message)
	}

	parameters := &eppapi.PublishEventParametersV1{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&parameters.PublishrequestV1); err != nil {
		var resp *eppapi.PublishEventResponses
		if err.Error() == requestBodyTooLargeErrorMessage {
			resp = ErrorResponseRequestBodyTooLarge(err.Error())
		} else {
			resp = ErrorResponseBadRequest(err.Error())
		}
		return nil, resp, errors.New(resp.Error.Message)
	}

	// validate the PublishRequestV1 for missing / incoherent values
	checkResp := t.checkParameters(parameters)
	if checkResp.Error != nil {
		return nil, checkResp, errors.New(checkResp.Error.Message)
	}

	appName := ParseApplicationNameFromPath(request.URL.Path)
	publishRequestData := &eppapi.PublishRequestData{
		PublishEventParameters: parameters,
		ApplicationName:        appName,
		URLPath:                request.URL.Path,
		Headers:                request.Header,
	}

	return publishRequestData, nil, nil
}

// WriteLegacyRequestsToCE transforms the legacy event to cloudevent from the given request.
// It also returns the original event-type without cleanup as the second return type.
func (t *Transformer) WriteLegacyRequestsToCE(writer http.ResponseWriter,
	publishData *eppapi.PublishRequestData) (*ceevent.Event, string) {
	uncleanedAppName := publishData.ApplicationName

	// clean the application name form non-alphanumeric characters
	// handle non-existing applications
	appName := application.GetCleanName(uncleanedAppName)
	// check if we need to use name from application CR.
	if t.isApplicationListerEnabled() {
		if appObj, err := t.applicationLister.Get(uncleanedAppName); err == nil {
			// handle existing applications
			appName = application.GetCleanTypeOrName(appObj)
		}
	}

	event, err := t.convertPublishRequestToCloudEvent(appName, publishData.PublishEventParameters)
	if err != nil {
		response := ErrorResponse(http.StatusInternalServerError, err)
		WriteJSONResponse(writer, response)
		return nil, ""
	}

	// prepare the original event-type without cleanup
	eventType := formatEventType(t.eventTypePrefix, publishData.ApplicationName,
		publishData.PublishEventParameters.PublishrequestV1.EventType,
		publishData.PublishEventParameters.PublishrequestV1.EventTypeVersion)

	return event, eventType
}

func (t *Transformer) WriteCEResponseAsLegacyResponse(writer http.ResponseWriter, statusCode int,
	event *ceevent.Event, msg string) {
	response := &eppapi.PublishEventResponses{}
	// Fail
	if !is2XXStatusCode(statusCode) {
		response.Error = &eppapi.Error{
			Status:  statusCode,
			Message: msg,
		}
		WriteJSONResponse(writer, response)
		return
	}

	// Success
	response.Ok = &eppapi.PublishResponse{EventID: event.ID()}
	WriteJSONResponse(writer, response)
}

// TransformPublishRequestToCloudEvent converts the given publish request to a CloudEvent with raw values.
func (t *Transformer) TransformPublishRequestToCloudEvent(publishRequestData *eppapi.PublishRequestData) (*ceevent.Event,
	error) {
	source := publishRequestData.ApplicationName
	publishRequest := publishRequestData.PublishEventParameters

	// instantiate a new cloudEvent object
	event := ceevent.New(ceevent.CloudEventsVersionV1)
	eventName := publishRequest.PublishrequestV1.EventType
	eventTypeVersion := publishRequest.PublishrequestV1.EventTypeVersion

	// set type by combining type and version (<type>.<version>) e.g. order.created.v1
	event.SetType(fmt.Sprintf("%s.%s", eventName, eventTypeVersion))
	event.SetSource(source)
	event.SetExtension(eventTypeVersionExtensionKey, eventTypeVersion)
	event.SetDataContentType(internal.ContentTypeApplicationJSON)

	// set cloudEvent time
	evTime, err := time.Parse(time.RFC3339, publishRequest.PublishrequestV1.EventTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse time from the external publish request")
	}
	event.SetTime(evTime)

	// set cloudEvent data
	if err := event.SetData(internal.ContentTypeApplicationJSON, publishRequest.PublishrequestV1.Data); err != nil {
		return nil, errors.Wrap(err, "failed to set data to CloudEvent data field")
	}

	// set the event id from the request if it is available
	// otherwise generate a new one
	if len(publishRequest.PublishrequestV1.EventID) > 0 {
		event.SetID(publishRequest.PublishrequestV1.EventID)
	} else {
		event.SetID(uuid.New().String())
	}

	return &event, nil
}

// convertPublishRequestToCloudEvent converts the given publish request to a CloudEvent.
func (t *Transformer) convertPublishRequestToCloudEvent(appName string,
	publishRequest *eppapi.PublishEventParametersV1) (*ceevent.Event, error) {
	if !application.IsCleanName(appName) {
		return nil, errors.New("application name should be cleaned from none-alphanumeric characters")
	}

	event := ceevent.New(ceevent.CloudEventsVersionV1)

	evTime, err := time.Parse(time.RFC3339, publishRequest.PublishrequestV1.EventTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse time from the external publish request")
	}
	event.SetTime(evTime)

	if err := event.SetData(internal.ContentTypeApplicationJSON, publishRequest.PublishrequestV1.Data); err != nil {
		return nil, errors.Wrap(err, "failed to set data to CloudEvent data field")
	}

	// set the event id from the request if it is available
	// otherwise generate a new one
	if len(publishRequest.PublishrequestV1.EventID) > 0 {
		event.SetID(publishRequest.PublishrequestV1.EventID)
	} else {
		event.SetID(uuid.New().String())
	}

	eventName := combineEventNameSegments(removeNonAlphanumeric(publishRequest.PublishrequestV1.EventType))
	prefix := removeNonAlphanumeric(t.eventTypePrefix)
	eventType := formatEventType(prefix, appName, eventName, publishRequest.PublishrequestV1.EventTypeVersion)
	event.SetType(eventType)
	event.SetSource(t.eventMeshNamespace)
	event.SetExtension(eventTypeVersionExtensionKey, publishRequest.PublishrequestV1.EventTypeVersion)
	event.SetDataContentType(internal.ContentTypeApplicationJSON)
	return &event, nil
}

// combineEventNameSegments returns an eventName with exactly two segments separated by "." if the given event-type
// has two or more segments separated by "." (e.g. "Account.Order.Created" becomes "AccountOrder.Created").
func combineEventNameSegments(eventName string) string {
	parts := strings.Split(eventName, ".")
	if len(parts) > 1 {
		businessObject := strings.Join(parts[0:len(parts)-1], "")
		operation := parts[len(parts)-1]
		eventName = fmt.Sprintf("%s.%s", businessObject, operation)
	}
	return eventName
}

// removeNonAlphanumeric returns an eventName without any non-alphanumerical character besides dot (".").
func removeNonAlphanumeric(eventType string) string {
	return regexp.MustCompile("[^a-zA-Z0-9.]+").ReplaceAllString(eventType, "")
}
