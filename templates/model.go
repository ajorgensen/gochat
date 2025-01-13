package templates

var Anthropic = Provider{
	Name: "Anthropic",
	Models: []Model{
		{
			Name: "Test",
		},
	},
}

var Robot = Provider{
	Name: "Robot",
	Models: []Model{
		{
			Name: "repeater",
		},
	},
}

var Providers = []Provider{
	Anthropic, Robot,
}

type Provider struct {
	Name   string
	Models []Model
}

type Model struct {
	Name string
}
