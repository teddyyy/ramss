package model

// Unit ...
type Unit struct {
	ID            string `json:"Id"`
	Description   string `json:"Description"`
	LoadState     string `json:"LoadState"`
	ActiveState   string `json:"ActiveState"`
	UnitFileState string `json:"UnitFileState"`
	MainPID       int    `json:"MainPID"`
}

// Units ...
type Units struct {
	Units []*Unit `json:"Systemd"`
}

// APISuccess ...
type APISuccess struct {
	Code    int
	Message string
}

// APIError ...
type APIError struct {
	Code    int
	Message string
}
