package datastore

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type fileDatastore struct {
	*baseDatastore
	userFile  string
	groupFile string
}

func NewFileDatastore(prefix string, userObj string, groupObj string) (Datastore, error) {
	return &fileDatastore{
		baseDatastore: newBaseDatastore(),
		userFile:      prefix + userObj,
		groupFile:     prefix + groupObj,
	}, nil
}

func (ds *fileDatastore) Load() error {
	log.Info("Loading user/group lists from files")

	log.Infof("loading users from '%s'", ds.userFile)
	uf, err := os.Open(ds.userFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warningf("failed to open %s file: %s", ds.userFile, err)
		} else {
			return fmt.Errorf("failed to open %s: %d", ds.userFile, err)
		}
	} else {
		defer uf.Close()
		decoder := json.NewDecoder(uf)
		err = decoder.Decode(&ds.users)
		if err != nil {
			return fmt.Errorf("failed to decode user list: %w", err)
		}
	}

	log.Infof("loading groups from '%s'", ds.groupFile)
	gf, err := os.Open(ds.groupFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warningf("failed to open %s file: %s", ds.groupFile, err)
		} else {
			return fmt.Errorf("failed to open %s: %d", ds.groupFile, err)
		}
	} else {
		defer gf.Close()
		decoder := json.NewDecoder(gf)
		err = decoder.Decode(&ds.groups)
		if err != nil {
			return fmt.Errorf("failed to decode group list: %w", err)
		}
	}

	return nil
}

func (ds *fileDatastore) Store() error {
	log.Info("Storing user/group lists in files")

	log.Infof("storing users in '%s'", ds.userFile)
	uf, err := os.Create(ds.userFile)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %w", ds.userFile, err)
	} else {
		defer uf.Close()
		encoder := json.NewEncoder(uf)
		encoder.SetIndent("", "    ")
		err = encoder.Encode(&ds.users)
		if err != nil {
			return fmt.Errorf("failed to encode user list to json: %w", err)
		}
	}

	log.Infof("storing groups in '%s'", ds.groupFile)
	gf, err := os.Create(ds.groupFile)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %w", ds.groupFile, err)
	} else {
		defer gf.Close()
		encoder := json.NewEncoder(gf)
		encoder.SetIndent("", "    ")
		err = encoder.Encode(&ds.groups)
		if err != nil {
			return fmt.Errorf("failed to encode group list to json: %w", err)
		}
	}
	return nil
}
