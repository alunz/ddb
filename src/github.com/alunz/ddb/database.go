package ddb

import "fmt"

type Database struct {
	targets map[string]Target
}

var instanced int = 0
var instance Database;

func getInstance() Database {

	if instanced == 0 {
		instance = Database{targets:make(map[string]Target)}
		instanced = 1;
	}

	for foo, bar := range instance.targets {
		fmt.Println(foo, bar);
	}

	return instance
}
