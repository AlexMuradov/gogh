package lib

import (
	"fmt"
	"log"
)

type Subnetworks struct {
	Subnet []Subnetwork
}

type Subnetwork struct {
	Addr string
	Name string
	Cidr string
}

func (s *Subnetworks) Add(subnetwork Subnetwork) error {
	s.Subnet = append(s.Subnet, subnetwork)
	return nil
}

func (s *Subnetworks) Populate(plan *TerraformPlan, layersData map[string]interface{}) error {
	// Assuming layersData["subnetwork"] contains the types of resources we're interested in
	subnetworkLayerData, ok := layersData["subnetwork"].([]interface{})
	if !ok {
		return fmt.Errorf("expected 'subnetwork' key to be of type []interface{}, got %T", layersData["subnetwork"])
	}

	seen := make(map[string]bool)
	for _, val := range subnetworkLayerData {
		strVal, ok := val.(string)
		if !ok {
			log.Printf("Warning: Non-string value found in subnetwork layer data, ignoring: %v", val)
			continue
		}
		seen[strVal] = true
	}

	processResources := func(resources []map[string]interface{}) {
		for _, resource := range resources {
			resourceType, ok := resource["type"].(string)
			if !ok {
				log.Printf("Warning: Resource type is not a string: %v", resource["type"])
				continue
			}
			// Check if the resource type is "aws_subnet" and in the seen map
			if resourceType == "aws_subnet" && seen[resourceType] {
				resourceName, _ := resource["name"].(string) // Assuming name is always present and a string
				log.Printf("AWS subnet match found: %s", resourceName)
				s.Add(Subnetwork{
					Addr: resourceType,
					Name: resourceName,
					Cidr: "0.0.0.0/0", // Example CIDR, adjust as necessary
				})
			}
		}
	}

	// Process child module resources
	for _, childModule := range plan.PlannedValues.RootModule.ChildModules {
		processResources(childModule.Resources)
	}

	return nil
}

func NewSubnetworkLayer(plan *TerraformPlan, layersData map[string]interface{}) error {
	s := &Subnetworks{}
	if err := s.Populate(plan, layersData); err != nil {
		return fmt.Errorf("Subnet layer can't be populated: %w", err)
	}

	fmt.Println(s)
	return nil
}
