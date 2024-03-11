//go:build !wasm

package config

import (
	"regexp"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
)

func Load(cmd *cobra.Command) (*Config, error) {
	conf := NewDefault()
	k := koanf.New(".")

	// Load default config
	if err := k.Load(structs.Provider(conf, "yaml"), nil); err != nil {
		return nil, err
	}

	// Find config file
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}
	if cfgFile == "" {
		cfgFile = ".changelog-generator.yaml"
	}

	// Parse config file
	parser := yaml.Parser()
	if err := k.Load(file.Provider(cfgFile), parser); err != nil {
		return nil, err
	}

	if err := k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return nil, err
	}

	for _, g := range conf.Groups {
		re, err := regexp.Compile(g.Regexp)
		if err != nil {
			return nil, err
		}
		g.re = re
	}

	for _, exclude := range conf.Filters.Exclude {
		re, err := regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
		conf.Filters.excludeRe = append(conf.Filters.excludeRe, re)
	}
	for _, exclude := range conf.Filters.Include {
		re, err := regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
		conf.Filters.includeRe = append(conf.Filters.includeRe, re)
	}

	Default = conf
	return conf, err
}
