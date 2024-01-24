package app

import (
	"fmt"
	"github.com/LeakIX/bbwatch/conf"
	"github.com/LeakIX/bbwatch/program_sources"
	"github.com/LeakIX/bbwatch/vulnerability_source"
	"github.com/gookit/color"
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
			log.Printf("Loaded vulnerability source %s", vulnSource.GetName())
		}
	}
	for progSourceName, progSourceConfig := range bbw.Config.ProgramSources {
		if progSource, found := program_sources.ProgramSources[progSourceName]; found {
			progSource.Configure(progSourceConfig)
			bbw.ProgramSources = append(bbw.ProgramSources, progSource)
			log.Printf("Loaded program source %s", progSource.GetName())
		}
	}
}

func (bbw *BBWatcher) Start() {
	for _, progSource := range bbw.ProgramSources {
		log.Printf("Processing programs from %s ...", progSource.GetName())
		for program := range progSource.GetPrograms(bbw.Config.BountyOnly) {
			if len(program.Assets) < 1 {
				continue
			}
			for _, vulnSource := range bbw.VulnerabilitySource {
				for _, vuln := range vulnSource.GetVulnerabilities(program) {
					fmt.Printf("[%s -> %s] [%s] Found vulnerability for program %s at %s -> %s \n",
						progSource.GetName(), vulnSource.GetName(), bbw.ColorSeverity(vuln.Severity), color.FgLightWhite.Render(program.Name), color.FgLightWhite.Render(vuln.Resource), color.FgLightWhite.Render(vuln.Identifier))
				}
			}
		}
	}
}

func (bbw *BBWatcher) ColorSeverity(severity string) string {
	switch severity {
	case "critical":
		return color.FgLightRed.Render(severity)
	case "high":
		return color.FgYellow.Render(severity)
	case "medium":
		return color.FgLightBlue.Render(severity)
	case "low":
		return color.FgLightWhite.Render(severity)
	}
	return color.FgLightWhite.Render(severity)
}
