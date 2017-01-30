package ddb

type Command struct {
	operation string
	target    string
	data      map[string]interface{}
}

func (c *Command) initFromJson(i interface{}) {

	if i == nil {
		return
	}

	m := i.(map[string]interface{})
	command, ok := m["command"].(string)

	if ok == true {
		c.operation = command
	}

	target, ok := m["target"].(string)

	if ok == true {
		c.target = target
	}

	data, ok := m["data"].(map[string]interface{})

	if ok == true {
		c.data = data
	}
}
