package datastore

import (
	"encoding/json"
	"fmt"

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
		kv:            consul.KV(),
		userKey:       prefix + userObj,
		groupKey:      prefix + groupObj,
	}, nil
}

func (ds *consulDatastore) Load() error {
	log.Info("Loading user/group lists from consul")
	log.Infof("loading users from '%s'", ds.userKey)

	pair, _, err := ds.kv.Get(ds.userKey, nil)
	if err != nil {
		return fmt.Errorf("error fetching users: %w", err)
	} else if pair == nil {
		log.Warningf("consul KV '%s' does not exist: %s", ds.userKey, err)
	} else {
		err = json.Unmarshal(pair.Value, &ds.users)
		if err != nil {
			return fmt.Errorf("failed to parse user list JSON from consul: %w", err)
		}
	}

	log.Infof("loading groups from '%s'", ds.groupKey)
	pair, _, err = ds.kv.Get(ds.groupKey, nil)
	if err != nil {
		return fmt.Errorf("error fetching groups: %w", err)
	} else if pair == nil {
		log.Warningf("consul KV '%s'' does noty exist: %s", ds.groupKey, err)
	} else {
		err = json.Unmarshal(pair.Value, &ds.groups)
		if err != nil {
			return fmt.Errorf("failed to parse group list JSON from consul: %w", err)
		}
	}

	return err
}

func (ds *consulDatastore) Store() error {
	log.Info("Storing user/group lists in consul")
	log.Infof("storing users to '%s'", ds.userKey)

	data, err := json.MarshalIndent(ds.users, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to convert user list to json: %w", err)
	}
	pair := consulapi.KVPair{
		Key:   ds.userKey,
		Value: data,
	}
	_, err = ds.kv.Put(&pair, nil)
	if err != nil {
		return fmt.Errorf("failed to PUT users in '%s': %w", ds.userKey, err)
	}
	data, err = json.MarshalIndent(ds.groups, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to convert group list to json: %w", err)
	}
	pair = consulapi.KVPair{
		Key:   ds.groupKey,
		Value: data,
	}
	_, err = ds.kv.Put(&pair, nil)
	if err != nil {
		return fmt.Errorf("failed to PUT groups in '%s': %w", ds.groupKey, err)
	}
	return nil
}
