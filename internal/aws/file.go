package aws

import (
	"encoding/json"
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"
)

type fileDatastore struct {
	*baseDatastore
	userFile string
	groupFile string
}

func NewFileDatastore(prefix string) (Datastore, error) {
	return &fileDatastore{
		baseDatastore: newBaseDatastore(),
		userFile: fmt.Sprintf("%sUsers.json", prefix),
		groupFile: fmt.Sprintf("%sGroups.json", prefix),
	}, nil
}

func addError(first error, second error) error {
	if first == nil {
		return second
	}
	return fmt.Errorf("%w; %s", first, second)
}

func (ds *fileDatastore) Load(client Client) (ret error) {
	log.Info("Loading user/group lists from files")

	log.Infof("loading users from '%s'", ds.userFile)
	uf, err := os.Open(ds.userFile)
	if err != nil {
		log.Errorf("failed to open %s file: %s", ds.userFile, err)
	} else {
		defer uf.Close();
		decoder := json.NewDecoder(uf)
		err = decoder.Decode(&ds.users)
		if err != nil {
			log.Errorf("failed to decode user list: %s", err)
			ret = addError(ret, err)
		}
	}

	log.Info("cleaning the just loaded user list")
	for name, _ := range ds.users {
		_, err := client.FindUserByEmail(name)
		if err == ErrUserNotFound {
			delete(ds.users, name)
			log.Debugf("VerifyUsers removed: %s", name)
			continue
		}
		if err != nil {
			log.Warningf("validate aws user failed for '%s' with: %s", name, err)
		}
	}

	log.Infof("loading groups from '%s'", ds.groupFile)
	gf, err := os.Open(ds.groupFile)
	if err != nil {
		log.Errorf("failed to open %s file: %s", ds.groupFile, err)
	} else {
		defer gf.Close();
		decoder := json.NewDecoder(gf)
		err = decoder.Decode(&ds.groups)
		if err != nil {
			log.Errorf("failed to decode group list: %s", err)
			ret = addError(ret, err)
		}
	}

	log.Info("cleaning the just loaded group list")
	for name, _ := range ds.groups {
		_, err := client.FindGroupByDisplayName(name)
		if err == ErrGroupNotFound {
			delete(ds.groups, name)
			log.Debugf("VeriftGroups removed: %s", name)
			continue
		}
		if err != nil {
			log.Warningf("validate aws group failed for '%s' with: %s", name, err)
		}
	}
	return
}

func (ds *fileDatastore) Store() error {
	log.Info("Storing user/group lists in files")

	log.Infof("storing users in '%s'", ds.userFile)
	uf, err := os.Create(ds.userFile)
	if err != nil {
		log.Errorf("failed to open %s for writing: %s", ds.userFile, err)
	} else {
		defer uf.Close();
		encoder := json.NewEncoder(uf)
		err = encoder.Encode(&ds.users)
		if err != nil {
			log.Error("failed to encode user list to json: %s", err)
		}
	}

	log.Infof("storing users in '%s'", ds.userFile)
	gf, err := os.Create(ds.userFile)
	if err != nil {
		log.Errorf("failed to open %s for writing: %s", ds.userFile, err)
	} else {
		defer gf.Close();
		encoder := json.NewEncoder(gf)
		err = encoder.Encode(&ds.groups)
		if err != nil {
			log.Error("failed to encode user list to json: %s", err)
		}
	}
	return err
}

