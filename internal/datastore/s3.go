package datastore

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

type s3Datastore struct {
	*baseDatastore
	s3       *s3.S3
	bucket   string
	userKey  string
	groupKey string
}

func NewS3Datastore(bucket string, userObj string, groupObj string) (Datastore, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &s3Datastore{
		baseDatastore: newBaseDatastore(),
		s3:       s3.New(sess),
		bucket:   bucket,
		userKey:  userObj,
		groupKey: groupObj,
	}, nil
}

func (ds *s3Datastore) Load() (ret error) {
	log.Info("Loading user/group lists from S3")
	log.Infof("loading users from bucket '%s' object '%s'", ds.bucket, ds.userKey)
	userResults, err := ds.s3.GetObject(&s3.GetObjectInput{
        Bucket: aws.String(ds.bucket),
        Key:    aws.String(ds.userKey),
    })
	if err != nil {
		// cast to awserr err
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				log.Warningf("S3 key '%s' does not exist:", ds.userKey, err)
			} else {
				log.Error("error fetching users:", err)
			}
		}
	} else {
		defer userResults.Body.Close()
		decoder := json.NewDecoder(userResults.Body)
		err = decoder.Decode(&ds.users)
		if err != nil {
			log.Errorf("failed to decode user list: %s", err)
			ret = addError(ret, err)
		}
	}

	log.Infof("loading groups from bucket '%s' object '%s'", ds.bucket, ds.groupKey)
	groupResult, err := ds.s3.GetObject(&s3.GetObjectInput{
        Bucket: aws.String(ds.bucket),
        Key:    aws.String(ds.groupKey),
    })
	if err != nil {
		// cast to awserr err to determin if its that the key does not exist
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				log.Warningf("S3 key '%s' does not exist:", ds.groupKey, err)
			} else {
				log.Error("error fetching users:", err)
			}
		}
	} else {
		defer groupResult.Body.Close()
		decoder := json.NewDecoder(groupResult.Body)
		err = decoder.Decode(&ds.groups)
		if err != nil {
			log.Errorf("failed to decode group list: %s", err)
			ret = addError(ret, err)
		}
	}
	return
}

func (ds *s3Datastore) Store() error {
	log.Infof("Storing user/group lists in S3 bucket: %s", ds.bucket)
	data, err := json.Marshal(ds.users)
    if err != nil {
		log.Error("failed to convert user list to json", err)
        return err
    }
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
        Bucket: aws.String(ds.bucket),
        Key:    aws.String(ds.userKey),
    }
	_, err = ds.s3.PutObject(input)
	if err != nil {
		log.Error("failed to PUT user list in S3: ", err)
        return err
	}

	data, err = json.Marshal(ds.groups)
    if err != nil {
		log.Error("failed to convert user list to json", err)
        return err
    }
	input = &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
        Bucket: aws.String(ds.bucket),
        Key:    aws.String(ds.groupKey),
    }
	_, err = ds.s3.PutObject(input)
	if err != nil {
		log.Error("failed to PUT group list in S3: ", err)
        return err
	}

	return err
}

