package task

type Tasks struct {
	Items []Task `json:"items"`
}

type SearchedTasks struct {
	Items []SearchedTask `json:"items"`
}
