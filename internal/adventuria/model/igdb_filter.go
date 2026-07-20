package model

import (
	"fmt"
	"strings"
	"time"
)

type IgdbFilter struct {
	GameTypes      []string
	Platforms      []string
	ReleaseDateMin time.Time
	ReleaseDateMax time.Time
}

func (i IgdbFilter) Build() string {
	var entries []string

	if len(i.GameTypes) > 0 {
		entries = append(entries, fmt.Sprintf("game_type = %s", i.buildFromSlice(i.GameTypes)))
	}
	if len(i.Platforms) > 0 {
		entries = append(entries, fmt.Sprintf("platforms = %s", i.buildFromSlice(i.Platforms)))
	}
	if !i.ReleaseDateMin.IsZero() {
		entries = append(entries, fmt.Sprintf("first_release_date > %d", i.ReleaseDateMin.Unix()))
	}
	if !i.ReleaseDateMax.IsZero() {
		entries = append(entries, fmt.Sprintf("first_release_date < %d", i.ReleaseDateMax.Unix()))
	}

	return strings.Join(entries, " & ")
}

func (i IgdbFilter) buildFromSlice(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ","))
}
