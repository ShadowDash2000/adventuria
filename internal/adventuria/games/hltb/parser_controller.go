package hltb

type ParserController struct {
	parser *Parser
}

func New() (*ParserController, error) {
	p, err := NewParser()
	if err != nil {
		return nil, err
	}

	return &ParserController{parser: p}, nil
}
