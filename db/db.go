package db

import (
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	cn "virtualizer/constants"
	u "virtualizer/utils"
)

func connect() (session *mgo.Session, err error) {

	connectURL := os.Getenv("MONGO_DB_URL")
	session, err = mgo.Dial(connectURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Panic("Can't connect to mongo, go error: " + err.Error() + "\n")
		return nil, err
	}

	session.SetSafe(&mgo.Safe{})
	return session, nil
}

func GetData(db string, collection string, condition []bson.M) (result string, err error) {

	session, err := connect()
	if err != nil {
		return "", err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(db).C(collection)

	var m bson.M
	var all []bson.M
	var bytes []byte
	if condition != nil {
		err = c.Find(bson.M{cn.MONGO_DB_AND_OPERATOR: condition}).One(&m)
		if err != nil {
			return "", err
		}

		// remove id
		resp := make(map[string]interface{})
		for k, v := range m {
			if k != cn.MONGO_DB_FIELD_ID {
				resp[k] = v
			}
		}

		bytes, err = json.Marshal(resp)
		if err != nil {
			return "", err
		}

	} else {
		err = c.Find(nil).All(&all)
		if err != nil {
			return "", err
		}

		var resp []map[string]interface{}
		for _, obj := range all {
			delete(obj, cn.MONGO_DB_FIELD_ID)
			resp = append(resp, obj)
		}

		bytes, err = json.Marshal(resp)
		if err != nil {
			return "", err
		}
	}

	result = u.BytesToString(bytes)
	logrus.WithFields(logrus.Fields{"result": result}).Debug("result from db")

	return result, nil
}

func InsertRow(row map[string]interface{}, dbName string, collectionName string) (err error) {
	session, err := connect()
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(dbName).C(collectionName)
	err = c.Insert(row)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"data": row}).Info("Inserted to db")
	return nil
}

func Delete(db string, collection string, condition bson.M) (err error) {
	session, err := connect()
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(db).C(collection)
	info, err := c.RemoveAll(condition)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"Matched": info.Matched, "Deleted": info.Removed, "Updated": info.Updated}).Info("Deleted from db")
	return nil
}
