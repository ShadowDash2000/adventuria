package result

type Result struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Ok() *Result {
	return &Result{
		Success: true,
	}
}

func Err(err string) *Result {
	return &Result{
		Success: false,
		Error:   err,
	}
}

func (r *Result) WithData(data any) *Result {
	r.Data = data
	return r
}

func (r *Result) HasError() bool {
	return r.Error != ""
}

func (r *Result) Ok() bool {
	if r == nil {
		return false
	}
	return r.Success
}

func (r *Result) Failed() bool {
	if r == nil {
		return false
	}
	return !r.Success || r.Error != ""
}
