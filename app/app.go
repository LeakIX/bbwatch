package app

import (
	"github.com/LeakIX/bbwatch/conf"
	"github.com/LeakIX/bbwatch/program_sources"
	"github.com/LeakIX/bbwatch/vulnerability_source"
	"log"
)

type BBWatcher struct {
	Config              *conf.BBWatchConfig
	ProgramSources      []program_sources.ProgramSourceInterface
	VulnerabilitySource []vulnerability_source.VulnerabilitySourceInterface
}

func NewBBWatcher(config *conf.BBWatchConfig) *BBWatcher {
	bbwatcher := &BBWatcher{
		Config:              config,
		ProgramSources:      make([]program_sources.ProgramSourceInterface, 0),
		VulnerabilitySource: make([]vulnerability_source.VulnerabilitySourceInterface, 0),
	}
	bbwatcher.LoadConfig()
	return bbwatcher
}

func (bbw *BBWatcher) LoadConfig() {
	for vulnSourceName, vulnSourceConfig := range bbw.Config.VulnerabilitySources {
		if vulnSource, found := vulnerability_source.VulnerabilitySources[vulnSourceName]; found {
			vulnSource.Configure(vulnSourceConfig)
			bbw.VulnerabilitySource = append(bbw.VulnerabilitySource, vulnSource)
			log.Printf("loaded vulnerability source %s", vulnSourceName)
		}
	}
	for progSourceName, progSourceConfig := range bbw.Config.ProgramSources {
		if progSource, found := program_sources.ProgramSources[progSourceName]; found {
			progSource.Configure(progSourceConfig)
			bbw.ProgramSources = append(bbw.ProgramSources, progSource)
			log.Printf("loaded program source %s", progSourceName)
		}
	}
}

func (bbw *BBWatcher) Start() {
	for _, progSource := range bbw.ProgramSources {
		for program := range progSource.GetPrograms(true) {
			if len(program.Assets) < 1 {
				continue
			}
			for _, vulnSource := range bbw.VulnerabilitySource {
				for _, vuln := range vulnSource.GetVulnerabilities(program) {
					log.Printf("Found vulnerability for program %s on %s at %s : %s (%s)", program.Name, progSource.GetName(), vuln.Resource, vuln.Identifier, vuln.Severity)
				}
			}
		}
	}
}
