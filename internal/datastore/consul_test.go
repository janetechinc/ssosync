package datastore

import (
	"fmt"
	"os"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

//
// Use a local consul in dev mode
// Run with `consul agent -dev`
//
// then run the tests with these exported:
// export CONSUL_HTTP_ADDR=localhost:8500
// export CONSUL_TEST_PREFIX=ssosync_datastore_test
// 
//
func TestConsul(t *testing.T) {
	// skip if CONSUL_HTTP_ADDR is not set
	if _, ok := os.LookupEnv("CONSUL_HTTP_ADDR"); !ok {
		t.Skip("CONSUL_HTTP_ADDR is not set")
	}
	if _, ok := os.LookupEnv("CONSUL_TEST_PREFIX"); !ok {
		t.Skip("CONSUL_TEST_PREFIX is not set")
	}

    const (
        noSuchFileName     = "no_file.json"
        invalidFileName    = "invalid.json"
        userFileName       = "Users.json"
        groupFileName      = "Groups.json"
    )

	consul, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		t.Errorf("can't connect to consul: %s", err)
	}

	kv := consul.KV()

	prefixCount := 0

	put := func(key string, value string) {
		pair := consulapi.KVPair{
			Key:   key,
			Value: []byte(value),
		}
		_, err = kv.Put(&pair, nil)
		if err != nil {
			t.Errorf("failed PUT key %s: %s", key, err)
		}
	}
	setup := func() string {
		prefix := fmt.Sprintf("%s%d/", os.Getenv("CONSUL_TEST_PREFIX"), prefixCount);
		prefixCount += 1;
		// create or insure absent files for the tests
	
        _, err := kv.Delete(prefix + noSuchFileName, nil)
        if err != nil {
            t.Fatalf("could not remove file %s: %s", prefix+noSuchFileName, err)
        }
		put(prefix + invalidFileName, "[x]")
		put(prefix + userFileName, "{\"user1@example.com\": true,\"user2@example.com\": true}")
		put(prefix + groupFileName, "{\"group1\": true,\"group2\": true}")
        return prefix
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
			desc:      "invalid user file",
			userFile:  invalidFileName,
			groupFile: noSuchFileName,
			success:   false,
			test:      LoadOnly,
		},
		{
			desc:      "invalid group file",
			userFile:  noSuchFileName,
			groupFile: invalidFileName,
			success:   false,
			test:      LoadOnly,
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
		prefix := setup()
		ds, err := NewConsulDatastore(prefix, data.userFile, data.groupFile)
		if err != nil {
			t.Errorf("failed to create datastore for test '%s': %s", data.desc, err)
		}

		runLoad := func() {
			if data.success {
				err := ds.Load()
				if err != nil {
					t.Errorf("%s - %s", data.desc, err)
				}
			} else {
				err := ds.Load()
				if err == nil {
					t.Errorf("%s - should have failed", data.desc)
				}
			}
		}
		
		runStore := func() {
			if data.success {
				err := ds.Load()
				if err != nil {
					t.Errorf("%s - %s", data.desc, err)
				}
			} else {
				err := ds.Load()
				if err == nil {
					t.Errorf("%s - should have failed", data.desc)
				}
			}
		}
	
		switch data.test {
		case LoadOnly:
			t.Run("Load "+data.desc, func(t *testing.T) { runLoad() })

		case StoreOnly:
			t.Run("Store "+data.desc, func(t *testing.T) { runStore() })

		case LoadStore:
			t.Run("Load "+data.desc, func(t *testing.T) { runLoad() })
			t.Run("Store "+data.desc, func(t *testing.T) { runStore() })

		case LoadStoreLoad:
			t.Run("Load "+data.desc, func(t *testing.T) { runLoad() })
			t.Run("Store "+data.desc, func(t *testing.T) { runStore() })
			t.Run("Load "+data.desc, func(t *testing.T) { runLoad() })

		default:
			t.Run(data.desc, func(t *testing.T) { t.Errorf("unknown test type!") })
		}
	}
}
