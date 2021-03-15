package config

import (
	"fmt"
	"regexp"
)

var cfg *Config

const (
	releaseRegexFormat = `%s(\d+\.\d+\.\d+)`
	featureRegexFormat = `%s(.*)`
	hotfixRegexFormat  = `%s(.*)`
)

type Config struct {
	DefaultVersion string
	MainBranch     string
	ReleasePrefix  string
	FeaturePrefix  string
	HotfixPrefix   string
	RemoteName     string
	MainTag        string
	MainRegex      *regexp.Regexp
	ReleaseRegex   *regexp.Regexp
	FeatureRegex   *regexp.Regexp
	HotfixRegex    *regexp.Regexp
}

func LoadConfig() error {
	cfg = defaultConfig()
	cfg.ReleaseRegex = regexp.MustCompile(fmt.Sprintf(releaseRegexFormat, cfg.ReleasePrefix))
	cfg.FeatureRegex = regexp.MustCompile(fmt.Sprintf(featureRegexFormat, cfg.FeaturePrefix))
	cfg.HotfixRegex = regexp.MustCompile(fmt.Sprintf(hotfixRegexFormat, cfg.HotfixPrefix))
	cfg.MainRegex = regexp.MustCompile(cfg.MainBranch)
	return nil
}

func defaultConfig() *Config {
	return &Config{
		DefaultVersion: "0.1.0",
		MainBranch:     "^(main|master)$",
		ReleasePrefix:  "release/",
		FeaturePrefix:  "feature/",
		HotfixPrefix:   "fix/",
		RemoteName:     "origin",
		MainTag:        "beta",
	}
}

func GetConfig() *Config {
	return cfg
}
