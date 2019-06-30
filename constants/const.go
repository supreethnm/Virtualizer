package constants

const (
	// Mongo DB
	MONGO_DB_FIELD_ID     = "_id"
	MONGO_DB_AND_OPERATOR = "$and"

	CONFIG_ENDPOINT = `
	[[services]]
    sname="Virtualizer"
    path="/updateConfig"
    type="text"
	method="post"
	
	[[services]]
    sname="JSONservice"
    path="/getConfig"
    type="text"
    method="get"
	`
	CONFIG_ENDPOINT_UPDATE_CONFIG = "updateConfig"
	CONFIG_ENDPOINT_GET_CONFIG    = "getConfig"

	DB_ENDPOINTS = `
    [[services]]
    sname="JSONservice"
    path="/insertData"
    type="text/json"
    method="post"

    [[services]]
    sname="JSONservice"
    path="/getData"
    type="text/json"
    method="get"

    [[services]]
    sname="JSONservice"
    path="/delete"
    type="text/json"
    method="delete"
	`
	DB_ENDPOINT_INSERT_DATA = "insertData"
	DB_ENDPOINT_GET_DATA    = "getData"
	DB_ENDPOINT_DELETE_DATA = "delete"

	// DB and Collection string
	STRING_DATABASE   = "database"
	STRING_COLLECTION = "collection"

	// TODO: provide flexibility for users to use other ports
	PORT = "8080"

	// TODO: revisit config file location
	//CONFIG_FILE = "/tmp/virtualizer/config.toml"

	// HTTP vars
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_DELETE = "DELETE"

	// Misc
	STRING_CONTENT_TYPE     = "content-type"
	STRING_XML              = "xml"
	STRING_APPLICATION_JSON = "application/json"

	// Default config string
	DEFAULT_CONFIG_STRING = `
# config.toml not found in config directory
# Below is a sample content.

[[services]]
sname="Service1"
path="/path/to/service1"
type="text/json"
method="get"
database="testdb"
collection="testcol"
delay=0
reference="""
{
	"json_key": "json_value"
}
"""
omit=["any_json_field"]
response="""
{
	"message" : "default message"
}
"""

# -------------------------------------------------------------------------------------------

[[services]]
sname="Service2"
path="/path/to/service2"
type="text/json"
method="get"
delay=0
database="testdb"
collection="testcol2"
reference="""
{
	"json_key": "json_value"
}
"""
response="""
{
	"message" : "default message"
}
"""

# -------------------------------------------------------------------------------------------

[[services]]
sname="XMLService"
path="/path/to/xmlservice"
type="text/xml"
method="post"
database="testdb"
collection="testcol3"
delay=0
reference="""
{
"RootObject.Body.Element": "RootObject.Body.Element",
}
"""
response="""
{
	"message" : "default message"
}
"""
`
)
