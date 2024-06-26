package subscribed

import (
	"fmt"
	"strings"

	"github.com/kyma-project/eventing-publisher-proxy/pkg/informers"
	kcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"

	emeventingv1alpha1 "github.com/kyma-project/eventing-manager/api/eventing/v1alpha1"
	emeventingv2alpha1 "github.com/kyma-project/eventing-manager/api/eventing/v1alpha2"
)

func SubscriptionGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Version:  emeventingv2alpha1.GroupVersion.Version,
		Group:    emeventingv2alpha1.GroupVersion.Group,
		Resource: "subscriptions",
	}
}

func SubscriptionV1alpha1GVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Version:  emeventingv1alpha1.GroupVersion.Version,
		Group:    emeventingv1alpha1.GroupVersion.Group,
		Resource: "subscriptions",
	}
}

// ConvertRuntimeObjToSubscriptionV1alpha1 converts a runtime.Object to a v1alpha1 version of Subscription object
// by converting to unstructured in between.
func ConvertRuntimeObjToSubscriptionV1alpha1(sObj runtime.Object) (*emeventingv1alpha1.Subscription, error) {
	sub := &emeventingv1alpha1.Subscription{}
	if subUnstructured, ok := sObj.(*unstructured.Unstructured); ok {
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(subUnstructured.Object, sub)
		if err != nil {
			return nil, err
		}
	}
	return sub, nil
}

// ConvertRuntimeObjToSubscription converts a runtime.Object to a Subscription object
// by converting to unstructured in between.
func ConvertRuntimeObjToSubscription(sObj runtime.Object) (*emeventingv2alpha1.Subscription, error) {
	sub := &emeventingv2alpha1.Subscription{}
	if subUnstructured, ok := sObj.(*unstructured.Unstructured); ok {
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(subUnstructured.Object, sub)
		if err != nil {
			return nil, err
		}
	}
	return sub, nil
}

// GenerateSubscriptionInfFactory generates DynamicSharedInformerFactory for Subscription.
func GenerateSubscriptionInfFactory(k8sConfig *rest.Config) dynamicinformer.DynamicSharedInformerFactory {
	subDynamicClient := dynamic.NewForConfigOrDie(k8sConfig)
	dFilteredSharedInfFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(subDynamicClient,
		informers.DefaultResyncPeriod,
		kcorev1.NamespaceAll,
		nil,
	)
	dFilteredSharedInfFactory.ForResource(SubscriptionGVR())
	return dFilteredSharedInfFactory
}

// ConvertEventsMapToSlice converts a map of Events to a slice of Events.
func ConvertEventsMapToSlice(eventsMap map[Event]bool) []Event {
	result := make([]Event, 0)
	for k := range eventsMap {
		result = append(result, k)
	}
	return result
}

// AddUniqueEventsToResult returns a map of unique Events which also contains the events eventsSubSet.
func AddUniqueEventsToResult(eventsSubSet []Event, uniqEvents map[Event]bool) map[Event]bool {
	if len(uniqEvents) == 0 {
		uniqEvents = make(map[Event]bool)
	}
	for _, event := range eventsSubSet {
		if !uniqEvents[event] {
			uniqEvents[event] = true
		}
	}
	return uniqEvents
}

// FilterEventTypeVersions returns a slice of Events:
// if the event source matches the appName for typeMatching standard
// if the <eventTypePrefix>.<appName> is present in the eventType for typeMatching exact.
func FilterEventTypeVersions(eventTypePrefix, appName string, subscription *emeventingv2alpha1.Subscription) []Event {
	events := make([]Event, 0)
	prefixAndAppName := fmt.Sprintf("%s.%s.", eventTypePrefix, appName)

	for _, eventType := range subscription.Spec.Types {
		if subscription.Spec.TypeMatching == emeventingv2alpha1.TypeMatchingExact {
			// in case of type matching exact, we have app name as a part of event type
			if strings.HasPrefix(eventType, prefixAndAppName) {
				eventTypeVersion := strings.ReplaceAll(eventType, prefixAndAppName, "")
				event := buildEvent(eventTypeVersion)
				events = append(events, event)
			}
		} else {
			// in case of type matching standard, the source must be app name
			if appName == subscription.Spec.Source {
				event := buildEvent(eventType)
				events = append(events, event)
			}
		}
	}
	return events
}

// it receives event and type version, e.g. order.created.v1 and returns `{Name: order.created, Version: v1}`.
func buildEvent(eventTypeAndVersion string) Event {
	lastDotIndex := strings.LastIndex(eventTypeAndVersion, ".")
	eventName := eventTypeAndVersion[:lastDotIndex]
	eventVersion := eventTypeAndVersion[lastDotIndex+1:]
	return Event{
		Name:    eventName,
		Version: eventVersion,
	}
}

// FilterEventTypeVersionsV1alpha1 returns a slice of Events for v1alpha1 version of Subscription resource:
// 1. if the eventType matches the format: <eventTypePrefix><appName>.<event-name>.<version>
// E.g. sap.kyma.custom.varkes.order.created.v0
// 2. if the eventSource matches BEBNamespace name.
func FilterEventTypeVersionsV1alpha1(eventTypePrefix, bebNs, appName string,
	filters *emeventingv1alpha1.BEBFilters,
) []Event {
	events := make([]Event, 0)
	if filters == nil {
		return events
	}
	for _, filter := range filters.Filters {
		if filter == nil {
			continue
		}

		prefix := preparePrefix(eventTypePrefix, appName)

		if filter.EventSource == nil || filter.EventType == nil {
			continue
		}

		// IMPORTANT revisit the filtration logic as part of https://github.com/kyma-project/kyma/issues/10761

		// filter by event-source if exists
		if len(strings.TrimSpace(filter.EventSource.Value)) > 0 && !strings.EqualFold(filter.EventSource.Value, bebNs) {
			continue
		}

		if !strings.HasPrefix(filter.EventType.Value, prefix) {
			continue
		}

		// remove the prefix
		eventTypeVersion := strings.TrimPrefix(filter.EventType.Value, prefix)

		eventTypeVersionArr := strings.Split(eventTypeVersion, ".")

		// join all segments except the last to form the eventType
		eventType := strings.Join(eventTypeVersionArr[:len(eventTypeVersionArr)-1], ".")

		// the last segment is the version
		version := eventTypeVersionArr[len(eventTypeVersionArr)-1]

		event := Event{
			Name:    eventType,
			Version: version,
		}
		events = append(events, event)
	}
	return events
}

func preparePrefix(eventTypePrefix string, appName string) string {
	if len(strings.TrimSpace(eventTypePrefix)) == 0 {
		return strings.ToLower(appName + ".")
	}
	return strings.ToLower(fmt.Sprintf("%s.%s.", eventTypePrefix, appName))
}
