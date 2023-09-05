package template

// Model used as a variable because it cannot load template file after packed, params still can pass file
const CustomColumnTemplate = `
package {{Package}}

import (
    "database/sql/driver"
)

type {{CustomStructName}} struct {

}

func (c {{CustomStructName}}) Value() (driver.Value, error) {
	//if options == nil {
	//	return "", nil
	//}

	//return json.Marshal(options)
	return c, nil
}

func (c *{{CustomStructName}}) Scan(value interface{}) error {
	//b, ok := value.([]byte)
	//if !ok {
	//	return fmt.Errorf("value is not []byte, value: %v", value)
	//}
	//if len(b) == 0 {
	//	return nil
	//}

	//return json.Unmarshal(b, &options)

	return nil
}

`
