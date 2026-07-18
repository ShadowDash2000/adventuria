package errs

import "net/http"

type AppError struct {
	Code       string
	Message    string
	Status     int
	Translates map[string]string
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) GetMessage(lang string) string {
	if msg, ok := e.Translates[lang]; ok {
		return msg
	}
	return e.Message
}

var (
	ErrSettingsNotFound = &AppError{
		Code:    "settings_not_found",
		Message: "Settings not found",
	}
	ErrSeasonNotFound = &AppError{
		Code:    "season_not_found",
		Message: "Season not found",
	}
	ErrPlayerNotFound = &AppError{
		Code:    "player_not_found",
		Message: "Player not found",
	}
	ErrProgressNotFound = &AppError{
		Code:    "progress_not_found",
		Message: "Player progress not found",
	}
	ErrActionNotFound = &AppError{
		Code:    "action_not_found",
		Message: "Action not found",
	}
	ErrItemNotFound = &AppError{
		Code:    "item_not_found",
		Message: "Item not found",
	}
	ErrInventoryNotFound = &AppError{
		Code:    "inventory_not_found",
		Message: "Inventory not found",
	}
	ErrGenreNotFound = &AppError{
		Code:    "genre_not_found",
		Message: "Genre not found",
	}
	ErrPlatformNotFound = &AppError{
		Code:    "platform_not_found",
		Message: "Platform not found",
	}
	ErrActivityFilterNotFound = &AppError{
		Code:    "activity_filter_not_found",
		Message: "Activity filter not found",
	}
	ErrCellNotFound = &AppError{
		Code:    "cell_not_found",
		Message: "Cell not found",
	}
	ErrReviewNotFound = &AppError{
		Code:    "review_not_found",
		Message: "Review not found",
	}
	ErrOutboxNotFound = &AppError{
		Code:    "outbox_not_found",
		Message: "Outbox not found",
	}
	ErrEffectNotFound = &AppError{
		Code:    "effect_not_found",
		Message: "Effect not found",
	}
	ErrCheapSharkNotFound = &AppError{
		Code:    "cheapshark_not_found",
		Message: "Cheapshark not found",
	}
	ErrHowLongToBeatNotFound = &AppError{
		Code:    "howlongtobeat_not_found",
		Message: "HowLongToBeat not found",
	}
	ErrSteamSpyNotFound = &AppError{
		Code:    "steam_spy_not_found",
		Message: "Steam Spy not found",
	}
	ErrCompanyNotFound = &AppError{
		Code:    "company_not_found",
		Message: "Company not found",
	}
	ErrTagNotFound = &AppError{
		Code:    "tag_not_found",
		Message: "Tag not found",
	}
	ErrThemeNotFound = &AppError{
		Code:    "theme_not_found",
		Message: "Theme not found",
	}
	ErrGameTypeNotFound = &AppError{
		Code:    "game_type_not_found",
		Message: "Game type not found",
	}
	ErrActivityNotFound = &AppError{
		Code:    "activity_not_found",
		Message: "Activity not found",
	}
	ErrPlayerStatsNotFound = &AppError{
		Code:    "player_stats_not_found",
		Message: "Player stats not found",
	}
	ErrActionEventNotFound = &AppError{
		Code:    "cell_event_not_found",
		Message: "Cell event not found",
	}

	ErrReviewCommentTooLong = &AppError{
		Code:    "review_comment_max_size",
		Message: "Review comment is too long",
		Status:  http.StatusBadRequest,
	}
	ErrReviewScoreInvalid = &AppError{
		Code:    "review_score_invalid",
		Message: "Invalid review score",
		Status:  http.StatusBadRequest,
	}

	ErrNotEnoughMoney = &AppError{
		Code:    "not_enough_money",
		Message: "Not enough money",
		Status:  http.StatusBadRequest,
	}
	ErrNotEnoughEnergy = &AppError{
		Code:    "not_enough_energy",
		Message: "Not enough energy",
		Status:  http.StatusBadRequest,
	}

	ErrUnknownAction = &AppError{
		Code:    "unknown_action",
		Message: "Unknown action",
		Status:  http.StatusBadRequest,
	}

	ErrUnknownCellType = &AppError{
		Code:    "unknown_cell_type",
		Message: "Unknown cell type",
	}
	ErrUnknownEffectType = &AppError{
		Code:    "unknown_effect_type",
		Message: "Unknown effect type",
	}
	ErrUnknownActionEventType = &AppError{
		Code:    "unknown_action_event_type",
		Message: "Unknown action event type",
	}

	ErrNoPendingOutbox = &AppError{
		Code:    "no_pending_outbox",
		Message: "No pending outbox",
	}

	ErrPlayerIsBusy = &AppError{
		Code:    "player_is_busy",
		Message: "Player is busy",
		Status:  http.StatusConflict,
	}

	ErrDontDoThat = &AppError{
		Code:    "dont_do_that",
		Message: "Don't do that",
		Status:  http.StatusNotImplemented,
	}

	ErrActionIsNotEventCompatible = &AppError{
		Code:    "action_is_not_event_compatible",
		Message: "Action is not event compatible",
	}
)
