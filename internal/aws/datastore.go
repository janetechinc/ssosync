package aws

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Datastore interface {
	Load(Client) error
	Store() error
	GetUsers() ([]*User, error)
	AddUser(*User) error
	DeleteUser(*User) error
	GetGroups() ([]*Group, error)
	AddGroup(*Group) error
	DeleteGroup(*Group) error
}

type datastoreUsers  map[string]*User
type datastoreGroups map[string]*Group
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

func (ds *baseDatastore) GetUsers() ([]*User, error) {
	v := make([]*User, 0, len(ds.users))
	for  _, value := range ds.users {
		v = append(v, value)
	 }
	 return v, nil
}

func (ds *baseDatastore) AddUser(user *User) error {
	log := log.WithFields(log.Fields{"user": user.Username})
	if _, ok := ds.users[user.Username]; !ok {
		log.Debug("adding user to user/group store")
		ds.users[user.Username] = user
	}
	return nil
}

func (ds *baseDatastore) DeleteUser(user *User) error {
	log := log.WithFields(log.Fields{"group": user.Username})
	log.Debug("deleting user from user/group store")
	delete(ds.users, user.Username)
	return nil
}

func (ds *baseDatastore) GetGroups() ([]*Group, error) {
	v := make([]*Group, 0, len(ds.groups))
	for  _, value := range ds.groups {
		v = append(v, value)
	 }
	 return v, nil
}

func (ds *baseDatastore) AddGroup(group *Group) error {
	log := log.WithFields(log.Fields{"group": group.DisplayName})
	if _, ok := ds.users[group.DisplayName]; !ok {
		log.Debug("adding group to user/group store")
		ds.groups[group.DisplayName] = group
	}
	return nil
}

func (ds *baseDatastore) DeleteGroup(group *Group) error {
	log := log.WithFields(log.Fields{"group": group.DisplayName})
	log.Debug("deleting group from user/group store")
	delete(ds.groups, group.DisplayName)
	return nil
}

