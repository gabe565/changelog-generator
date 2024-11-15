//go:build !wasm

package config

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"gabe565.com/utils/must"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
)

func Load(cmd *cobra.Command) (*Config, error) {
	k := koanf.New(".")
	conf := NewDefault()

	// Load default config
	if err := k.Load(structs.Provider(conf, "yaml"), nil); err != nil {
		return nil, err
	}

	// Find config file
	cfgFiles := make([]string, 0, 4)
	cfgFile := must.Must2(cmd.Flags().GetString(ConfigFlag))
	if cfgFile != "" {
		cfgFiles = append(cfgFiles, cfgFile)
	} else {
		repoPath := must.Must2(cmd.Flags().GetString(RepoFlag))
		cfgFiles = append(cfgFiles,
			filepath.Join(repoPath, ".changelog-generator.yaml"),
			filepath.Join(repoPath, ".changelog-generator.yml"),
			filepath.Join(repoPath, ".github", "changelog-generator.yaml"),
			filepath.Join(repoPath, ".github", "changelog-generator.yml"),
			filepath.Join(repoPath, ".goreleaser.yaml"),
			filepath.Join(repoPath, ".goreleaser.yml"),
		)
	}

	// Parse config file
	parser := yaml.Parser()
	for _, cfgFile := range cfgFiles {
		subK := k
		if isGoReleaser(cfgFile) {
			subK = koanf.New(".")
		}

		if err := subK.Load(file.Provider(cfgFile), parser); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}

		if isGoReleaser(cfgFile) {
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
		if g.Regexp != "" {
			re, err := regexp.Compile(g.Regexp)
			if err != nil {
				return nil, err
			}
			g.re = re
		}
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

	if conf.Tag.Regexp != "" {
		conf.Tag.Regexp = "^" + conf.Tag.Regexp + "$"
		var err error
		if conf.Tag.re, err = regexp.Compile(conf.Tag.Regexp); err != nil {
			return nil, err
		}
	}

	return conf, nil
}

func isGoReleaser(path string) bool {
	base := filepath.Base(path)
	return base == ".goreleaser.yaml" || base == ".goreleaser.yml"
}
