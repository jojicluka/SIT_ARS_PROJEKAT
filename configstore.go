package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sort"

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

	data, _, err := kv.List(constructKeyGroup(id, Version), nil)
	if err != nil || data == nil {
		return nil, errors.New("key not found")
	}

	entries := []map[string]string{}
	for _, pair := range data {
		group := &map[string]string{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *group)
	}

	labels := []map[string]string{}
	for _, pair := range data {
		group := &map[string]string{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		labels = append(labels, *group)
	}
	post := &Group{entries, labels, Version, id}

	return post, nil
}

func (ps *ConfigStore) Get(id string, Version string) (*Config, error) {
	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKeyConfig(id, Version), nil)
	if err != nil || pair == nil {
		return nil, errors.New("config not found")
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
	data, _, err := kv.List(constructKeyConfig(id, Version), nil)
	//fmt.Println(konf, " KONF")
	//fmt.Println(err)
	if err != nil || data == nil {
		return nil, errors.New("configuration does not exist")
	} else {
		_, fail := kv.Delete(constructKeyConfig(id, Version), nil)
		if fail != nil {
			return nil, fail
		}
		return map[string]string{"Deleted config with ID: ": id}, nil
	}
}

func (ps *ConfigStore) DeleteGroup(id string, Version string) (map[string]string, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(constructKeyGroup(id, Version), nil)
	if err != nil || data == nil {
		return nil, errors.New("group does not exist")
	} else {
		_, fail := kv.DeleteTree(constructKeyGroup(id, Version), nil)
		if fail != nil {
			return nil, fail
		}
		return map[string]string{"Deleted group: ": id}, err
	}
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

	groupID := uuid.New().String()

	for _, v := range post.Entries {
		label := ""
		stringList := []string{}
		for k, val := range v {
			stringList = append(stringList, k+":"+val)
		}
		sort.Strings(stringList)
		for _, v := range stringList {
			label += v + ";"
		}
		label = label[:len(label)-1]
		fmt.Println(label, " ?")
		sid := constructKeyGroupLabels(groupID, post.Version, label) + uuid.New().String()
		post.Id = groupID

		data, err := json.Marshal(post.Labels)
		if err != nil {
			return nil, err
		}

		p := &api.KVPair{Key: sid, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range post.Labels {
		entry := ""
		stringList := []string{}
		for k, val := range v {
			stringList = append(stringList, k+":"+val)
		}
		sort.Strings(stringList)
		for _, v := range stringList {
			entry += v + ";"
		}
		entry = entry[:len(entry)-1]
		fmt.Println(entry)
		sid := constructKeyGroupLabels(groupID, post.Version, entry) + uuid.New().String()
		post.Id = groupID

		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		p := &api.KVPair{Key: sid, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

func (ps *ConfigStore) FilterLabel(label string) (*[]Config, error) {
	kv := ps.cli.KV()

	pairs, _, err := kv.List("configs", nil)
	if err != nil {
		return nil, err
	}

	var retVal []Config

	for _, config := range pairs {
		c := &Config{}
		err = json.Unmarshal(config.Value, c)
		if err != nil {
			return nil, err
		}
		for _, l := range c.Labels {
			if l == label {
				retVal = append(retVal, *c)
			}
		}
	}

	return &retVal, err
}

func (ps *ConfigStore) GetGroupByLabel(id string, version string, label string) ([]map[string]string, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(constructKeyGroupLabels(id, version, label), nil)
	if err != nil {
		return nil, err
	}

	configs := []map[string]string{}
	for _, pair := range data {
		config := &map[string]string{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, *config)
	}

	return configs, nil
}
