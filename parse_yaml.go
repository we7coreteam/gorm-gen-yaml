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
	generatedTable map[string]interface{}
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
	obj.generatedTable = make(map[string]interface{})
	return obj, nil
}

type DbRelation struct {
	Relation    []Relation `yaml:"relation"`
	RelationMap map[string]Relation
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
	self.yaml.RelationMap = make(map[string]Relation)

	for _, relation := range self.yaml.Relation {
		self.yaml.RelationMap[relation.Table] = relation
	}

	return nil
}

func (self *yamlGenerate) generateFromRelation(relation Relation) {
	if _, exists := self.generatedTable[relation.Table]; exists {
		return
	}

	for _, relate := range relation.Relate {
		if trelation, exists := self.yaml.RelationMap[relate.Table]; exists {
			self.generateFromRelation(trelation)
		} else {
			self.generatedTable[relate.Table] = self.generateTableModel(relate.Table)
		}
	}

	//找到所有relate,生成模型
	self.generatedTable[relation.Table] = self.generateTableModel(relation.Table)
}

func (self *yamlGenerate) GenerateV1(opt ...gen.ModelOpt) {
	for _, relation := range self.yaml.Relation {
		self.generateFromRelation(relation)
	}
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

func (self *yamlGenerate) generateTableModel(tableName string) interface{} {
	return 1
}
