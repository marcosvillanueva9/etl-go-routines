package model

type Yaml struct {
	Source struct {
		Connection struct {
			Path     string `yaml:"path"`
			Filetype string `yaml:"filetype"`
		} `yaml:"connection"`
		Columns struct {
			Origin      string `yaml:"origin"`
			Destination string `yaml:"destination"`
		} `yaml:"columns"`
	} `yaml:"source"`
	Destination struct {
		Connection struct {
			Path string `yaml:"path"`
		} `yaml:"connection"`
	} `yaml:"destination"`
	Transform []Connectors `yaml:"transform"`
}

type Connectors struct {
	Connector Connector `yaml:"connector"`
}

type Connector struct {
	Name   string `yaml:"name"`
	Params map[string]interface{} `yaml:"params"`
	Outcolumn struct {
		Destination string `yaml:"destination"`
	} `yaml:"outcolumn"`
}
