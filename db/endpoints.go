package db

//TODO: refactor the code
// Just dumped few endpoints for moving forward

func GetDBEndpoints() (result string) {

	result = `
    [[services]]
    sname="JSONservice"
    path="/insertData"
    type="text/json"
    method="post"
        [[services.operations]]
            opname=""
            delay=0
            monitoring=true
                [[services.operations.outputs]]
                tagvalue=""
                response=""				
                [services.operations.outputs.variables]


    [[services]]
    sname="JSONservice"
    path="/getData"
    type="text/json"
    method="get"
        [[services.operations]]
            opname=""
            delay=0
            monitoring=true
                [[services.operations.outputs]]
                tagvalue=""
                response=""	
                [services.operations.outputs.variables]

    [[services]]
    sname="JSONservice"
    path="/delete"
    type="text/json"
    method="delete"
        [[services.operations]]
            opname=""
            delay=0
            monitoring=true
                [[services.operations.outputs]]
                tagvalue=""
                response=""				
                [services.operations.outputs.variables]
    `

	return
}
