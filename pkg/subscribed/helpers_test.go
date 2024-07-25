package subscribed

import (
	"reflect"
	"testing"

	emeventingv2alpha1 "github.com/kyma-project/eventing-manager/api/eventing/v1alpha2"
)

func TestFilterEventTypeVersions(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		appName        string
		subscription   *emeventingv2alpha1.Subscription
		expectedEvents []Event
	}{
		{
			name:           "should return no events when there is no subscription",
			appName:        "fooapp",
			subscription:   &emeventingv2alpha1.Subscription{},
			expectedEvents: make([]Event, 0),
		}, {
			name:    "should return a slice of events when eventTypes are provided",
			appName: "foovarkes",
			subscription: &emeventingv2alpha1.Subscription{
				Spec: emeventingv2alpha1.SubscriptionSpec{
					Source: "foovarkes",
					Types: []string{
						"order.created.v1",
						"order.created.v2",
					},
				},
			},
			expectedEvents: []Event{
				NewEvent("order.created", "v1"),
				NewEvent("order.created", "v2"),
			},
		}, {
			name:    "should return no event if app name is different than subscription source",
			appName: "foovarkes",
			subscription: &emeventingv2alpha1.Subscription{
				Spec: emeventingv2alpha1.SubscriptionSpec{
					Source: "diff-source",
					Types: []string{
						"order.created.v1",
						"order.created.v2",
					},
				},
			},
			expectedEvents: []Event{},
		}, {
			name:    "should return event types if event type consists of eventType and appName for typeMaching exact",
			appName: "foovarkes",
			subscription: &emeventingv2alpha1.Subscription{
				Spec: emeventingv2alpha1.SubscriptionSpec{
					Source:       "/default/sap.kyma/tunas-develop",
					TypeMatching: emeventingv2alpha1.TypeMatchingExact,
					Types: []string{
						"sap.kyma.custom.foovarkes.order.created.v1",
						"sap.kyma.custom.foovarkes.order.created.v2",
					},
				},
			},
			expectedEvents: []Event{
				NewEvent("order.created", "v1"),
				NewEvent("order.created", "v2"),
			},
		}, {
			name:    "should return no event if app name is not part of external event types",
			appName: "foovarkes",
			subscription: &emeventingv2alpha1.Subscription{
				Spec: emeventingv2alpha1.SubscriptionSpec{
					Source:       "/default/sap.kyma/tunas-develop",
					TypeMatching: emeventingv2alpha1.TypeMatchingExact,
					Types: []string{
						"sap.kyma.custom.difffoovarkes.order.created.v1",
						"sap.kyma.custom.difffoovarkes.order.created.v2",
					},
				},
			},
			expectedEvents: []Event{},
		}, {
			name:    "should return event type only with 'sap.kyma.custom' prefix and appname",
			appName: "foovarkes",
			subscription: &emeventingv2alpha1.Subscription{
				Spec: emeventingv2alpha1.SubscriptionSpec{
					Source:       "/default/sap.kyma/tunas-develop",
					TypeMatching: emeventingv2alpha1.TypeMatchingExact,
					Types: []string{
						"foo.prefix.custom.foovarkes.order.created.v1",
						"sap.kyma.custom.foovarkes.order.created.v2",
						"sap.kyma.custom.diffvarkes.order.created.v2",
					},
				},
			},
			expectedEvents: []Event{
				NewEvent("order.created", "v2"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotEvents := FilterEventTypeVersions("sap.kyma.custom", tc.appName, tc.subscription)
			if !reflect.DeepEqual(tc.expectedEvents, gotEvents) {
				t.Errorf("Received incorrect events, Wanted: %v, Got: %v", tc.expectedEvents, gotEvents)
			}
		})
	}
}

func TestBuildEventType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		appName string
		want    Event
	}{
		{
			name:    "should return no events when there is no subscription",
			appName: "order.created.v1",
			want: Event{
				Name:    "order.created",
				Version: "v1",
			},
		}, {
			name:    "should return a slice of events when eventTypes are provided",
			appName: "product.order.created.v1",
			want: Event{
				Name:    "product.order.created",
				Version: "v1",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			event := buildEvent(tc.appName)
			if !reflect.DeepEqual(tc.want, event) {
				t.Errorf("Received incorrect events, Wanted: %v, Got: %v", tc.want, event)
			}
		})
	}
}

func TestConvertEventsMapToSlice(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		inputMap     map[Event]bool
		wantedEvents []Event
	}{
		{
			name: "should return events from the map in a slice",
			inputMap: map[Event]bool{
				NewEvent("foo", "v1"): true,
				NewEvent("bar", "v2"): true,
			},
			wantedEvents: []Event{
				NewEvent("foo", "v1"),
				NewEvent("bar", "v2"),
			},
		}, {
			name:         "should return no events for an empty map of events",
			inputMap:     map[Event]bool{},
			wantedEvents: []Event{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotEvents := ConvertEventsMapToSlice(tc.inputMap)
			for _, event := range gotEvents {
				found := false
				for _, wantEvent := range tc.wantedEvents {
					if event == wantEvent {
						found = true
						continue
					}
				}
				if !found {
					t.Errorf("incorrect slice of events, wanted: %v, got: %v", tc.wantedEvents, gotEvents)
				}
			}
		})
	}
}

func TestAddUniqueEventsToResult(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name                   string
		eventsSubSet           []Event
		givenUniqEventsAlready map[Event]bool
		wantedUniqEvents       map[Event]bool
	}{
		{
			name: "should return unique events along with the existing ones",
			eventsSubSet: []Event{
				NewEvent("foo", "v1"),
				NewEvent("bar", "v1"),
			},
			givenUniqEventsAlready: map[Event]bool{
				NewEvent("bar-already-existing", "v1"): true,
			},
			wantedUniqEvents: map[Event]bool{
				NewEvent("foo", "v1"):                  true,
				NewEvent("bar", "v1"):                  true,
				NewEvent("bar-already-existing", "v1"): true,
			},
		}, {
			name: "should return unique new events from the subset provided only",
			eventsSubSet: []Event{
				NewEvent("foo", "v1"),
				NewEvent("bar", "v1"),
			},
			givenUniqEventsAlready: nil,
			wantedUniqEvents: map[Event]bool{
				NewEvent("foo", "v1"): true,
				NewEvent("bar", "v1"): true,
			},
		}, {
			name:         "should return existing unique events when an empty subset provided",
			eventsSubSet: []Event{},
			givenUniqEventsAlready: map[Event]bool{
				NewEvent("foo", "v1"): true,
				NewEvent("bar", "v1"): true,
			},
			wantedUniqEvents: map[Event]bool{
				NewEvent("foo", "v1"): true,
				NewEvent("bar", "v1"): true,
			},
		}, {
			name:                   "should return no unique events when an empty subset provided",
			eventsSubSet:           []Event{},
			givenUniqEventsAlready: map[Event]bool{},
			wantedUniqEvents:       map[Event]bool{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotUniqEvents := AddUniqueEventsToResult(tc.eventsSubSet, tc.givenUniqEventsAlready)
			if !reflect.DeepEqual(tc.wantedUniqEvents, gotUniqEvents) {
				t.Errorf("incorrect unique events, wanted: %v, got: %v", tc.wantedUniqEvents, gotUniqEvents)
			}
		})
	}
}

func NewEvent(name, version string) Event {
	return Event{
		Name:    name,
		Version: version,
	}
}
