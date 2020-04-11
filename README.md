# planner

Planner helps development teams prioritize their works and communicate realistic task completion forecasts.

It supports a very narrow set of use cases (basically only mine at the moment), and this is the polar opposite of a general purpose project management tool. More precisely, it follows the following usage principles:
- The edition is purely text-based. For me, it's quicker to edit, easier to automate and version.
- Support some kind of visual output, namely in the form of Gantt chart. While I dislike editing anything other than text, I also think that it's often a poor medium to summarize and communicate ideas.
- Support the same format for its input and its output, as explained in the workflow section.
- Doesn't support multiple conflicting constraints on the same item. For example, you can only specify the effort needed to complete each task, not the desired completion dates. It's limiting, but more simple and enough for my personal needs.

# workflow

- Install planner
```shell script
go get github.com/ostapneko/planner/planner
```
- Create a planning file. See the dedicated section for the specs.
- Run planner to generate completion dates and optionally a Gantt chart.
```shell script
planner -o output-planning.yaml -g gantt input-planning.yaml
```
- Run [PlantUML](https://plantuml.com/) to generate a visual output:
```shell script
java -jar plantuml.jar gantt
```
Example output

![Gantt chart](doc/example-gantt.png "Gantt chart")