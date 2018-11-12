# **Virtualizer:**
Used for mocking REST based services

### How to run the Virtualizer:
1. Make sure we have Virtualizer executable in place
2. Make sure the config.toml file is in the same directory of Virtualizer executable (Virtualizer will read this config.toml file)
3. Run the Virtualizer executable. An example is shown below
Example:
```
$ ./virtualizer
INFO[0000] Request for GET /EndPoint1 
INFO[0000] Request for POST /EndPoint2
INFO[0000] Request for POST /insertData                 
INFO[0000] Request for GET /getData                     
INFO[0000] Request for DELETE /delete                   
INFO[0000] Endpoint created at path: /EndPoint1
INFO[0000] Endpoint created at path: /EndPoint2
INFO[0000] Endpoint created at path: /insertData        
INFO[0000] Endpoint created at path: /getData           
INFO[0000] Endpoint created at path: /delete            
INFO[0000] Listening on port 8080
```

**NOTE: Virtualizer needs a mongo db running in the background WITH NO CREDENTIALS setup**


## Config.toml file:

A Sample config.toml file that creates Service1, Service2 services:

```
[services]]
sname="Service1"
path="/EndPoint1"
method="get"
database="test"
collection="service1"
delay=0
reference="""
{
"ipAddress": "ipAddress"
}
"""
omit=["ipAddress"]
response="""
{
"status":"failed"
}
"""

# -------------------------------------------------------------------------------------------

[[services]]
sname="Service2"
path="/EndPoint2"
method="post"
delay=0
database="test"
collection="service2"
reference="""
{
"service2.accountNumber": "service2.accountNumber", 
"service2.siteCode": "service2.siteCode"
}
"""
response="""
{
"status":"failed"
}
"""
```

### Tags and its uses:
```
"[[services]]"                      :   This tag indicates the starting point to create a new service

"sname"                             :   This tag is used for naming the service

"path"                              :   This tag is used for creating a endpoint and should be unique

"type"                              :   This tag defines the HTTP response format (Eg. JSON or XML).
                                        Currently, Virtualizer supports JSON (by default) and XML formats only.
                                        For XML, assign the value "text/xml" or "application/xml".
                                        For JSON, assign the value "text/json" or "application/json"

"method"                            :   Tag that specifies HTTP method user intend to setup for a service 

"[[services.operations]]"           :   Tag used to indicate Virtualizer the operation(s) in a service

"database"                          :   To specify the database where Virtualizer can find the response(s), Collection name shoule also be specified under "collection" tag.

"collection"                        :   Collection/Table name where Virtualizer can find the response

"delay"                             :   Intended delay for getting the response

"[[services.operations.outputs]]"   :   Tag used for defining the output related params

"references"                        :   To map keys in the request to a response (to be specified in JSON format)
                                        Should specify the quick path to the keys of request and response
                                        Example: 
                                        Say, for the following JSON data:
                                        {
                                            "key1": "value1"
                                            "key2": {
                                                "key3": "value3"
                                                "key4": "value4"
                                                "key5": {
                                                    "key6": "value6"
                                                }
                                            }
                                        }

                                        The path to get to "value1" is "key1". Similarly, path to get to "value6" is "key2.key5.key6"

                                        The KEY part in this json data shoule be for the REQUEST and the VALUE part should be for the RESPONSE

                                        Example:
                                        Say, a POST request contains following json body:									  
                                        {
                                            "service2": {
                                                "accountNumber": "1234567890",
                                                "systemID": "Gateway",
                                                "siteCode": "East"
                                            }
                                        }

                                        And user expects following response (Which should be presend in the database):
                                        {
                                            "service2": {
                                                "accountNumber": "1234567890",
                                                "accountStatus": "Active",
                                                "customerClassification": "Retail",
                                                "siteCode": "East"
                                            }
                                        }

                                        Here, request to response reference keys are mapped as example shown below:
                                        {
                                            "service2.accountNumber": "service2.accountNumber", 
                                            "service2.siteCode": "service2.siteCode"
                                        }


                                        So, "service2.accountNumber" in the REQUEST is mapped to "service2.accountNumber" in the RESPONSE.
                                        Similarly, "service2.siteCode" in the REQUEST is mapped to "service2.siteCode" in the RESPONSE.

                                        Now, Virtualizer looks for the value in the specifed path of the request and retrieves the data from the DB only if the value matches with the specified path of the response.

                                        NOTE: For GET request the "key" reference will be searched in query param of the request
                                        Example:
                                        Say, the CURL for a get request is as shown below:

                                        $ curl "http://localhost:8080/EndPoint1?ipAddress=10.10.5.5"

                                        And user expects following response (Which should be presend in the database):
                                        {
                                            "ipAddress": "10.10.5.5"
                                            "accountNumber": "1234567890",
                                            "siteCode": "East"
                                        }

                                        The reference data now looks like:
                                        {
                                            "ipAddress": "ipAddress"
                                        }


"omit"                              :   Specify to omit any field from the response. Its an array type that contains "quick path" to the field that needs to excluded from the response
                                        Example:
                                        Say, user expects following response:
                                        {
                                            "accountNumber": "1234567890",
                                            "companyId": 3,
                                            "siteCode": "East"
                                        }

                                        But, the response in DB is:
                                        {
                                            "ipAddress": "10.10.5.5"
                                            "accountNumber": "1234567890",
                                            "companyId": 3,
                                            "siteCode": "East"
                                        }

                                        Simply add `["ipAddress"]` to "omit" tag


"response"                          :   A default response if the data is not found in the DB

```

