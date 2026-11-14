package behaviors

import (
	"bytes"
)

type BehaviorBase struct {
	payloadJSON []byte
	dataJSON    []byte
}

func NewBehaviorBase(payloadJSON, dataJSON []byte) BehaviorBase {
	return BehaviorBase{
		payloadJSON: payloadJSON,
		dataJSON:    dataJSON,
	}
}

func (p BehaviorBase) Payload() []byte {
	return bytes.Clone(p.payloadJSON)
}

func (p BehaviorBase) Data() []byte {
	return bytes.Clone(p.dataJSON)
}
