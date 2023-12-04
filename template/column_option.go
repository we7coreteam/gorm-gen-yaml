package template

// Model used as a variable because it cannot load template file after packed, params still can pass file
var ColumnOptionTemplate = map[string]string{"common": `package {{Package}}

import (
    "database/sql/driver"
)

type {{OptionStructName}} struct {

}

func (c {{OptionStructName}}) Value() (driver.Value, error) {
	return c, nil
}

func (c *{{OptionStructName}}) Scan(value interface{}) error {
	return nil
}

`,
	"json": `package {{Package}}

type {{OptionStructName}} struct {

}
`,
}
