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

func (this *Subnetworks) Add(subnetwork Subnetwork) error {
	this.Subnet = append(this.Subnet, subnetwork)
	return nil
}

func (n *Subnetworks) Populate(plan *TerraformPlan, layersData map[string]interface{}) error {

	subnetwork := layersData["subnetwork"].([]interface{})

	seen := make(map[string]bool)

	for _, val := range subnetwork {
		seen[val.(string)] = true
	}

	for _, childModule := range plan.PlannedValues.RootModule.ChildModules {
		for _, resource := range childModule.Resources {
			if seen[resource.Type] {
				log.Printf("network match found: %s", resource.Type)
				n.Add(Subnetwork{
					Addr: resource.Type,
					Name: "My test network",
					Cidr: "0.0.0.0/0",
				})
			}
		}
	}

	return nil
}

func NewSubnetworkLayer(plan *TerraformPlan, layersData map[string]interface{}) error {

	n := &Subnetworks{}
	err := n.Populate(plan, layersData)

	if err != nil {
		log.Fatalf("Subenet layer can't be populated %v", err)
	}

	fmt.Println(n)

	return nil
}
