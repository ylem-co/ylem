package readmodel

type Period string

const (
	DAY     Period = "day"
	WEEK    Period = "week"
	MONTH   Period = "month"
	QUARTER Period = "quarter"
	YEAR    Period = "year"
)

var ValidPeriods = []interface{}{
	string(DAY),
	string(WEEK),
	string(MONTH),
	string(QUARTER),
	string(YEAR),
}
