package params

type CResponse struct {
	FileName string     `json:"file_name"`
	Result   string     `json:"result"`
	Message  string     `json:"mgs"`
	Hints    []HintsMap `json:"hints_map"`
}

type HintsMap struct {
	Hints     string `json:"hints"`
	CheckName string `json:"check_name"`
	Level     int32  `json:"level"`
}
