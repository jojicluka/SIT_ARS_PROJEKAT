package main

type Config struct {
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
	//Labels  map[string]string `json:"labels"`
	Id string `json:"Id"`
}

type Group struct {
	Configs []*GroupConfig `json:"configs"`
	Version string         `json:"version"`
	Id      string         `json:"Id"`
}

type GroupConfig struct {
	Entries map[string]string `json:"entries"`
	Values  map[string]string `json:"values"`
}
