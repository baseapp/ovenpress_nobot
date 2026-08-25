package settings

import (
	"git.gammaspectra.live/git/go-away/utils"
	"maps"
	"strings"
)

type Settings struct {
	Bind Bind `yaml:"bind"`

	Backends map[string]Backend `yaml:"backends"`

	BindDebug   string `yaml:"bind-debug"`
	BindMetrics string `yaml:"bind-metrics"`

	Strings utils.Strings `yaml:"strings"`

	// Links to add to challenge/error pages like privacy/impressum.
	Links []Link `yaml:"links"`

	ChallengeTemplate string `yaml:"challenge-template"`

	// ChallengeTemplateOverrides Key/Value overrides for the current chosen template
	ChallengeTemplateOverrides map[string]string `yaml:"challenge-template-overrides"`

	ClientIpHeader  string `yaml:"client-ip-header"`
	BackendIpHeader string `yaml:"backend-ip-header"`
	AccessLog       string `yaml:"access-log"`
}

type Link struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

var DefaultSettings = Settings{
	Strings:           DefaultStrings,
	ChallengeTemplate: "anubis",
	ChallengeTemplateOverrides: func() map[string]string {
		m := make(map[string]string)
		maps.Copy(m, map[string]string{
			"Theme": "",
			"Logo":  "",
		})
		return m
	}(),

	Bind: Bind{
		Address:         ":8080",
		Network:         "tcp",
		SocketMode:      "0770",
		Proxy:           false,
		TLSAcmeAutoCert: "",
	},
	Backends: make(map[string]Backend),
}

func (s *Settings) SelectBackend(host string) (Backend, bool) {
	backend, ok := s.Backends[host]
	if !ok {
		// do wildcard match
		parts := strings.Split(host, ".")
		if len(parts) > 1 {
			wildcard := "*." + strings.Join(parts[1:], ".")
			backend, ok = s.Backends[wildcard]
		}

		if !ok {
			// return fallback
			backend, ok = s.Backends["*"]
		}
	}
	return backend, ok
}
