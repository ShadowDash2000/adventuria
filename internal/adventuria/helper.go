package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

func GetRecordById(table, id string, expand []string) (*core.Record, error) {
	record, err := PocketBase.FindRecordById(GameCollections.Get(table), id)
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

func UpdateRecordsFromViewCollection(
	records []*core.Record,
	dstCollection *core.Collection,
	primaryKey string,
	fieldsToUpdate []string,
) error {
	var dstRecords []*core.Record
	for _, record := range records {
		dstRecord, err := PocketBase.FindRecordById(dstCollection, record.GetString(primaryKey))
		if err != nil {
			return err
		}

		for _, field := range fieldsToUpdate {
			dstRecord.Set(field, record.Get(field))
		}

		dstRecords = append(dstRecords, dstRecord)
	}

	for _, dstRecord := range dstRecords {
		err := PocketBase.Save(dstRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

// normalized mod (0..n-1)
func mod(a, m int) int {
	return ((a % m) + m) % m
}

// floor division
func floorDiv(a, m int) int {
	r := mod(a, m)
	return (a - r) / m
}
