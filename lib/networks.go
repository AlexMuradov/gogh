package lib

import (
	"fmt"
	"log"
)

type Networks struct {
	Net []Network
}

type Network struct {
	Addr string
	Name string
	Cidr string
}

func (n *Networks) Add(network Network) {
	n.Net = append(n.Net, network)
}

// Populate uses the common logic for populating networks
func (n *Networks) Populate(plan *TerraformPlan, layersData map[string]interface{}) error {
	// Define a resource processor specifically for networks
	networkProcessor := func(resource map[string]interface{}) bool {
		resourceType, _ := resource["type"].(string)
		// Check if the resource type is what we're interested in
		resourceName, _ := resource["name"].(string) // Assuming 'name' is always present
		resourceCIDR, _ := resource["values"].(map[string]interface{})["cidr_block"].(string)
		n.Add(Network{
			Addr: resourceType,
			Name: resourceName,
			Cidr: resourceCIDR, // Example CIDR, adjust as necessary
		})
		return true
	}

	// Use the common PopulateCommon function with the networkProcessor
	err := PopulateCommon(plan, layersData, "network", networkProcessor)
	if err != nil {
		return fmt.Errorf("failed to populate networks: %w", err)
	}

	return nil
}

func NewNetworkLayer(plan *TerraformPlan, layersData map[string]interface{}) error {
	n := &Networks{}
	err := n.Populate(plan, layersData)
	if err != nil {
		log.Fatalf("Network layer can't be populated: %v", err)
	}

	fmt.Println(n)
	return nil
}
