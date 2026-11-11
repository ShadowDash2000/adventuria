package repository

import (
	effectsRepo "adventuria/internal/adventuria/effects/repository"
	"adventuria/internal/adventuria/errs"
	inventoryModel "adventuria/internal/adventuria/inventories/model"
	itemsRepo "adventuria/internal/adventuria/items/repository"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionInventory)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	InventoryToRecord(inventory, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToInventory(record), nil
}

func (r *Repository) Update(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionInventory, inventory.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrInventoryNotFound
		}
		return nil, err
	}

	InventoryToRecord(inventory, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToInventory(record), nil
}

func (r *Repository) GetAllByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItem, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Where(dbx.HashExp{schema.InventorySchema.Player: playerId}).
		All(&records)
	if err != nil {
		return nil, err
	}

	errs := pb.ExpandRecords(records, []string{schema.InventorySchema.Item}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand records: %v", errs)
	}

	items := make([]*model.InventoryItem, len(records))
	for i, record := range records {
		inventory := RecordToInventory(record)
		item := itemsRepo.RecordToItem(record.ExpandedOne(schema.InventorySchema.Item))
		items[i] = model.RestoreInventoryItem(inventory, item)
	}

	return items, nil
}

func (r *Repository) GetAllDetailedByPlayerID(ctx context.Context, playerId string) ([]*inventoryModel.InventoryItemDetails, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Where(dbx.HashExp{schema.InventorySchema.Player: playerId}).
		OrderBy(schema.InventorySchema.Activated + " ASC").
		All(&records)
	if err != nil {
		return nil, err
	}

	errs := pb.ExpandRecords(records, []string{
		schema.InventorySchema.Item,
		pbhelper.DotExpand(schema.InventorySchema.Item, schema.ItemSchema.Effects),
	}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand records: %v", errs)
	}

	items := make([]*inventoryModel.InventoryItemDetails, len(records))
	for i, inventoryRecord := range records {
		inventory := RecordToInventory(inventoryRecord)

		itemRecord := inventoryRecord.ExpandedOne(schema.InventorySchema.Item)
		item := itemsRepo.RecordToItem(inventoryRecord.ExpandedOne(schema.InventorySchema.Item))

		effects := effectsRepo.RecordsToEffects(itemRecord.ExpandedAll(schema.ItemSchema.Effects))

		items[i] = &inventoryModel.InventoryItemDetails{
			Created:   inventoryRecord.GetDateTime("created").Time(),
			IconPath:  pbhelper.BuildFilePath(schema.CollectionItems, item.ID(), item.Icon()),
			Inventory: inventory,
			Item:      item,
			Effects:   effects,
		}
	}

	return items, nil
}

func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Where(dbx.HashExp{schema.InventorySchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		return err
	}

	err = pb.DeleteWithContext(ctx, &record)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteByIDs(ctx context.Context, ids []string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Where(dbx.In(
			schema.InventorySchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return err
	}

	for _, record := range records {
		err = pb.DeleteWithContext(ctx, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.Inventory, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Where(dbx.HashExp{schema.InventorySchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrInventoryNotFound
		}
		return nil, err
	}

	return RecordToInventory(&record), nil
}

func (r *Repository) GetPlayerInventoryItemByID(ctx context.Context, playerId, itemId string) (*model.InventoryItem, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.InventorySchema.Id:     itemId,
			schema.InventorySchema.Player: playerId,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrInventoryNotFound
		}
		return nil, err
	}

	errs := pb.ExpandRecord(&record, []string{schema.InventorySchema.Item}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand record: %v", errs)
	}

	inventory := RecordToInventory(&record)
	item := itemsRepo.RecordToItem(record.ExpandedOne(schema.InventorySchema.Item))

	return model.RestoreInventoryItem(inventory, item), nil
}

func (r *Repository) GetPlayerInventoryItemWithEffectsByID(ctx context.Context, playerId, itemId string) (*inventoryModel.InventoryItemDetails, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var inventoryRecord core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.InventorySchema.Id:     itemId,
			schema.InventorySchema.Player: playerId,
		}).
		Limit(1).
		One(&inventoryRecord)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrInventoryNotFound
		}
		return nil, err
	}

	errs := pb.ExpandRecord(&inventoryRecord, []string{
		schema.InventorySchema.Item,
		pbhelper.DotExpand(schema.InventorySchema.Item, schema.ItemSchema.Effects),
	}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand record: %v", errs)
	}

	inventory := RecordToInventory(&inventoryRecord)

	itemRecord := inventoryRecord.ExpandedOne(schema.InventorySchema.Item)
	item := itemsRepo.RecordToItem(itemRecord)

	effects := effectsRepo.RecordsToEffects(itemRecord.ExpandedAll(schema.ItemSchema.Effects))

	return &inventoryModel.InventoryItemDetails{
		Created:   inventoryRecord.GetDateTime("created").Time(),
		IconPath:  pbhelper.BuildFilePath(schema.CollectionItems, item.ID(), item.Icon()),
		Inventory: inventory,
		Item:      item,
		Effects:   effects,
	}, nil
}

func (r *Repository) GetPlayerUsedSlots(ctx context.Context, playerId string) (int, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var count int
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Select("count(*)").
		InnerJoin(
			schema.CollectionItems,
			dbx.NewExp(
				pbhelper.Eq(
					pbhelper.DotExpand(schema.CollectionInventory, schema.InventorySchema.Item),
					pbhelper.DotExpand(schema.CollectionItems, schema.ItemSchema.Id),
				),
			),
		).
		Where(dbx.HashExp{
			pbhelper.DotExpand(schema.CollectionInventory, schema.InventorySchema.Player): playerId,
			pbhelper.DotExpand(schema.CollectionItems, schema.ItemSchema.IsUsingSlot):     true,
		}).
		Row(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetAllDroppableUsingSlotByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItem, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionInventory).
		WithContext(ctx).
		Select(pbhelper.DotExpand(schema.CollectionInventory, "*")).
		InnerJoin(
			schema.CollectionItems,
			dbx.NewExp(
				pbhelper.Eq(
					pbhelper.DotExpand(schema.CollectionInventory, schema.InventorySchema.Item),
					pbhelper.DotExpand(schema.CollectionItems, schema.ItemSchema.Id),
				),
			),
		).
		Where(dbx.HashExp{
			pbhelper.DotExpand(schema.CollectionInventory, schema.InventorySchema.Player): playerId,
			pbhelper.DotExpand(schema.CollectionItems, schema.ItemSchema.CanDrop):         true,
			pbhelper.DotExpand(schema.CollectionItems, schema.ItemSchema.IsUsingSlot):     true,
		}).
		All(&records)
	if err != nil {
		return nil, err
	}

	errs := pb.ExpandRecords(records, []string{schema.InventorySchema.Item}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand records: %v", errs)
	}

	items := make([]*model.InventoryItem, len(records))
	for i, record := range records {
		inventory := RecordToInventory(record)
		item := itemsRepo.RecordToItem(record.ExpandedOne(schema.InventorySchema.Item))
		items[i] = model.RestoreInventoryItem(inventory, item)
	}

	return items, nil
}
