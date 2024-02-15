package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func TF2Json(execDir string, planfile string) error {

	cmd := exec.Command("terraform", "init")
	cmd.Dir = execDir
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to initialize Terraform: %v", err)
	}

	cmd = exec.Command("terraform", "plan", "-out=plan")
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to plan: %v", err)
	}

	cmd = exec.Command("terraform", "show", "-json", "plan")
	cmd.Dir = execDir
	jsonOutput, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to convert Terraform plan to JSON: %v", err)
	}

	err = os.WriteFile(planfile, jsonOutput, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON output to file: %v", err)
	}

	return nil
}

type TerraformPlan struct {
	PlannedValues struct {
		RootModule struct {
			Resources []struct {
				Type string `json:"type"`
			} `json:"resources"`
			ChildModules []struct {
				Resources []struct {
					Type string `json:"type"`
				} `json:"resources"`
			} `json:"child_modules"`
		} `json:"root_module"`
	} `json:"planned_values"`
}

func (this *TerraformPlan) LoadFromFile(file string) (*TerraformPlan, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if !json.Valid(content) {
		return nil, err
	}

	err = json.Unmarshal(content, this)
	if err != nil {
		return nil, err
	}

	return this, nil
}

func NewTerraformPlan(file string) (*TerraformPlan, error) {
	plan := &TerraformPlan{}
	plan, err := plan.LoadFromFile(file)

	if err != nil {
		log.Fatalf("Failed to load terraform plan %v", err)
	}

	return plan, err
}
