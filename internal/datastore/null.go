package datastore

import ()

type nullDatastore struct {
	*baseDatastore
}

func NewNullDatastore() Datastore {
	return &nullDatastore{
		baseDatastore: newBaseDatastore(),
	}
}

func (ds *nullDatastore) Load() error {
	return nil
}

func (ds *nullDatastore) Store() error {
	return nil
}
