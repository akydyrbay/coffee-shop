package models

import "fmt"

type Status struct {
	ErrorMessage error `json:"error,omitempty"`
	Code         int   `json:"code"`
}

var (
	IdAlreadyExists     Status = Status{ErrorMessage: fmt.Errorf("ID already exists"), Code: 409}
	NoIdGiven           Status = Status{ErrorMessage: fmt.Errorf("no ID given"), Code: 404}
	NotFound            Status = Status{ErrorMessage: fmt.Errorf("not found"), Code: 404}
	Success             Status = Status{ErrorMessage: nil, Code: 200}
	Created             Status = Status{ErrorMessage: nil, Code: 201}
	BadRequest          Status = Status{ErrorMessage: fmt.Errorf("bad request"), Code: 400}
	InternalServerError Status = Status{ErrorMessage: fmt.Errorf("internal server error"), Code: 500}
	Conflict            Status = Status{ErrorMessage: fmt.Errorf("conflict"), Code: 409}
	Deleted             Status = Status{ErrorMessage: nil, Code: 204}
)
