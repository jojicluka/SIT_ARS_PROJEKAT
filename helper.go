package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	configs = "configs/%s/%s"
	groups  = "groups/%s/%s"
)

func generateKeyConfig(Version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(configs, id, Version), id
}

func constructKeyConfig(id string, Version string) string {
	return fmt.Sprintf(configs, id, Version)
}

func generateKeyGroup(Version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(groups, id, Version), id
}

func constructKeyGroup(id string, Version string) string {
	return fmt.Sprintf(groups, id, Version)
}

func decodeBody(r io.Reader) ([]*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt []*Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func decodeConfigBody(r io.Reader) (*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt *Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func decodeGroupBody(r io.Reader) (*Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt *Group
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}
