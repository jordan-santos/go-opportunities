package handler

import (
	"fmt"
)

type CreateOpeningRequest struct {
	Role     string `json:"role"`
	Company  string `json:"company"`
	Location string `json:"location"`
	Remote   *bool  `json:"remote"`
	Link     string `json:"link"`
	Salary   int64  `json:"salary"`
}

func (req *CreateOpeningRequest) Validate() error {
	if req.Role == "" {
		return errParamIsRequired("role", "string")
	}

	if req.Company == "" {
		return errParamIsRequired("company", "string")
	}

	if req.Location == "" {
		return errParamIsRequired("location", "string")
	}

	if req.Remote == nil {
		return errParamIsRequired("remote", "bool")
	}

	if req.Link == "" {
		return errParamIsRequired("link", "string")
	}

	if req.Salary <= 0 {
		return errParamIsRequired("salary", "int")
	}

	return nil
}

func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param: %s (type: %s) is required", name, typ)
}

type UpdateOpeningRequest struct {
	Role     string `json:"role"`
	Company  string `json:"company"`
	Location string `json:"location"`
	Remote   *bool  `json:"remote"`
	Link     string `json:"link"`
	Salary   int64  `json:"salary"`
}

func (req *UpdateOpeningRequest) Validate() error {
	if req.Role != "" || req.Company != "" || req.Location != "" || req.Remote != nil || req.Salary > 0 {
		return nil
	}

	return fmt.Errorf("at least one param is required")
}
