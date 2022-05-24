package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type ConfigStore struct {
	cli *api.Client
}

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}

func (ps *ConfigStore) GetGroup(id string, Version string) (*Group, error) {
	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKeyGroup(id, Version), nil)
	if err != nil {
		return nil, err
	}

	post := &Group{}
	err = json.Unmarshal(pair.Value, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *ConfigStore) Get(id string, Version string) (*Config, error) {
	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKeyConfig(id, Version), nil)
	if err != nil {
		return nil, err
	}

	post := &Config{}
	err = json.Unmarshal(pair.Value, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *ConfigStore) Delete(id string, Version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(constructKeyConfig(id, Version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) DeleteGroup(id string, Version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(constructKeyGroup(id, Version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *ConfigStore) Post(post *Config) (*Config, error) {
	kv := ps.cli.KV()

	sid, rid := generateKeyConfig(post.Version)
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *ConfigStore) PostGroup(post *Group) (*Group, error) {
	kv := ps.cli.KV()

	sid, rid := generateKeyGroup(post.Version)
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return post, nil
}
