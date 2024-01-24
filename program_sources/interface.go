package program_sources

type ProgramSourceInterface interface {
	GetName() string
	Configure(config map[string]interface{})
	GetPrograms(bountyOnly bool) chan Program
}

var ProgramSources = make(map[string]ProgramSourceInterface)

type Program struct {
	Name     string
	Platform string
	Assets   []Asset
	Reward   bool
}

type Asset struct {
	Domain   string
	Wildcard bool
}
