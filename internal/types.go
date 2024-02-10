package internal

type Matter struct {
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
}

type Script struct {
	Name   string
	Matter Matter
}
