package constants

const (
	// Mongo DB
	MONGO_DB_HOST         = "localhost"
	MONGO_DB_FIELD_ID     = "_id"
	MONGO_DB_AND_OPERATOR = "$and"

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

	// Config
	CONFIG_FILE = "config.toml"
	PORT        = "8080"

	// HTTP vars
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_DELETE = "DELETE"

	// Misc
	STRING_CONTENT_TYPE     = "content-type"
	STRING_XML              = "xml"
	STRING_APPLICATION_JSON = "application/json"
)
