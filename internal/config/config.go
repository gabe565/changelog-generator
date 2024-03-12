package config

var Default *Config

type Config struct {
	Sort    string   `yaml:"sort"`
	Abbrev  int      `yaml:"abbrev"`
	Groups  []*Group `yaml:"groups"`
	Divider string   `yaml:"divider"`
	Filters Filters  `yaml:"filters"`
}

func NewDefault() *Config {
	return &Config{
		Sort:   SortAscending,
		Abbrev: 8,
	}
}
