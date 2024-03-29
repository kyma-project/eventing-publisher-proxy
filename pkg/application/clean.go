package application

import (
	"regexp"
	"strings"

	kymaappconnv1alpha1 "github.com/kyma-project/kyma/components/central-application-gateway/pkg/apis/applicationconnector/v1alpha1"
)

const (
	// TypeLabel is an optional label for the application to determine its type.
	TypeLabel = "application-type"
)

// invalidApplicationNameSegment used to match and replace none-alphanumeric characters in the application name.
var invalidApplicationNameSegment = regexp.MustCompile(`\W|_`)

// GetCleanTypeOrName cleans the application name form none-alphanumeric characters and returns it
// if the application type label exists, it will be cleaned and returned instead of the application name.
func GetCleanTypeOrName(application *kymaappconnv1alpha1.Application) string {
	if application == nil {
		return ""
	}
	applicationName := application.Name
	for k, v := range application.Labels {
		if strings.ToLower(k) == TypeLabel {
			applicationName = v
			break
		}
	}
	return GetCleanName(applicationName)
}

// GetTypeOrName returns the application name.
// if the application type label exists, it will be returned instead of the application name.
func GetTypeOrName(application *kymaappconnv1alpha1.Application) string {
	if application == nil {
		return ""
	}
	applicationName := application.Name
	for k, v := range application.Labels {
		if strings.ToLower(k) == TypeLabel {
			applicationName = v
			break
		}
	}
	return applicationName
}

// GetCleanName cleans the name form none-alphanumeric characters and returns the clean name.
func GetCleanName(name string) string {
	return invalidApplicationNameSegment.ReplaceAllString(name, "")
}

// IsCleanName returns true if the name contains alphanumeric characters only, otherwise returns false.
func IsCleanName(name string) bool {
	return !invalidApplicationNameSegment.MatchString(name)
}
