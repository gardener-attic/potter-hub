package kubeval

type SecretResponse struct {
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Data struct {
		Kubeconfig string `yaml:"kubeconfig"`
	} `yaml:"data"`
}
