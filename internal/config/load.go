//go:build !wasm

package config

import (
	"errors"
	"os"
	"regexp"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
)

type configFile struct {
	path         string
	isGoReleaser bool
}

func Load(cmd *cobra.Command) (*Config, error) {
	conf := NewDefault()
	k := koanf.New(".")

	// Load default config
	if err := k.Load(structs.Provider(conf, "yaml"), nil); err != nil {
		return nil, err
	}

	// Find config file
	var cfgFiles []configFile
	cfgFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}
	if cfgFile != "" {
		cfgFiles = append(cfgFiles, configFile{path: cfgFile})
	} else {
		cfgFiles = []configFile{
			{path: ".changelog-generator.yaml"},
			{path: ".changelog-generator.yml"},
			{path: ".goreleaser.yaml", isGoReleaser: true},
			{path: ".goreleaser.yml", isGoReleaser: true},
		}
	}

	// Parse config file
	parser := yaml.Parser()
	for _, cfgFile := range cfgFiles {
		subK := k
		if cfgFile.isGoReleaser {
			subK = koanf.New(".")
		}

		if err := subK.Load(file.Provider(cfgFile.path), parser); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}

		if cfgFile.isGoReleaser {
			if changelogConf := subK.Get("changelog"); changelogConf != nil {
				if err := k.Load(confmap.Provider(changelogConf.(map[string]any), "."), nil); err != nil {
					return nil, err
				}
			}
		}

		break
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
	if len(conf.Groups) == 0 {
		conf.Groups = append(conf.Groups, &Group{})
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
