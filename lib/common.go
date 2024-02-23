package lib

import (
	"fmt"
	"log"
)

// ResourceProcessor defines a function type for processing resources
type ResourceProcessor func(resource map[string]interface{}) bool

// ProcessResources generic resource processing function
func ProcessResources(resources []map[string]interface{}, processor ResourceProcessor) {
	for _, resource := range resources {
		if processor(resource) {
			// If the processor returns true, log the match
			resourceType, _ := resource["type"].(string)
			log.Printf("Resource match found: %s", resourceType)
		}
	}
}

// PopulateCommon generic population logic for both networks and subnetworks
func PopulateCommon(plan *TerraformPlan, layersData map[string]interface{}, layerKey string, processor ResourceProcessor) error {
	layerData, ok := layersData[layerKey].([]interface{})
	if !ok {
		return fmt.Errorf("expected '%s' key to be of type []interface{}, got %T", layerKey, layersData[layerKey])
	}

	seen := make(map[string]bool)
	for _, val := range layerData {
		strVal, ok := val.(string)
		if !ok {
			log.Printf("Warning: Non-string value found in %s layer data, ignoring: %v", layerKey, val)
			continue
		}
		seen[strVal] = true
	}

	// Process root module resources
	ProcessResources(plan.PlannedValues.RootModule.Resources, func(resource map[string]interface{}) bool {
		resourceType, ok := resource["type"].(string)
		if !ok {
			log.Printf("Warning: Resource type is not a string: %v", resource["type"])
			return false
		}
		return seen[resourceType] && processor(resource)
	})

	// Process child module resources
	for _, childModule := range plan.PlannedValues.RootModule.ChildModules {
		ProcessResources(childModule.Resources, func(resource map[string]interface{}) bool {
			resourceType, ok := resource["type"].(string)
			if !ok {
				log.Printf("Warning: Resource type is not a string: %v", resource["type"])
				return false
			}
			return seen[resourceType] && processor(resource)
		})
	}

	return nil
}
