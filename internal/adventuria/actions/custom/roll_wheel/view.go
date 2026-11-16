package roll_wheel

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"context"
)

var _ actions.View = (*RollWheel)(nil)

func (r *RollWheel) GetView(ctx context.Context, actor *player.Actor) (any, error) {
	currentCell, err := r.cells.GetByPlayer(ctx, actor.Player())
	if err != nil {
		return nil, err
	}

	activitiesState := actor.Player().LastAction().State().Activities
	if activitiesState == nil {
		return nil, errs.ErrNoActiveActivity
	}

	activities, err := r.activities.GetDetailedByIDs(ctx, activitiesState.Ids)
	if err != nil {
		return nil, err
	}

	return struct {
		Items         []activityViewDetailed `json:"items"`
		AudioPresetId string                 `json:"audio_preset_id,omitempty"`
	}{
		Items:         toActivityViewDetailedList(activities),
		AudioPresetId: currentCell.AudioPreset(),
	}, nil
}
