package validator

type Config struct {
	Dns struct {
		Allow []DnsRule `yaml:"allow"`
		Deny  []DnsRule `yaml:"deny"`
	} `yaml:"dns"`
}
