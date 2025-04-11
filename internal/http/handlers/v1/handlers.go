package handlers

import (
	"adventuria/internal/adventuria"
	"errors"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"net/http"
)

type Handlers struct {
	Game adventuria.Game
}

func New(g adventuria.Game) *Handlers {
	return &Handlers{Game: g}
}

func (h *Handlers) RollHandler(e *core.RequestEvent) error {
	n, diceRolls, currentCell, err := h.Game.Roll(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"roll":      n,
		"diceRolls": diceRolls,
		"cellId":    currentCell.ID(),
	})
	return nil
}

func (h *Handlers) ChooseGameHandler(e *core.RequestEvent) error {
	data := struct {
		Game string `json:"game"`
	}{}
	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	if data.Game == "" {
		e.JSON(http.StatusBadRequest, "You must choose a game")
		return nil
	}

	err := h.Game.ChooseGame(data.Game, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GetNextStepTypeHandler(e *core.RequestEvent) error {
	nextStepType, err := h.Game.GetNextStepType(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"nextStepType": nextStepType,
	})

	return nil
}

func (h *Handlers) UpdateActionHandler(e *core.RequestEvent) error {
	data := struct {
		ActionID string `form:"actionId"`
		Comment  string `form:"comment"`
	}{}

	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	var file *filesystem.File
	files, err := e.FindUploadedFiles("icon")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	} else if len(files) > 0 {
		file = files[0]
	}

	err = h.Game.UpdateAction(data.ActionID, data.Comment, file, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) RerollHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `form:"comment"`
	}{}

	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	var file *filesystem.File
	files, err := e.FindUploadedFiles("result-file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	} else if len(files) > 0 {
		file = files[0]
	}

	err = h.Game.Reroll(data.Comment, file, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) DropHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `form:"comment"`
	}{}

	err := e.BindBody(&data)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	var file *filesystem.File
	files, err := e.FindUploadedFiles("result-file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	} else if len(files) > 0 {
		file = files[0]
	}

	err = h.Game.Drop(data.Comment, file, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})

	return nil
}

func (h *Handlers) DoneHandler(e *core.RequestEvent) error {
	data := struct {
		Comment string `form:"comment"`
	}{}

	if err := e.BindBody(&data); err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	var file *filesystem.File
	files, err := e.FindUploadedFiles("result-file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		e.JSON(http.StatusInternalServerError, err.Error())
		return nil
	} else if len(files) > 0 {
		file = files[0]
	}

	err = h.Game.Done(data.Comment, file, e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, struct{}{})
	return nil
}

func (h *Handlers) GetLastActionHandler(e *core.RequestEvent) error {
	isInJail, action, err := h.Game.GetLastAction(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"title":    action.Value(),
		"isInJail": isInJail,
	})

	return nil
}

func (h *Handlers) GetRollEffectsHandler(e *core.RequestEvent) error {
	effects, err := h.Game.GetItemsEffects(e.Auth.Id, adventuria.EffectUseOnRoll)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, effects.Map())

	return nil
}

func (h *Handlers) UseItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"itemId"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err = h.Game.UseItem(e.Auth.Id, data.ItemId)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) DropItemHandler(e *core.RequestEvent) error {
	data := struct {
		ItemId string `json:"itemId"`
	}{}

	err := e.BindBody(&data)
	if err != nil {
		e.JSON(http.StatusBadRequest, err.Error())
		return nil
	}

	err = h.Game.DropItem(e.Auth.Id, data.ItemId)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) StartTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StartTimer(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) StopTimerHandler(e *core.RequestEvent) error {
	err := h.Game.StopTimer(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, "")

	return nil
}

func (h *Handlers) GetTimeLeftHandler(e *core.RequestEvent) error {
	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(e.Auth.Id)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"time":               time,
		"isActive":           isActive,
		"nextTimerResetDate": nextTimerResetDate,
	})

	return nil
}

func (h *Handlers) GetTimeLeftByUserHandler(e *core.RequestEvent) error {
	userId := e.Request.PathValue("userId")

	time, isActive, nextTimerResetDate, err := h.Game.GetTimeLeft(userId)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"time":               time,
		"isActive":           isActive,
		"nextTimerResetDate": nextTimerResetDate,
	})

	return nil
}
