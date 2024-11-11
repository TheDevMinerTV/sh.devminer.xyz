package internal

type Matter struct {
	Description  string   `yaml:"description"`
	Tags         []string `yaml:"tags"`
	DownloadOnly bool     `yaml:"download-only"`
}

type Script struct {
	Name   string
	Matter Matter
}
