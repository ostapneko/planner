package gantt

import (
	"fmt"
	"github.com/ostapneko/planner"
	"io"
	"log"
	"sort"
	"time"
)

type Color string

var colors = []Color{
	"LightCoral",
	"LightGreen",
	"LightBlue",
	"LightPink",
	"Gold",
	"LightGray",
}

type drawer struct {
	devToColor map[planner.DeveloperId]Color
}

func newDrawer(devToColor map[planner.DeveloperId]Color) *drawer {
	return &drawer{devToColor: devToColor}
}

func (g *drawer) drawLine(firstDay planner.Day, lastDay planner.Day, name string, developerId planner.DeveloperId) string {
	firstDayDate := dayToPlantUMLDate(firstDay)
	lastDayDate := dayToPlantUMLDate(lastDay)
	color := g.devToColor[developerId]
	line := fmt.Sprintf("[<font:sans>%s (%s)] is colored in %s and starts on %s and ends on %s\n", name, developerId, color, firstDayDate, lastDayDate)
	return line
}

func (g *drawer) drawMilestone(day planner.Day, name string) string {
	return fmt.Sprintf("[<font:sans>%s] happens on %s\n", name, dayToPlantUMLDate(day))
}

type writer struct {
	planning *planner.Planning
	w io.Writer
	drawer *drawer
}

func (writer *writer) writeStr(str string) {
	_, _ = writer.w.Write([]byte(str))
}

func newWriter(planning *planner.Planning, w io.Writer, drawer *drawer) *writer {
	return &writer{planning: planning, w: w, drawer: drawer}
}

func ToPlantUML(planning *planner.Planning, w io.Writer) {
	devToColor := map[planner.DeveloperId]Color{}
	for i, developer := range planning.Developers {
		color := colors[i % len(colors)]
		devToColor[developer.Id] = color
	}

	drawer := newDrawer(devToColor)
	writer := newWriter(planning, w, drawer)

	writer.header()
	writer.closedDays()
	writer.projectStart()
	writer.tasks()
	writer.offdays()
	writer.supportWeeks()
	writer.end()
}

func (writer *writer) section(str string) {
	writer.writeStr(fmt.Sprintf("-- %s --\n", str))
}

func (writer *writer) header() {
	writer.writeStr("@startgantt\nprintscale daily\n")
}

func (writer *writer) closedDays() {
	writer.writeStr("saturday are closed\nsunday are closed\n")

	for _, holiday := range writer.planning.Holidays {
		writer.writeStr(fmt.Sprintf("%s is closed\n", dayToPlantUMLDate(holiday)))
	}
}

func (writer *writer) projectStart() {
	startDay := writer.planning.StartDay
	writer.writeStr(fmt.Sprintf("Project starts on %s\n", dayToPlantUMLDate(startDay)))
}

func (writer *writer) tasks() {
	writer.section("Roadmap")
	for _, task := range writer.planning.Tasks {
		for developerId, attribution := range task.Attributions {
			firstDay := attribution.FirstDay
			lastDay := attribution.LastDay
			if firstDay == nil || lastDay == nil {
				continue
			}

			line := writer.drawer.drawLine(*firstDay, *lastDay, task.Name, developerId)
			writer.writeStr(line)
		}
		milestone := writer.drawer.drawMilestone(*task.LastDay, fmt.Sprintf("%s completed", task.Name))
		writer.writeStr(milestone)
	}
}

func (writer *writer) supportWeeks() {
	writer.section("Support Weeks")
	for i, week := range writer.planning.SupportWeeks {
		name := fmt.Sprintf("Support Week %d", i)
		line := writer.drawer.drawLine(week.FirstDay, week.LastDay, name, week.DevId)
		writer.writeStr(line)
	}
}

func (writer *writer) end() {
	writer.writeStr("@endgantt")
}

func (writer *writer) offdays() {
	writer.section("Vacations")
	for _, developer := range writer.planning.Developers {
		// start
		if developer.Starts != nil {
			ms := writer.drawer.drawMilestone(*developer.Starts, fmt.Sprintf("%s starts", developer.Id))
			writer.writeStr(ms)
		}

		// vacations
		// find contiguous days and make a line out of them

		days := developer.OffDays
		sort.Sort(days)
		var firstDay *planner.Day
		var lastDay *planner.Day

		i := 0
		for _, day := range days {
			d := day
			if firstDay == nil {
				firstDay = &d
				lastDay = &d
				continue
			}

			if lastDay == nil {
				log.Fatalf("Unreachable code")
			}

			if int(d) == int(*lastDay) + 1 {
				lastDay = &d
				continue
			}

			i++
			line := writer.drawer.drawLine(*firstDay, *lastDay, fmt.Sprintf("%s - %d", developer.Id, i), developer.Id)
			writer.writeStr(line)
			firstDay = &d
			lastDay = &d
		}

		if firstDay != nil && lastDay != nil {
			i++
			line := writer.drawer.drawLine(*firstDay, *lastDay, fmt.Sprintf("%s - %d", developer.Id, i), developer.Id)
			writer.writeStr(line)
		}

		// end
		if developer.Leaves != nil {
			ms := writer.drawer.drawMilestone(*developer.Leaves, fmt.Sprintf("%s leaves", developer.Id))
			writer.writeStr(ms)
		}
	}
}

// from  18307 (nb of days since epoch) -> 15/02/2020
func dayToPlantUMLDate(day planner.Day) string {
	epoch := time.Unix(0, 0)
	t := epoch.Add(time.Duration(int(day)) * time.Hour * 24)
	return t.Format("2006/01/02")
}
