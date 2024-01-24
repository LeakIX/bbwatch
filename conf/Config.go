package conf

type BBWatchConfig struct {
	VulnerabilitySources map[string]VulnerabilitySources `yaml:"vulnerability_sources"`
	ProgramSources       map[string]ProgramSources       `yaml:"program_sources"`
}

type VulnerabilitySources map[string]interface{}
type ProgramSources map[string]interface{}
