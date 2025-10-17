package hltb

import "github.com/forbiddencoding/howlongtobeat"

type Parser struct {
	client *howlongtobeat.Client
}

func NewParser() (*Parser, error) {
	c, err := howlongtobeat.New()
	if err != nil {
		return nil, err
	}

	return &Parser{client: c}, nil
}
