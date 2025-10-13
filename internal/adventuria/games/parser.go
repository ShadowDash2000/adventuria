package games

type Parser interface {
	ParseGames() (chan []Game, error)
	ParsePlatforms() (chan []Platform, error)
	ParseCompanies() (chan []Company, error)
}
