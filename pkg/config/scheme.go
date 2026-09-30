package config

import . "forge.pi.home.arpa/govalds/bot/pkg/marshal"

type Config struct {
	Token   String   `yaml:"token"`
	Logging *Logging `yaml:"logging"`
}

type Logging struct {
	Disabled    bool         `yaml:"disabled"`
	ClientTrace bool         `yaml:"client_trace"`
	Level       *String      `yaml:"level"`
	Timestamp   *bool        `yaml:"timestamp"`
	JSON        bool         `yaml:"json"`
	Colors      *bool        `yaml:"colors"`
	File        *LoggingFile `yaml:"file"`
}

type LoggingFile struct {
	Path        String `yaml:"path"`
	Placeholder String `yaml:"placeholder"`
}
