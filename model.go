package main

type Config struct {
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
	Id      string            `json:"Id"`
}

type Group struct {
	Entries []map[string]string `json:"entries"`
	Version string              `json:"version"`
	Id      string              `json:"Id"`
}
