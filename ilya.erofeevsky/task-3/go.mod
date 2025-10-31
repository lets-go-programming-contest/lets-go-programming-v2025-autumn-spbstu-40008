module github.com/task-3

go 1.22

require (
	golang.org/x/text v0.14.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/task-3/config => ./config
replace github.com/task-3/internal/data => ./internal/data
replace github.com/task-3/internal/structures => ./internal/structures