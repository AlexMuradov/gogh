package main

import (
	"fmt"
	"log"
	"time"

	//"github.com/ajstarks/svgo"
	"github.com/architecthub-io/gogh/lib"
)

func main() {

	asciiArt := `
         _              _            _              _       _  
        /\ \           /\ \         /\ \           / /\    / /\
       /  \ \         /  \ \       /  \ \         / / /   / / /
      / /\ \_\       / /\ \ \     / /\ \_\       / /_/   / / / 
     / / /\/_/      / / /\ \ \   / / /\/_/      / /\ \__/ / /  
    / / / ______   / / /  \ \_\ / / / ______   / /\ \___\/ /   
   / / / /\_____\ / / /   / / // / / /\_____\ / / /\/___/ /    
  / / /  \/____ // / /   / / // / /  \/____ // / /   / / /     
 / / /_____/ / // / /___/ / // / /_____/ / // / /   / / /      
/ / /______\/ // / /____\/ // / /______\/ // / /   / / /       
\/___________/ \/_________/ \/___________/ \/_/    \/_/        
`
	fmt.Println(asciiArt)
	fmt.Printf("working")

	done := make(chan bool)

	// Start the goroutine that prints dots
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// Print a dot and sleep for a specified interval
				fmt.Printf(".")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Loads configuration
	config, err := lib.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Converts terraform into JSON
	if err := lib.TF2Json(config.TerraformFile, config.PlanFile); err != nil {
		log.Fatalf("Failed to convert Terraform to JSON: %v", err)
	}

	// Exctracts resources from terraform JSON and stores it
	plan, err := lib.NewTerraformPlan(config.PlanFile)
	if err != nil {
		log.Fatalf("Failed to load Terraform plan: %v", err)
	}

	// Initializes layers map by reading Layers JSON and storing it
	// This map is used for pairing infrastructure code with resources
	// that needs to be create on the layer level
	layers, err := lib.InitLayers(config.Layers)

	// Initializes network layer by using network key and slice values
	// within that key. It then pairs network items in the that slice
	// by going through all items in the infrastructure code using hashmaps

	layerInitializers := []func(*lib.TerraformPlan, map[string]interface{}) error{ // Replace PlanType and LayersType with actual types
		lib.NewNetworkLayer,
		lib.NewSubnetworkLayer,
	}

	for _, initializer := range layerInitializers {
		err := initializer(plan, layers)
		if err != nil {
			log.Fatalf("error %v", err)
		}
	}

	// 	http.Handle("/circle", http.HandlerFunc(circle))
	// 	err = http.ListenAndServe(":2003", nil)
	// 	if err != nil {
	// 		log.Fatal("ListenAndServe:", err)
	// 	}
	// }

	// func circle(w http.ResponseWriter, req *http.Request) {
	// 	w.Header().Set("Content-Type", "image/svg+xml")
	// 	s := svg.New(w)
	// 	s.Start(500, 500)
	// 	s.Circle(250, 250, 125, "fill:none;stroke:black")
	// 	ant := "https://icon.icepanel.io/GCP/svg/Anthos.svg"

	// pathData := "M109"
	// s.Path(pathData, "fill:none;stroke:black;stroke-width:2")
	// s.Text(100, 100, "Hello, SVG!", "text-anchor:middle;fill:black;font-size:24px;font-family:montserrat")
	// s.Line(10, 10, 490, 490, "stroke:black;stroke-width:2;stroke-dasharray:10,10")
	// s.Rect(10, 200, 60, 60, "fill:none;stroke:#E4CCFF;stroke-width:1;rx:12; fill:#FBF7FF")
	// s.Image(20, 205, 40, 40, ant)
	// s.End()
	done <- true
}
