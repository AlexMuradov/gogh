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
	Zone string
}

func (n *Subnetworks) Add(subnetwork Subnetwork) {
	n.Subnet = append(n.Subnet, subnetwork)
}

// Populate uses the common logic for populating networks
func (n *Subnetworks) Populate(plan *TerraformPlan, layersData map[string]interface{}) error {
	// Define a resource processor specifically for networks
	Processor := func(resource map[string]interface{}) bool {
		resourceType, _ := resource["type"].(string)
		// Check if the resource type is what we're interested in
		resourceName, _ := resource["name"].(string) // Assuming 'name' is always present
		resourceCIDR, _ := resource["values"].(map[string]interface{})["cidr_block"].(string)
		resourceZone, _ := resource["values"].(map[string]interface{})["availability_zone"].(string)
		n.Add(Subnetwork{
			Addr: resourceType,
			Name: resourceName,
			Cidr: resourceCIDR,
			Zone: resourceZone,
		})
		return true
	}

	// Use the common PopulateCommon function with the networkProcessor
	err := PopulateCommon(plan, layersData, "subnetwork", Processor)
	if err != nil {
		return fmt.Errorf("failed to populate subnetworks: %w", err)
	}

	return nil
}

func NewSubnetworkLayer(plan *TerraformPlan, layersData map[string]interface{}) error {
	n := &Subnetworks{}
	err := n.Populate(plan, layersData)
	if err != nil {
		log.Fatalf("Subetwork layer can't be populated: %v", err)
	}

	fmt.Println(n)
	return nil
}
