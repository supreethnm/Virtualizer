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
	Operations []Operation   //`bson:"Operations"`
}

type Operation struct {
	Opname           string //`bson:"Opname"`
	Database         string
	Collection       string
	Delay            int      // `bson:"Delay"`
	Outputs          []Output // `bson:"Output"`
	Monitoring       bool
	MultipleResponse int
	//Path string

}

type Output struct {
	Variables map[string]string
	Tagvalue  string //`bson:"TagName"`
	Response  string //`bson:"Response"`
	Reference string //map[string]interface{}
	Omit      []string
}
