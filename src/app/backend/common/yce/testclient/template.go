package testclient

const (

	// Data below is not standard deployment json and service json, it's yce's deployment and service

	CHECK_TEMPLATE = `{"name":"nginx-template"}`
	//CREATE_TEMPLATE = `{"name":"nginx-template","deployment":{"apiVersion":"extensions/v1beta1","kind":"Deployment","metadata":{"name":"nginx-template"}},"service":{"apiVersion":"v1","kind":"Service","metadata":{"name":"nginx-svc-template"}}}`
	CREATE_TEMPLATE = `{"name":"nginx-template", "deployment":{}, "service":{}}`
	//CREATE_TEMPLATE = `{"name":"nginx-template"}`
	DELETE_TEMPLATE = `{"name":"nginx-template","id":7}`
	UPDATE_TEMPLATE = `{"name":"nginx-template","deployment":{"apiVersion":"extensions/v1beta1","kind":"Deployment","metadata":{"name":"nginx-template","namespace":"test-template"}},"service":{"apiVersion":"v1","kind":"Service","metadata":{"name":"nginx-svc-template", "namespace":"test-template"}}}`

)
