package condition

import (
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
)

type Condition struct {
	Less int    `json:"less"`
	More int    `json:"more"`
	Is   int    `json:"is"`
	Key  string `json:"key"`
	src  any
}

func NewCondition(s string, src any) *Condition {
	cond := new(Condition)
	json.Unmarshal([]byte(s), &cond)
	cond.src = src
	return cond
}

func (c *Condition) Compare() (bool, error) {
	var srcValue int

	switch c.src.(type) {
	case int:
		srcValue = c.src.(int)
	case struct{}:
		if c.Key == "" {
			return false, errors.New("empty condition key")
		}

		var srcMap map[string]any
		var ok bool
		mapstructure.Decode(c.src, &srcMap)
		srcValue, ok = srcMap[c.Key].(int)
		if !ok {
			return false, errors.New("key not found in src")
		}
	default:
		panic("invalid condition src type")
	}

	res := false
	if srcValue < c.Less {
		res = true
	}
	if srcValue > c.More {
		res = true
	}
	if srcValue == c.Is {
		res = true
	}

	return res, nil
}
