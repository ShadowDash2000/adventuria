package steam

import steamstore "github.com/ShadowDash2000/steam-store-go"

type Parser struct {
	client *steamstore.Client
}

func NewParser() *Parser {
	return &Parser{
		client: steamstore.New(),
	}
}
