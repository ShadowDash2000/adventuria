package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

func GetRecordById(locator PocketBaseLocator, table, id string, expand []string) (*core.Record, error) {
	collection, err := locator.Collections().Get(table)
	if err != nil {
		return nil, err
	}

	record, err := locator.PocketBase().FindRecordById(collection, id)
	if err != nil {
		return nil, err
	}

	if expand != nil {
		errs := locator.PocketBase().ExpandRecord(record, expand, nil)
		if errs != nil {
			for _, err := range errs {
				return nil, err
			}
		}
	}

	return record, nil
}
