package main

import (
	"fmt"
	"regexp"
)

var cfg *Config

const (
	releaseRegexFormat = `%s(\d+\.\d+\.\d+)`
	featureRegexFormat = `%s(.*)`
)

type Config struct {
	DefaultVersion string
	MainBranch     string
	ReleasePrefix  string
	FeaturePrefix  string
	HotfixPrefix   string
	RemoteName     string
	MainTag        string
	ReleaseRegex   *regexp.Regexp
	FeatureRegex   *regexp.Regexp
}

func LoadConfig() error {
	cfg = defaultConfig()
	cfg.ReleaseRegex = regexp.MustCompile(fmt.Sprintf(releaseRegexFormat, cfg.ReleasePrefix))
	cfg.FeatureRegex = regexp.MustCompile(fmt.Sprintf(featureRegexFormat, cfg.FeaturePrefix))
	return nil
}

func defaultConfig() *Config {
	return &Config{
		DefaultVersion: "0.1.0",
		MainBranch:     "main",
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
