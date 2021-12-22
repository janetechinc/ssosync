package datastore

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Datastore interface {
	Load() error
	Store() error
	GetUsers() ([]string, error)
	AddUser(string) error
	DeleteUser(string) error
	GetGroups() ([]string, error)
	AddGroup(string) error
	DeleteGroup(string) error
}

type datastoreUsers  map[string]bool
type datastoreGroups map[string]bool
type baseDatastore struct {
	users  datastoreUsers
	groups datastoreGroups
}

func newBaseDatastore() (*baseDatastore) {
	return &baseDatastore{
		users: datastoreUsers{},
		groups: datastoreGroups{},
	}
}

func NewDatastore() (Datastore, error) {
	return nil, fmt.Errorf("no Datastore implementstions available!")
}

func (ds *baseDatastore) GetUsers() ([]string, error) {
    users := make([]string, 0, len(ds.users))
    for name, _ := range ds.users {
        users = append(users, name)
    }
	return users, nil
/*
	v := make([]*User, 0, len(ds.users))	
	for  name, _ := range ds.users {
		user, err := client.FindUserByEmail(name)
		if err == ErrUserNotFound {
			delete(ds.users, name)
			log.Debugf("GetUsers removed non-existant user from list: %s", name)
			continue
		}
		if err != nil {
			log.Fatalf("GetUsers failed for '%s' with: %s", name, err)
		}
		v = append(v, user)
	 }
	 return v, nil
*/
}

func (ds *baseDatastore) AddUser(user string) error {
	log := log.WithFields(log.Fields{"user": user})
	if _, ok := ds.users[user]; !ok {
		log.Debug("adding user to datastore")
		ds.users[user] = true
	}
	return nil
}

func (ds *baseDatastore) DeleteUser(user string) error {
	log := log.WithFields(log.Fields{"group": user})
	log.Debug("deleting user from datastore")
	delete(ds.users, user)
	return nil
}

func (ds *baseDatastore) GetGroups() ([]string, error) {
    groups := make([]string, 0, len(ds.groups))
    for name, _ := range ds.groups {
        groups = append(groups, name)
    }
	return groups, nil
/*
	v := make([]*Group, 0, len(ds.groups))
	for  name, _ := range ds.groups {
		group, err := client.FindGroupByDisplayName(name)
		if err == ErrGroupNotFound {
			delete(ds.groups, name)
			log.Debugf("GetGroups removed non-existant group from list: %s", name)
			continue
		}
		if err != nil {
			log.Fatalf("GetGroups failed for '%s' with: %s", name, err)
		}
		v = append(v, group)
	 }
	 return v, nil
*/
}

func (ds *baseDatastore) AddGroup(group string) error {
	log := log.WithFields(log.Fields{"group": group})
	if _, ok := ds.groups[group]; !ok {
		log.Debug("adding group to datastore")
		ds.groups[group] = true
	}
	return nil
}

func (ds *baseDatastore) DeleteGroup(group string) error {
	log := log.WithFields(log.Fields{"group": group})
	log.Debug("deleting group from datastore")
	delete(ds.groups, group)
	return nil
}

