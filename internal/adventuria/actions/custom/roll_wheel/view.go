package roll_wheel

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*RollWheel)(nil)

func (r *RollWheel) GetView(ctx context.Context, _ *model.Events, player *model.Player) (any, error) {
	currentCell, err := r.cells.GetByPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	activitiesState := player.LastAction().State().Activities
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
