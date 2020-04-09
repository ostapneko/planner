package main

import (
	"encoding/json"
	"github.com/ostapneko/planner"
	"github.com/ostapneko/planner/gantt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "out",
				Aliases:   []string{"o"},
				Usage:     "output file with tasks completed",
				TakesFile: true,
				Required:  true,
			},
			&cli.StringFlag{
				Name:      "gantt",
				Aliases:   []string{"g"},
				Usage:     "output gantt chart in PlantUML dialect",
				TakesFile: true,
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "json or yaml",
				Value:   "yaml",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				log.Fatalf("Require the input planning as argument")
			}

			inputFile := c.Args().Get(0)

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

			planner.ForecastCompletion(planning)

			planningOutput := planner.NewPlanningInput(planning)

			var doc []byte

			format := c.String("format")
			if format == "yaml" {
				doc, _ = yaml.Marshal(planningOutput)
			} else if format == "json" {
				doc, _ = json.Marshal(planningOutput)
			} else {
				log.Fatalf("Unsupported format %s", format)
			}

			outFile := c.String("out")
			err = ioutil.WriteFile(outFile, doc, 0644)

			if err != nil {
				log.Fatalf("error writing to file %s", outFile)
			}

			if c.IsSet("gantt") {
				gantFile := c.String("gantt")
				file, err  := os.OpenFile(gantFile, os.O_WRONLY | os.O_CREATE, 0644)
				if err != nil {
					log.Fatalf("could not write Gantt chart to %s", gantFile)
				}
				gantt.ToPlantUML(planning, file)
				_ = file.Close()
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
