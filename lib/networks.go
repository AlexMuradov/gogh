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

func (this *Networks) Add(network Network) error {
	this.Net = append(this.Net, network)
	return nil
}

func (n *Networks) Populate(plan *TerraformPlan, layersData map[string]interface{}) error {

	network := layersData["network"].([]interface{})

	seen := make(map[string]bool)

	for _, val := range network {
		seen[val.(string)] = true
	}

	for _, val := range plan.PlannedValues.RootModule.Resources {
		if seen[val.Type] {
			log.Printf("network match found: %s", val.Type)
			n.Add(Network{
				Addr: val.Type,
				Name: "My test network",
				Cidr: "0.0.0.0/0",
			})
		}
	}

	for _, childModule := range plan.PlannedValues.RootModule.ChildModules {
		for _, resource := range childModule.Resources {
			if seen[resource.Type] {
				log.Printf("network match found: %s", resource.Type)
				n.Add(Network{
					Addr: resource.Type,
					Name: "My test network",
					Cidr: "0.0.0.0/0",
				})
			}
		}
	}

	return nil
}

func NewNetworkLayer(plan *TerraformPlan, layersData map[string]interface{}) error {

	n := &Networks{}
	err := n.Populate(plan, layersData)

	if err != nil {
		log.Fatalf("Network layer can't be populated %v", err)
	}

	fmt.Println(n)

	return nil
}
