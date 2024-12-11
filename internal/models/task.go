package models

type Task struct {
	ID        uint32 `json:"task_id"`
	InputData string `json:"input_data"`
	Status    string `json:"status"`
	Result    string `json:"result"`
}
