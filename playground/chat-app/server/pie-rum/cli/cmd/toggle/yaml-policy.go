package toggle

// omitempty => nil | value

type YamlISwap struct {
	Sub   bool   `yaml:"sub"`
	SubIn string `yaml:"sub_in,"`
}

type YamlPolicy struct {
	Name     string     `yaml:"name"`
	Activate bool       `yaml:"activate"`
	Deactive bool       `yaml:"deactive"`
	Delete   bool       `yaml:"delete"`
	Swap     *YamlISwap `yaml:"swap,"`
}

type YamlProfile struct {
	Policy YamlPolicy `yaml:"policy"`
}
type YamlKit struct {
	Policy YamlPolicy `yaml:"policy"`
}
type YamlService struct {
	Policy YamlPolicy `yaml:"policy"`
}

type YamlDispatcher struct {
	Policy YamlPolicy `yaml:"policy"`
}
