package game_state

import "adventuria/internal/adventuria/npcs"

type activeNPCView struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Color  string `json:"color"`
}

func activeNPCToView(npc *npcs.ActiveNPC) *activeNPCView {
	return &activeNPCView{
		Id:     npc.Id,
		Name:   npc.Name,
		Avatar: npc.Avatar,
		Color:  npc.Color,
	}
}

func activeNPCsToView(npcs []*npcs.ActiveNPC) []*activeNPCView {
	res := make([]*activeNPCView, len(npcs))
	for i, npc := range npcs {
		res[i] = activeNPCToView(npc)
	}
	return res
}
