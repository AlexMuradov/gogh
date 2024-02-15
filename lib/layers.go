package lib

import (
	"encoding/json"
	"os"
)

// !!! TODO: Use interface{} to avoid repeating Populate() Add() function for each layer

// This function is used to initialize layers map by reading Layers JSON and storing it

func InitLayers(layersData string) (map[string]interface{}, error) {

	file, err := os.Open(layersData)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var Layers map[string]interface{}

	err = json.NewDecoder(file).Decode(&Layers)
	if err != nil {
		return nil, err
	}
	return Layers, nil
}
