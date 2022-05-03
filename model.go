package main

type Config struct {
	entries map[string]string
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
}