## DB operations
In addition to user specified endpoints, Virtualizer creates 3 more endpoint for DB operations:
1. `insertData`
2. `getData`
3. `delete`

### Insert data into the DB

Below request to the Virtualizer inserts data into the DB
```
curl -X POST \
  'http://localhost:8080/insertData?database=test&collection=service1' \
  -H 'Content-Type: application/json' \
  -d '{
    "ipAddress":"10.10.5.5",
    "accountNumber": "1234567890",
    "companyId":3,
    "siteCode":"East"
}'
```
`database` name and the `collection` name should be specified as part of query params. 
Here, Below data will be inserted into `service1` collection in `test` database:
```
{
    "ipAddress":"10.10.5.5",
    "accountNumber": "1234567890",
    "companyId":3,
    "siteCode":"East"
}
```
**NOTE: User can also feed XML data to Virtualizer. However, Virtualizer converts the XML data to JSON format and inserts into the DB. User can always initialize "type" (defined above) to get the response in a particular format.**
### Get data from the DB

Below request to the Virtualizer retrieves all documents in a given collection
```
curl -X GET 'http://localhost:8080/getData?database=test&collection=service1'
```
`database` name and the `collection` name should be specified as part of query params. 
Here, name of the database is `test` and the collection is `service1`

### Delete data from the DB

Below request to the Virtualizer deletes **all** the documents in a given collection
```
curl -X DELETE \
  'http://localhost:8080/delete?database=test&collection=service1' \
  -H 'Content-Type: application/json' \
  -d '{}'
```
**Observe, the request body here is empty JSON object.**

`database` name and the `collection` name should be specified as part of query params. 
Here, name of the database is `test` and the collection is `service1`. 

Now,
Below request to the Virtualizer deletes the document(s) in a given collection with a condition specified in request body
```
curl -X DELETE \
  'http://localhost:8080/delete?database=test&collection=service2' \
  -H 'Content-Type: application/json' \
  -d '{
    "service2.accountNumber": "1234567890", 
    "service2.accountStatus": "Active"
}'
```
`database` name and the `collection` name should be specified as part of query params. 
Here, name of the database is `test` and the collection is `service2`


## Miscellaneous


Step to redirect a source to talk to Virtualizer:
1. Login to source instance (make sure you have the root privilege)
2. Change the hosts file (/etc/hosts)

Added lines similar to below example.

Example:
```
192.168.1.4 actual_service_host
```

Here, the IP of the Virtualizer is "192.168.1.4"

