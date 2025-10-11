package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

func GetRecordById(table, id string, expand []string) (*core.Record, error) {
	collection, err := GameCollections.Get(table)
	if err != nil {
		return nil, err
	}

	record, err := PocketBase.FindRecordById(collection, id)
	if err != nil {
		return nil, err
	}

	if expand != nil {
		errs := PocketBase.ExpandRecord(record, expand, nil)
		if errs != nil {
			for _, err := range errs {
				return nil, err
			}
		}
	}

	return record, nil
}
