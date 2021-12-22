package datastore

import (
	"encoding/json"

	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"

)

type consulDatastore struct {
	*baseDatastore
	kv       *consulapi.KV
	userKey  string
	groupKey string
}

func NewConsulDatastore(prefix string, userObj string, groupObj string) (Datastore, error) {
	consul, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &consulDatastore{
		baseDatastore: newBaseDatastore(),
		kv:       consul.KV(),
		userKey:  prefix+userObj,
		groupKey: prefix+groupObj,
	}, nil
}

func (ds *consulDatastore) Load() error {
	log.Info("Loading user/group lists from consul")
	log.Infof("loading users from '%s'", ds.userKey)

	pair, _, err := ds.kv.Get(ds.userKey, nil)
	if err != nil {
		log.Error("error fetching users:", err)
	} else if pair == nil {
		log.Warningf("consul KV '%s' does not exist:", ds.userKey, err)
	} else {
		err = json.Unmarshal(pair.Value, &ds.users)
		if err != nil {
			log.Error("failed to parse user list JSON from consul:", err)
		}
	}

	log.Infof("loading groups from '%s'", ds.groupKey)
	pair, _, err = ds.kv.Get(ds.groupKey, nil)
	if err != nil {
		log.Error("error fetching groups:", err)
	} else if pair == nil {
		log.Warningf("consul KV '%s'' does noty exist:", ds.groupKey, err)
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

