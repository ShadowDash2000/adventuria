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
		Message: "season not found",
	}
	ErrPlayerNotFound = &AppError{
		Code:    "player_not_found",
		Message: "player not found",
	}
	ErrProgressNotFound = &AppError{
		Code:    "progress_not_found",
		Message: "player progress not found",
	}
	ErrActionNotFound = &AppError{
		Code:    "action_not_found",
		Message: "action not found",
	}
	ErrItemNotFound = &AppError{
		Code:    "item_not_found",
		Message: "item not found",
	}
	ErrInventoryNotFound = &AppError{
		Code:    "inventory_not_found",
		Message: "inventory not found",
	}
	ErrGenreNotFound = &AppError{
		Code:    "genre_not_found",
		Message: "genre not found",
	}
	ErrActivityFilterNotFound = &AppError{
		Code:    "activity_filter_not_found",
		Message: "activity filter not found",
	}
	ErrCellNotFound = &AppError{
		Code:    "cell_not_found",
		Message: "cell not found",
	}
	ErrReviewNotFound = &AppError{
		Code:    "review_not_found",
		Message: "review not found",
	}
	ErrOutboxNotFound = &AppError{
		Code:    "outbox_not_found",
		Message: "outbox not found",
	}
	ErrEffectNotFound = &AppError{
		Code:    "effect_not_found",
		Message: "effect not found",
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

	ErrNoPendingOutbox = &AppError{
		Code:    "no_pending_outbox",
		Message: "No pending outbox",
	}

	ErrPlayerIsBusy = &AppError{
		Code:    "player_is_busy",
		Message: "Player is busy",
		Status:  http.StatusConflict,
	}
)
