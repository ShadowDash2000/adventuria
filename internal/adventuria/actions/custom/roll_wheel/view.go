package roll_wheel

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*RollWheel)(nil)

func (r *RollWheel) GetView(ctx context.Context, _ *model.Events, player *model.Player) (any, error) {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	activities, err := r.activities.GetDetailedByIDs(ctx, player.LastAction().ItemsList())
	if err != nil {
		return nil, err
	}

	return struct {
		Items         []activityViewDetailed `json:"items"`
		AudioPresetId string                 `json:"audio_preset_id,omitempty"`
	}{
		Items:         toActivityViewDetailedList(activities),
		AudioPresetId: currentCell.Data().AudioPreset(),
	}, nil
}
