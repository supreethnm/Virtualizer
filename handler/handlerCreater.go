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
	"virtualizer/db"
	u "virtualizer/utils"
)

func PostHandler(service c.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// initialize response headers
		w.Header().Set("content-type", service.Type)

		// read request body
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.WithFields(logrus.Fields{}).Error(err.Error() + "\nVirtualizer: Error reading request body!")
			http.Error(w, err.Error()+"\nVirtualizer: Error reading request body!", http.StatusInternalServerError)
			return
		}
		// if XML type then convert to JSON
		if strings.Contains(strings.ToLower(r.Header.Get("content-type")), "xml") {
			reqBodyBytes, err = u.ToJsonBytes(reqBodyBytes)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// check if the request is for db insert
		if strings.Contains(r.URL.Path, "insertData") {
			var row map[string]interface{}
			err = json.Unmarshal(reqBodyBytes, &row)
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error() + "\nVirtualizer: Error UnMarshalling request body!")
				http.Error(w, err.Error()+"\nVirtualizer: Error UnMarshalling request body!", http.StatusInternalServerError)
				return
			}
			err = db.InsertRow(row, r.URL.Query().Get("database"), r.URL.Query().Get("collection"))
			if err != nil {
				logrus.WithFields(logrus.Fields{}).Error(err.Error() + "\nVirtualizer: Error inserting data into the DB!")
				http.Error(w, err.Error()+"\nVirtualizer: Error inserting data into the DB!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}

		for _, op := range service.Operations {

			// delay
			ch := make(chan bool, 1)
			u.AddDelay(time.Duration(op.Delay), ch)

			for _, output := range op.Outputs {

				var bsonC []bson.M
				var tempMap map[string]interface{}
				err = json.Unmarshal([]byte(output.Reference), &tempMap)
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
				logrus.WithFields(logrus.Fields{"DatabaseName": op.Database, "CollectionName": op.Collection}).Debug("POST: fetching data from DB...")

				// fetch data from db
				data, err := db.GetData(op.Database, op.Collection, bsonC)
				if err != nil {
					// check for any default response
					if len(output.Response) > 0 {
						logrus.WithFields(logrus.Fields{"output.Response": output.Response}).Debug("default response")
						// check if response type is XML
						if strings.Contains(strings.ToLower(service.Type), "xml") {
							resp, err := u.ToXmlBytes([]byte(output.Response))
							if err != nil {
								logrus.WithFields(logrus.Fields{}).Error(err.Error())
								http.Error(w, err.Error(), http.StatusInternalServerError)
								return
							}
							output.Response = string(resp)
						}
						fmt.Fprintf(w, output.Response)
						break
					}
					logrus.WithFields(logrus.Fields{}).Error(err.Error())
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}

				// omit the field if specified
				bData := []byte(data)
				for _, ov := range output.Omit {
					path := strings.Split(ov, ".")
					bData = jsonparser.Delete(bData, path...)
				}

				// check if response type is XML
				if strings.Contains(strings.ToLower(service.Type), "xml") {
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

				// TODO: Observe the usage of outputs. Remove it if no need
				break
			}

		}
	}
}

func GetHandler(service c.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// initialize response headers
		w.Header().Set("content-type", service.Type)

		// check if we need to get data from db
		if strings.Contains(r.URL.Path, "getData") {

			result, err := db.GetData(r.URL.Query().Get("database"), r.URL.Query().Get("collection"), nil)
			if err != nil {
				http.Error(w, err.Error()+"\nVirtualizer: Error getting data from the DB!", http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, result)
			return
		}

		// initialize response headers
		w.Header().Set("content-type", service.Type)

		// Get the query params
		queryParams, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, op := range service.Operations {

			// delay
			ch := make(chan bool, 1)
			u.AddDelay(time.Duration(op.Delay), ch)

			for _, output := range op.Outputs {

				var bsonC []bson.M
				var tempMap map[string]interface{}
				err = json.Unmarshal([]byte(output.Reference), &tempMap)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				for k, v := range tempMap {
					temp := bson.M{v.(string): queryParams.Get(k)}
					bsonC = append(bsonC, temp)
				}
				logrus.WithFields(logrus.Fields{"bsonC": bsonC}).Debug("GET: condition data")
				logrus.WithFields(logrus.Fields{"DatabaseName": op.Database, "CollectionName": op.Collection}).Debug("GET: fetching data from DB...")
				// fetch data from db
				data, err := db.GetData(op.Database, op.Collection, bsonC)
				if err != nil {
					// check for any default response
					if len(output.Response) > 0 {
						logrus.WithFields(logrus.Fields{"output.Response": output.Response}).Debug("default response")
						// check if response type is XML
						if strings.Contains(strings.ToLower(service.Type), "xml") {
							resp, err := u.ToXmlBytes([]byte(output.Response))
							if err != nil {
								logrus.WithFields(logrus.Fields{}).Error(err.Error())
								http.Error(w, err.Error(), http.StatusInternalServerError)
								return
							}
							output.Response = string(resp)
						}
						fmt.Fprintf(w, output.Response)
						break
					}
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}

				// omit the field if specified
				bData := []byte(data)
				for _, ov := range output.Omit {
					path := strings.Split(ov, ".")
					bData = jsonparser.Delete(bData, path...)
				}

				// check if response type is XML
				if strings.Contains(strings.ToLower(service.Type), "xml") {
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

				// TODO: Observer the usage of outputs. Remove it if no need
				break
			}

		}
	}
}

func DeleteHandler(service c.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if we need to delete data from db
		if strings.Contains(r.URL.Path, "delete") {

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

			err = db.Delete(r.URL.Query().Get("database"), r.URL.Query().Get("collection"), bsonC)
			if err != nil {
				http.Error(w, err.Error()+"\nVirtualizer: Error deleting data in the DB!", http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/text")
			logrus.WithFields(logrus.Fields{}).Debug("deleted data from the db")
			fmt.Fprintf(w, "Deleted")

			return
		}

		http.Error(w, "Error", http.StatusInternalServerError)
	}
}
