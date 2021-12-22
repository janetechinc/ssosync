package datastore

import (
	"encoding/json"

	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"

)

type consulDatastore struct {
	*baseDatastore
	kv        *consulapi.KV
}

func NewConsulDatastore() (Datastore, error) {
	consul, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &consulDatastore{
		baseDatastore: newBaseDatastore(),
		kv:        consul.KV(),
	}, nil
}

func (ds *consulDatastore) Load() error {
	log.Info("Loading user/group lists from consul")
	log.Info("loading users from 'aws-ssosync/users'")

	pair, _, err := ds.kv.Get("aws-ssosync/users", nil)
	if err != nil {
		log.Error("error fetching consul KV aws-ssosync/users!", err)
	} else if pair == nil {
		log.Error("consul KV aws-ssosync/users does not exist!", err)
	} else {
		err = json.Unmarshal(pair.Value, &ds.users)
		if err != nil {
			log.Error("failed to parse user list JSON from consul", err)
		}
	}

	log.Info("loading groups from 'aws-ssosync/groups'")
	pair, _, err = ds.kv.Get("aws-ssosync/groups", nil)
	if err != nil {
		log.Error("error fetching consul KV aws-ssosync/users!", err)
	} else if pair == nil {
		log.Warn("consul KV aws-ssosync/groups does noty exist!", err)
	} else {
		err = json.Unmarshal(pair.Value, &ds.groups)
		if err != nil {
			log.Error("failed to parse group list JSON from consul", err)
		}
	}

	return err
}

func (ds *consulDatastore) Store() error {
	log.Info("Storing user/group lists in consul")
	data, err := json.MarshalIndent(ds.users, "", "    ")
	if err != nil {
		log.Error("failed to convert user list to json", err)
		return err
	}
	pair := consulapi.KVPair{
		Key: "aws-ssosync/users",
		Value: data,
	}
	_, err = ds.kv.Put(&pair, nil)
	if err != nil {
		log.Error("failed to PUT consul KV for key 'aws-ssosync/users'")
	}
	data, err = json.MarshalIndent(ds.groups, "", "    ")
	if err != nil {
		log.Error("failed to convert group list to json", err)
		return err
	}
	pair = consulapi.KVPair{
		Key: "aws-ssosync/groups",
		Value: data,
	}
	_, err = ds.kv.Put(&pair, nil)
	if err != nil {
		log.Error("failed to PUT consul KV for key 'aws-ssosync/groups'")
	}
	return err
}

