package custom

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/activities"
	"adventuria/internal/adventuria/activity_filters"
	"adventuria/internal/adventuria/board"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/effects/custom/add_game_genre"
	"adventuria/internal/adventuria/effects/custom/add_items_to_inventory"
	"adventuria/internal/adventuria/effects/custom/balance_change"
	"adventuria/internal/adventuria/effects/custom/cell_points_divide"
	"adventuria/internal/adventuria/effects/custom/change_dices"
	"adventuria/internal/adventuria/effects/custom/change_game_by_id"
	"adventuria/internal/adventuria/effects/custom/change_game_price_filter"
	"adventuria/internal/adventuria/effects/custom/choose_activity"
	"adventuria/internal/adventuria/effects/custom/coins_for_all"
	"adventuria/internal/adventuria/effects/custom/debuff_block"
	"adventuria/internal/adventuria/effects/custom/discount_price_divide"
	"adventuria/internal/adventuria/effects/custom/drop_block"
	"adventuria/internal/adventuria/effects/custom/drop_inventory"
	"adventuria/internal/adventuria/effects/custom/give_wheel_on_done"
	"adventuria/internal/adventuria/effects/custom/give_wheel_on_new_lap"
	"adventuria/internal/adventuria/effects/custom/go_to_jail"
	"adventuria/internal/adventuria/effects/custom/jail_escape"
	"adventuria/internal/adventuria/effects/custom/no_coins_for_done"
	"adventuria/internal/adventuria/effects/custom/no_time_limit"
	"adventuria/internal/adventuria/effects/custom/nothing"
	"adventuria/internal/adventuria/effects/custom/paid_movement_in_radius"
	"adventuria/internal/adventuria/effects/custom/points_change"
	"adventuria/internal/adventuria/effects/custom/replace_dice_roll"
	"adventuria/internal/adventuria/effects/custom/reroll_block"
	"adventuria/internal/adventuria/effects/custom/return_to_prev_cell"
	"adventuria/internal/adventuria/effects/custom/roll_reverse"
	"adventuria/internal/adventuria/effects/custom/safe_drop"
	"adventuria/internal/adventuria/effects/custom/save_from_hole"
	"adventuria/internal/adventuria/effects/custom/stay_on_cell_after_done"
	"adventuria/internal/adventuria/effects/custom/teleport_to_closest_cell_by_type"
	"adventuria/internal/adventuria/effects/custom/teleport_to_random_cell"
	"adventuria/internal/adventuria/genres"
	"adventuria/internal/adventuria/inventories"
	"adventuria/internal/adventuria/items"
	"adventuria/internal/adventuria/outboxes"
	"adventuria/internal/adventuria/players"
)

func RegisterEffects(
	actions *actions.Actions,
	cells *cells.Cells,
	genres *genres.Genres,
	activityFilters *activity_filters.ActivityFilters,
	inventories *inventories.Inventories,
	items *items.Items,
	activities *activities.Activities,
	players *players.Players,
	outboxes *outboxes.Outboxes,
	board *board.Board,
) {
	effects.Register(
		add_game_genre.NewDef(actions, cells, genres, activityFilters),
		add_items_to_inventory.NewDef(inventories, items),
		cell_points_divide.NewDef(),
		change_dices.NewDef(),
		change_game_by_id.NewDef(cells, actions, activities),
		change_game_price_filter.NewDef(actions, cells, activityFilters),
		choose_activity.NewDef(actions, activities),
		coins_for_all.NewDef(players, outboxes),
		balance_change.NewDef(),
		debuff_block.NewDef(),
		discount_price_divide.NewDef(),
		drop_block.NewDef(),
		drop_inventory.NewDef(inventories),
		go_to_jail.NewDef(actions, board),
		jail_escape.NewDef(),
		no_coins_for_done.NewDef(),
		no_time_limit.NewDef(actions, cells),
		nothing.NewDef(cells),
		paid_movement_in_radius.NewDef(actions, cells, board),
		points_change.NewDef(),
		replace_dice_roll.NewDef(),
		reroll_block.NewDef(),
		return_to_prev_cell.NewDef(actions, board),
		roll_reverse.NewDef(),
		safe_drop.NewDef(cells),
		save_from_hole.NewDef(cells),
		stay_on_cell_after_done.NewDef(cells, actions),
		teleport_to_random_cell.NewDef(actions, cells, board),
		teleport_to_closest_cell_by_type.NewDef(actions, board),
	)
}

func RegisterPersistentEffects() {
	effects.RegisterPersistent(
		give_wheel_on_done.NewDef(),
		give_wheel_on_new_lap.NewDef(),
	)
}
