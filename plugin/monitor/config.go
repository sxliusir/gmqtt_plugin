package monitor

// Config is the configuration for the monitor plugin.
type Config struct {
	// add your config fields
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Db string `yaml:"db"`
	Charset string `yaml:"charset"`
}

// Validate validates the configuration, and return an error if it is invalid.
func (c *Config) Validate() error {
	//panic("implement me")
	return nil
}

// DefaultConfig is the default configuration.
var DefaultConfig = Config{
	Host: "127.0.0.1",
	Port: "3306",
	User: "root",
	Password: "root",
	Db: "demo",
	Charset: "utf8",
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type cfg Config
	var v = &struct {
		Monitor cfg `yaml:"monitor"`
	}{
		Monitor: cfg(DefaultConfig),
	}
	if err := unmarshal(v); err != nil {
		return err
	}
	empty := cfg(Config{})
	if v.Monitor == empty {
		v.Monitor = cfg(DefaultConfig)
	}
	*c = Config(v.Monitor)
	return nil
}
