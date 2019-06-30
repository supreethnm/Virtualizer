package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/buger/jsonparser"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2/bson"

	c "virtualizer/configuration"
	cn "virtualizer/constants"
	"virtualizer/db"
	u "virtualizer/utils"
)

func PostHandler(service c.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// initialize response headers
		w.Header().Set(cn.STRING_CONTENT_TYPE, service.Type)

		// read request body
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.WithFields(logrus.Fields{}).Error(err.Error() + "\nVirtualizer: Error reading request body!")
			http.Error(w, err.Error()+"\nVirtualizer: Error reading request body!", http.StatusInternalServerError)
			return
		}
		// if XML type then convert to JSON
		if strings.Contains(strings.ToLower(r.Header.Get(cn.STRING_CONTENT_TYPE)), cn.STRING_XML) {
			reqBodyBytes, err = u.ToJsonBytes(reqBodyBytes)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// check for config update
		if strings.Contains(r.URL.Path, cn.CONFIG_ENDPOINT_UPDATE_CONFIG) {
			logrus.WithFields(logrus.Fields{"data": string(reqBodyBytes)}).Info("Updating virtualizer config...")
			err = ioutil.WriteFile(u.GetConfigFilePath(), reqBodyBytes, 0755)
			if err != nil {
				errString := "Unable to update virtualizer config file: " + err.Error()
				logrus.WithFields(logrus.Fields{}).Error(errString)
				http.Error(w, errString, http.StatusInternalServerError)
				return
			}
			logrus.WithFields(logrus.Fields{}).Info("Successfully updated virtualizer config")
			w.WriteHeader(http.StatusCreated)
			infoStr := "Please restart your service/docker to reflect the updated config changes!"
			w.Write([]byte(infoStr))
			logrus.WithFields(logrus.Fields{}).Info(infoStr)
			return
		}

		// check if the request is for db insert
		if strings.Contains(r.URL.Path, cn.DB_ENDPOINT_INSERT_DATA) {
			var row map[string]interface{}
			err = json.Unmarshal(reqBodyBytes, &row)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error() + "\nVirtualizer: Error UnMarshalling request body!")
				http.Error(w, err.Error()+"\nVirtualizer: Error UnMarshalling request body!", http.StatusInternalServerError)
				return
			}
			err = db.InsertRow(row, r.URL.Query().Get(cn.STRING_DATABASE), r.URL.Query().Get(cn.STRING_COLLECTION))
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error() + "\nVirtualizer: Error inserting data into the DB!")
				http.Error(w, err.Error()+"\nVirtualizer: Error inserting data into the DB!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}

		// delay
		ch := make(chan bool, 1)
		u.AddDelay(time.Duration(service.Delay), ch)

		var bsonC []bson.M
		var tempMap map[string]interface{}
		err = json.Unmarshal([]byte(service.Reference), &tempMap)
		if err != nil {
			logrus.WithFields(logrus.Fields{}).Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for k, v := range tempMap {
			value := gjson.Get(string(reqBodyBytes), k)
			temp := bson.M{v.(string): value.String()}
			bsonC = append(bsonC, temp)
		}

		logrus.WithFields(logrus.Fields{"bsonC": bsonC}).Debug("POST: condition data")
		logrus.WithFields(logrus.Fields{"DatabaseName": service.Database, "CollectionName": service.Collection}).Debug("POST: fetching data from DB...")

		// fetch data from db
		data, err := db.GetData(service.Database, service.Collection, bsonC)
		if err != nil {
			// check for any default response
			if len(service.Response) > 0 {
				logrus.WithFields(logrus.Fields{"service.Response": service.Response}).Debug("default response")
				// check if response type is XML
				if strings.Contains(strings.ToLower(service.Type), cn.STRING_XML) {
					resp, err := u.ToXmlBytes([]byte(service.Response))
					if err != nil {
						logrus.WithFields(logrus.Fields{}).Error(err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					service.Response = string(resp)
				}
				// TODO: Do we need appropriate HTTP code for this?
				fmt.Fprintf(w, service.Response)
				return
			}
			logrus.WithFields(logrus.Fields{}).Error(err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// omit the field if specified
		bData := []byte(data)
		for _, ov := range service.Omit {
			path := strings.Split(ov, ".")
			bData = jsonparser.Delete(bData, path...)
		}

		// check if response type is XML
		if strings.Contains(strings.ToLower(service.Type), cn.STRING_XML) {
			bData, err = u.ToXmlBytes(bData)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data = string(bData)
		logrus.WithFields(logrus.Fields{"data": data}).Debug("response from db")
		fmt.Fprintf(w, data)

	}
}

func GetHandler(service c.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// initialize response headers
		w.Header().Set(cn.STRING_CONTENT_TYPE, service.Type)

		// check if we need to get data from db
		if strings.Contains(r.URL.Path, cn.DB_ENDPOINT_GET_DATA) {

			result, err := db.GetData(r.URL.Query().Get(cn.STRING_DATABASE), r.URL.Query().Get(cn.STRING_COLLECTION), nil)
			if err != nil {
				http.Error(w, err.Error()+"\nVirtualizer: Error getting data from the DB!", http.StatusInternalServerError)
				return
			}
			w.Header().Add(cn.STRING_CONTENT_TYPE, cn.STRING_APPLICATION_JSON)
			fmt.Fprintf(w, result)
			return
		}

		// check for config update
		if strings.Contains(r.URL.Path, cn.CONFIG_ENDPOINT_GET_CONFIG) {
			data, err := ioutil.ReadFile(u.GetConfigFilePath())
			if err != nil {
				errString := "Error in reading virtualizer config file: " + err.Error()
				logrus.WithFields(logrus.Fields{}).Error(errString)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(cn.DEFAULT_CONFIG_STRING))
				return
			}
			logrus.WithFields(logrus.Fields{"data": string(data)}).Info("virtualizer config file data")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}

		// initialize response headers
		w.Header().Set(cn.STRING_CONTENT_TYPE, service.Type)

		// Get the query params
		queryParams, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delay
		ch := make(chan bool, 1)
		u.AddDelay(time.Duration(service.Delay), ch)

		var bsonC []bson.M
		var tempMap map[string]interface{}
		err = json.Unmarshal([]byte(service.Reference), &tempMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for k, v := range tempMap {
			temp := bson.M{v.(string): queryParams.Get(k)}
			bsonC = append(bsonC, temp)
		}
		logrus.WithFields(logrus.Fields{"bsonC": bsonC}).Debug("GET: condition data")
		logrus.WithFields(logrus.Fields{"DatabaseName": service.Database, "CollectionName": service.Collection}).Debug("GET: fetching data from DB...")
		// fetch data from db
		data, err := db.GetData(service.Database, service.Collection, bsonC)
		if err != nil {
			// check for any default response
			if len(service.Response) > 0 {
				logrus.WithFields(logrus.Fields{"service.Response": service.Response}).Debug("default response")
				// check if response type is XML
				if strings.Contains(strings.ToLower(service.Type), cn.STRING_XML) {
					resp, err := u.ToXmlBytes([]byte(service.Response))
					if err != nil {
						logrus.WithFields(logrus.Fields{}).Error(err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					service.Response = string(resp)
				}
				// TODO: Do we need appropriate HTTP code for this?
				fmt.Fprintf(w, service.Response)
				return
			}
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// omit the field if specified
		bData := []byte(data)
		for _, ov := range service.Omit {
			path := strings.Split(ov, ".")
			bData = jsonparser.Delete(bData, path...)
		}

		// check if response type is XML
		if strings.Contains(strings.ToLower(service.Type), cn.STRING_XML) {
			bData, err = u.ToXmlBytes(bData)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data = string(bData)
		logrus.WithFields(logrus.Fields{"data": data}).Debug("response from db")
		fmt.Fprintf(w, data)

	}
}

func DeleteHandler(service c.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if we need to delete data from db
		if strings.Contains(r.URL.Path, cn.DB_ENDPOINT_DELETE_DATA) {

			reqBodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error() + "\nVirtualizer: Error reading request body!")
				http.Error(w, err.Error()+"\nVirtualizer: Error reading request body!", http.StatusInternalServerError)
				return
			}

			var bsonC bson.M
			err = bson.UnmarshalJSON(reqBodyBytes, &bsonC)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = db.Delete(r.URL.Query().Get(cn.STRING_DATABASE), r.URL.Query().Get(cn.STRING_COLLECTION), bsonC)
			if err != nil {
				http.Error(w, err.Error()+"\nVirtualizer: Error deleting data in the DB!", http.StatusInternalServerError)
				return
			}
			w.Header().Add(cn.STRING_CONTENT_TYPE, "application/text")
			logrus.WithFields(logrus.Fields{}).Debug("deleted data from the db")
			fmt.Fprintf(w, "Deleted")

			return
		}

		http.Error(w, "Error", http.StatusInternalServerError)
	}
}
