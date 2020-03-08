package main

import (
	"adalongcorp.com/planner"
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

	var planning planner.Planning

	err = yaml.Unmarshal(dat, &planning)

	if err != nil {
		log.Fatalf("error parsing planning: %s", err)
	}

	log.Printf("%+v", planning)
}
