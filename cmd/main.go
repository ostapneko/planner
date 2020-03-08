package main

import (
	"adalongcorp.com/planner"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage %s <INPUT FILE>", os.Args[0])
	}

	inputFile := os.Args[1]

	dat, err := ioutil.ReadFile(inputFile)

	if err != nil {
		log.Fatalf("could not read file %s", inputFile)
	}

	var planningInput planner.PlanningInput

	err = yaml.Unmarshal(dat, &planningInput)

	if err != nil {
		log.Fatalf("error parsing planning: %s", err)
	}

	planning, err := planner.NewPlanning(planningInput)

	if err != nil {
		log.Fatalf("error transforming planning input into planning: %s", err)
	}

	doc, _ := json.Marshal(planning)

	log.Printf("%s", doc)
}
