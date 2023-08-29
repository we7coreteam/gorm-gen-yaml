package gorm_gen_yaml

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/gen"
	"io"
	"os"
)

type yamlGenerate struct {
	yaml           *DbRelation
	gen            *gen.Generator
	tableModel     []interface{}
	tableModelSize int
}

func NewYamlGenerate(generator *gen.Generator, path string) (*yamlGenerate, error) {
	obj := &yamlGenerate{
		tableModelSize: 100,
	}
	err := obj.loadFromFile(path)
	if err != nil {
		return nil, err
	}
	obj.gen = generator
	obj.tableModel = make([]interface{}, obj.tableModelSize)
	return obj, nil
}

type DbRelation struct {
	Relation []Relation `yaml:"relation"`
}

type Relation struct {
	Table  string          `yaml:"table"`
	Relate []RelationTable `yaml:"relate"`
	isGen  bool
}

type RelationTable struct {
	Table      string `yaml:"table"`
	ForeignKey string `yaml:"foreign_key"`
	Type       string `yaml:"type"`
}

func (self *yamlGenerate) loadFromFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("%s file not found", path))
	}
	content, _ := io.ReadAll(file)
	self.yaml = &DbRelation{}
	err = yaml.Unmarshal(content, self.yaml)
	if err != nil {
		return err
	}
	fmt.Printf("%+v \n", self.yaml)

	return nil
}

func (self *yamlGenerate) Generate(opt ...gen.ModelOpt) []interface{} {
	tableModels := make([]interface{}, 100)
	for _, item := range self.yaml.Relation {
		if len(item.Relate) > 0 {
			for _, relateItem := range item.Relate {
				fmt.Println(relateItem.Table)
			}
		}
	}
	//self.gen.GenerateModel(hasOne.Table, nil)
	return tableModels
}

func (self *yamlGenerate) generateTableModel(tableName string) {

}
