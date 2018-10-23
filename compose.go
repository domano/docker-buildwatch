package main

type ComposeFile struct {
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Image string `yaml:"image"`
	Build struct {
		Context string `yaml:"context"`
	} `yaml:"build"`
}
