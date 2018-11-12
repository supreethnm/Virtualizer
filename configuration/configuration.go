package configuration

import "gopkg.in/mgo.v2/bson"

type Config struct {
	Services []Service
}

type Service struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Sname      string        //`bson:"Sname"`
	Path       string        //`bson:"EndPoint"`
	Type       string        //`bson:"Type"`
	Method     string        //http method
	Database   string
	Collection string
	Delay      int    // `bson:"Delay"`
	Response   string //`bson:"Response"`
	Reference  string //map[string]interface{}
	Omit       []string
}
