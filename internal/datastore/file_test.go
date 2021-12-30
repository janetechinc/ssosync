package datastore

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestFile(t *testing.T) {

	const (
		noSuchFileName     = "no_file.json"
		emptyFileName      = "empty.json"
		unreadableFileName = "unreadable.json"
		userFileName       = "Users.json"
		groupFileName      = "Groups.json"
	)

	setup := func(t *testing.T) string {
		prefix := t.TempDir() + "/"
		// create or insure absent files for the tests
		err := os.Remove(prefix + noSuchFileName)
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("could not remove file %s: %s", prefix+noSuchFileName, err)
		}
		f, err := os.Create(prefix + emptyFileName)
		if err != nil {
			t.Fatalf("could not create file %s: %s", prefix+emptyFileName, err)
		}
		f.Close()
		f, err = os.Create(prefix + unreadableFileName)
		if err != nil {
			t.Fatalf("could not create file %s: %s", prefix+unreadableFileName, err)
		}
		err = f.Chmod(0)
		if err != nil {
			t.Fatalf("could not set file mode on %s: %s", prefix+unreadableFileName, err)
		}
		f.Close()
		f, err = os.Create(prefix + userFileName)
		if err != nil {
			t.Fatalf("could not create valid user file %s: %s", prefix+userFileName, err)
		}
		_, err = f.WriteString("{\"user1@example.com\": true,\"user2@example.com\": true}")
		if err != nil {
			t.Fatalf("could not write valid user file %s: %s", prefix+userFileName, err)
		}
		f.Close()
		f, err = os.Create(prefix + groupFileName)
		if err != nil {
			t.Fatalf("could not create valid group file %s: %s", prefix+groupFileName, err)
		}
		_, err = f.WriteString("{\"group1\": true,\"group2\": true}")
		if err != nil {
			t.Fatalf("could not write valid user file %s: %s", prefix+groupFileName, err)
		}
		f.Close()
		return prefix
	}

	load := func(t *testing.T, ds Datastore, success bool) {
		if success {
			err := ds.Load()
			if err != nil {
				t.Errorf("%s", err)
			}
		} else {
			err := ds.Load()
			if err == nil {
				t.Errorf("should have failed")
			}
		}
	}

	store := func(t *testing.T, ds Datastore, success bool) {
		if success {
			err := ds.Load()
			if err != nil {
				t.Errorf("%s", err)
			}
		} else {
			err := ds.Load()
			if err == nil {
				t.Errorf("should have failed")
			}
		}
	}

	log.SetLevel(log.ErrorLevel)
	type TestType uint
	const (
		LoadOnly TestType = iota
		StoreOnly
		LoadStore
		LoadStoreLoad
	)
	tests := []struct {
		desc      string
		groupFile string
		userFile  string
		success   bool // true if success is expected
		test      TestType
	}{
		{
			desc:      "no existing file",
			userFile:  noSuchFileName,
			groupFile: noSuchFileName,
			success:   true,
			test:      LoadOnly,
		},
		{
			desc:      "empty user file",
			userFile:  emptyFileName,
			groupFile: noSuchFileName,
			success:   false,
			test:      LoadOnly,
		},
		{
			desc:      "empty group file",
			userFile:  noSuchFileName,
			groupFile: emptyFileName,
			success:   false,
			test:      LoadOnly,
		},
		{
			desc:      "unreadable user file",
			userFile:  unreadableFileName,
			groupFile: noSuchFileName,
			success:   false,
			test:      LoadStore,
		},
		{
			desc:      "unreadable group file",
			userFile:  noSuchFileName,
			groupFile: unreadableFileName,
			success:   false,
			test:      LoadStore,
		},
		{
			desc:      "valid user file",
			userFile:  userFileName,
			groupFile: noSuchFileName,
			success:   true,
			test:      LoadStoreLoad,
		},
		{
			desc:      "valid group file",
			userFile:  noSuchFileName,
			groupFile: groupFileName,
			success:   true,
			test:      LoadStoreLoad,
		},
	}

	for _, data := range tests {
		data := data
		prefix := setup(t)
		ds, err := NewFileDatastore(prefix, data.userFile, data.groupFile)
		if err != nil {
			t.Errorf("failed to create datastore for test '%s': %s", data.desc, err)
		}
		switch data.test {
		case LoadOnly:
			t.Run("Load "+data.desc, func(t *testing.T) { load(t, ds, data.success) })

		case StoreOnly:
			t.Run("Store "+data.desc, func(t *testing.T) { store(t, ds, data.success) })

		case LoadStore:
			t.Run("Load "+data.desc, func(t *testing.T) { load(t, ds, data.success) })
			t.Run("Store "+data.desc, func(t *testing.T) { store(t, ds, data.success) })

		case LoadStoreLoad:
			t.Run("Load "+data.desc, func(t *testing.T) { load(t, ds, data.success) })
			t.Run("Store "+data.desc, func(t *testing.T) { store(t, ds, data.success) })
			t.Run("Load "+data.desc, func(t *testing.T) { load(t, ds, data.success) })

		default:
			t.Run(data.desc, func(t *testing.T) { t.Errorf("unknown test type!") })
		}
	}
}
