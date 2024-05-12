package config

type Provider interface {
	Name() string
	ReadConfigSectionDefinition(name string, v interface{}) error
	Parse() error
}

type ProviderSetup interface {
	AddProvider(provider Provider) error
	BuildDefinitions() (DefinitionSetup, error)
}

type DefinitionSetup interface {
	AddConfigSectionDefinition(name string, v interface{}) error
	BuildConfiguration() (Configuration, error)
}

type Configuration interface {
	Parse() error
	Parsed() bool
}
