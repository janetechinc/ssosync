package datastore

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

//
// This test requires s3 permissions
//
// then run the tests with these exported:
// export CONSUL_HTTP_ADDR=localhost:8500
// export CONSUL_TEST_PREFIX=ssosync_datastore_test
//
// required AWS auth env vars for this test to run:
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY
// AWS_REGION
//
// plus these needed for the test:
// DATASTORE_S3_BUCKET - aws bucket to use
// DATASTORE_S3_FOLDER - folder in bucket
//
// permissions required for bucket/folder are:
// PutObject, DeleteObject, ListBucket, GetObject
//
func TestS3(t *testing.T) {
	reqEnvs := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
		"DATASTORE_S3_BUCKET",
		"DATASTORE_S3_FOLDER",
	}
	for _, env := range reqEnvs {
		if _, ok := os.LookupEnv(env); !ok {
			t.Skipf("%s is not set", env)
		}
	}

	const (
		noSuchFileName  = "no_file.json"
		invalidFileName = "invalid.json"
		userFileName    = "Users.json"
		groupFileName   = "Groups.json"
	)

	sess, err := session.NewSession()
	if err != nil {
		t.Errorf("can't connect to consul: %s", err)
	}

	s3client := s3.New(sess)

	prefixCount := 0
	bucket := os.Getenv("DATASTORE_S3_BUCKET")
	folder := os.Getenv("DATASTORE_S3_FOLDER")
    if !strings.HasSuffix(folder, "/") {
        folder = folder+"/"
    }

	put := func(key string, value string) {
		input := &s3.PutObjectInput{
			Body:   aws.ReadSeekCloser(strings.NewReader(value)),
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}
		_, err = s3client.PutObject(input)
		if err != nil {
			t.Errorf("failed PUT key %s: %s", key, err)
		}
	}
	setup := func() string {
		prefix := fmt.Sprintf("%s%d-", folder, prefixCount)
		prefixCount += 1
		// create or insure absent files for the tests
		input := &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(prefix + noSuchFileName),
		}
		_, err := s3client.DeleteObject(input)
		if err != nil {
			t.Fatalf("could not remove key %s: %s", prefix+noSuchFileName, err)
		}
		put(prefix+invalidFileName, "[x]")
		put(prefix+userFileName, "{\"user1@example.com\": true,\"user2@example.com\": true}")
		put(prefix+groupFileName, "{\"group1\": true,\"group2\": true}")
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
		ds, err := NewS3Datastore(bucket, prefix+data.userFile, prefix+data.groupFile)
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
