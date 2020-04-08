package gantt

import (
	"fmt"
	"github.com/ostapneko/planner"
	"strings"
	"time"
)

func ToPlantUML(planning *planner.Planning) string {
	var b strings.Builder

	b.WriteString("@startgantt\nprintscale daily\n")

	firstDay, _ := planner.DateToDay("01/01/3000")

	for _, task := range planning.Tasks {
		for _, attribution := range task.Attributions {
			if attribution.FirstDay != nil && *attribution.FirstDay < firstDay {
				firstDay = *attribution.FirstDay
			}
		}
	}

	b.WriteString(fmt.Sprintf("Project starts on %s\n", dayToPlantUMLDate(firstDay)))

	for _, task := range planning.Tasks {
		for developerId, attribution := range task.Attributions {
			if attribution.FirstDay == nil || attribution.LastDay == nil {
				continue
			}

			firstDayDate := dayToPlantUMLDate(*attribution.FirstDay)
			lastDayDate := dayToPlantUMLDate(*attribution.LastDay)
			b.WriteString(fmt.Sprintf("[%s (%s)] starts on %s and ends on %s\n", task.Name, developerId, firstDayDate, lastDayDate))
		}
	}

	b.WriteString("@endgantt")
	return b.String()
}

// from  18307 (nb of days since epoch) -> 15/02/2020
func dayToPlantUMLDate(day planner.Day) string {
	epoch := time.Unix(0, 0)
	t := epoch.Add(time.Duration(int(day)) * time.Hour * 24)
	return t.Format("2006/01/02")
}
