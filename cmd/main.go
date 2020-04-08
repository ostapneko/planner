package main

import (
	"fmt"
	"github.com/ostapneko/planner"
	"github.com/ostapneko/planner/gantt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
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

	err = planner.CheckPlanning(planning)

	if err != nil {
		log.Fatalf("inconsistent planning: %s", err)
	}

	now, err := planner.DateToDay(time.Now().Format("02/01/2006"))
	if err != nil {
		log.Fatalf("could not parse date, something is seriously wrong: %s", err)
	}

	planner.ForecastCompletion(planning, now)

	//planningOutput := planner.NewPlanningInput(planning)
	//doc, _ := yaml.Marshal(planningOutput)

	doc := gantt.ToPlantUML(planning)

	fmt.Println(string(doc))
}
