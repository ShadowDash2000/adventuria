package model

type ActivityDetailed struct {
	activity   *Activity
	platforms  []*Platform
	developers []*Company
	publishers []*Company
	genres     []*Genre
	tags       []*Tag
	themes     []*Theme
}

func RestoreActivityDetailed(
	activity *Activity,
	platforms []*Platform,
	developers []*Company,
	publishers []*Company,
	genres []*Genre,
	tags []*Tag,
	themes []*Theme,
) *ActivityDetailed {
	return &ActivityDetailed{
		activity:   activity,
		platforms:  platforms,
		developers: developers,
		publishers: publishers,
		genres:     genres,
		tags:       tags,
		themes:     themes,
	}
}

func (a *ActivityDetailed) Activity() *Activity {
	return a.activity
}

func (a *ActivityDetailed) Platforms() []*Platform {
	return a.platforms
}

func (a *ActivityDetailed) Developers() []*Company {
	return a.developers
}

func (a *ActivityDetailed) Publishers() []*Company {
	return a.publishers
}

func (a *ActivityDetailed) Genres() []*Genre {
	return a.genres
}

func (a *ActivityDetailed) Tags() []*Tag {
	return a.tags
}

func (a *ActivityDetailed) Themes() []*Theme {
	return a.themes
}
