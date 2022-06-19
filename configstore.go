package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//	"github.com/google/uuid"
	"os"
	"reflect"
	//	"sort"
	"strings"

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

	sid := constructKeyGroup(id, Version)
	pair, _, err := kv.Get(sid, nil)
	if pair == nil {
		return nil, errors.New("Could not find a group.")
	}
	configgroup := &Group{}
	err = json.Unmarshal(pair.Value, configgroup)
	if err != nil {
		return nil, err
	}
	return configgroup, nil
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

	sid, rid := generateKeyGroup(post.Version)
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	pairs, _, err := kv.Get(sid, nil)
	if pairs != nil {
		return nil, errors.New("group already exists. ")
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return post, nil
}

//func (ps *ConfigStore) FilterLabel(label string) (*[]Config, error) {
//	kv := ps.cli.KV()
//
//	pairs, _, err := kv.List("configs", nil)
//	if err != nil {
//		return nil, err
//	}
//
//	var retVal []Config
//
//	for _, config := range pairs {
//		c := &Config{}
//		err = json.Unmarshal(config.Value, c)
//		if err != nil {
//			return nil, err
//		}
//		for _, l := range c.Labels {
//			if l == label {
//				retVal = append(retVal, *c)
//			}
//		}
//	}
//
//	return &retVal, err
//}

func (ps *ConfigStore) GetGroupByLabel(id string, version string, label string) ([]*GroupConfig, error) {
	kv := ps.cli.KV()
	configList := []*GroupConfig{}
	sid := constructKeyGroup(id, version)

	data, _, err := kv.Get(sid, nil)
	if data == nil {
		return nil, err
	}

	labelList := strings.Split(label, ";")
	configStoreLabels := make(map[string]string)
	for _, label := range labelList {
		part := strings.Split(label, ":")
		if part != nil {
			configStoreLabels[part[0]] = part[1]
		}
	}

	group := &Group{}
	err = json.Unmarshal(data.Value, group)
	if err != nil {
		return nil, err
	}

	for _, config := range group.Configs {
		if len(config.Entries) == len(configStoreLabels) {
			if reflect.DeepEqual(config.Entries, configStoreLabels) {
				configList = append(configList, config)
			}
		}
	}
	//configs := []map[string]string{}
	//for _, pair := range data {
	//	config := &map[string]string{}
	//	err = json.Unmarshal(pair.Value, config)
	//	if err != nil {
	//		return nil, err
	//	}
	//	configs = append(configs, *config)
	//}

	return configList, nil
}

func (ps *ConfigStore) UpdateConfigGroup(group *Group) (*Group, error) {
	kv := ps.cli.KV()
	data, err := json.Marshal(group)

	sid := constructKeyGroup(group.Id, group.Version)
	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}
