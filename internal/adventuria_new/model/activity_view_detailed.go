package model

type ActivityViewDetailed struct {
	activity   *Activity
	platforms  []*Platform
	developers []*Developer
	publishers []*Publisher
	genres     []*Genre
	tags       []*Tag
	themes     []*Theme
}

func RestoreActivityViewDetailed(
	activity *Activity,
	platforms []*Platform,
	developers []*Developer,
	publishers []*Publisher,
	genres []*Genre,
	tags []*Tag,
	themes []*Theme,
) *ActivityViewDetailed {
	return &ActivityViewDetailed{
		activity:   activity,
		platforms:  platforms,
		developers: developers,
		publishers: publishers,
		genres:     genres,
		tags:       tags,
		themes:     themes,
	}
}

func (a *ActivityViewDetailed) Activity() *Activity {
	return a.activity
}

func (a *ActivityViewDetailed) Platforms() []*Platform {
	return a.platforms
}

func (a *ActivityViewDetailed) Developers() []*Developer {
	return a.developers
}

func (a *ActivityViewDetailed) Publishers() []*Publisher {
	return a.publishers
}

func (a *ActivityViewDetailed) Genres() []*Genre {
	return a.genres
}

func (a *ActivityViewDetailed) Tags() []*Tag {
	return a.tags
}

func (a *ActivityViewDetailed) Themes() []*Theme {
	return a.themes
}
