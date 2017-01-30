package ddb


import (
	"fmt"
	"github.com/google/uuid"
)

type Target struct {
	name string
	data map[string]interface{}
}

func (target *Target) insert(data map[string]interface{}) (id string){
	foo := uuid.Must(uuid.NewRandom())

	id = foo.String()

	fmt.Println("Generated ID is: ", id)

	target.data[id] = data

	return
}

func (target Target) read(data map[string]interface{}) map[string]interface{} {

	var id string
	var ok bool

	if id, ok = data["_id"].(string); !ok {
		fmt.Println("Could not get _id from Request")
		return make(map[string]interface{})
	}

	fmt.Println("Found id:", id, "in Request")

	result, ok := target.data[id];

	if ok == false {
		return make(map[string]interface{})
	}

	if foo, ok := result.(map[string]interface{}); ok {
		return foo
	} else {
		return make(map[string]interface{})
	}

}

func (target Target) remove(data map[string]interface{}) bool {

	return false;
}