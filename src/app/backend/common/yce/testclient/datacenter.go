package testclient

const (
	CHECK_DATACENTER = `{"name":"test-datacenter","orgId":"1"}`
	CREATE_DATACENTER = `{"name":"test-datacenter","nodePort":["30000","32767"],"host":"192.168.1.110","port":8080,"orgId":3","op":1}`
	DELETE_DATACENTER = `{"name":"test-datacenter","orgId":"3","op":1}`
	UPDATE_DATACENTER = `{"name":"test-datacenter","nodePort":["31000","32767"],"host":"192.168.1.110","port":8080,"orgId":3","op":1}`
)
