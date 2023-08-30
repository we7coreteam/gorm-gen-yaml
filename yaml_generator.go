package gorm_gen_yaml

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"io"
	"os"
)

type yamlGenerator struct {
	yaml           *DbRelation
	gen            *gen.Generator
	tableModelSize int
	generatedTable map[string]interface{}
}

func NewYamlGenerator(path string) *yamlGenerator {
	obj := &yamlGenerator{
		tableModelSize: 100,
	}
	err := obj.loadFromFile(path)
	if err != nil {
		return nil
	}
	obj.generatedTable = make(map[string]interface{})
	return obj
}

type DbRelation struct {
	Relation    []Relation `yaml:"relation"`
	RelationMap map[string]Relation
}

type Relation struct {
	Table  string          `yaml:"table"`
	Relate []RelationTable `yaml:"relate"`
}

type RelationTable struct {
	Table      string `yaml:"table"`
	ForeignKey string `yaml:"foreign_key"`
	Type       string `yaml:"type"`
}

func (self *yamlGenerator) UseGormGenerator(g *gen.Generator) *yamlGenerator {
	self.gen = g
	return self
}

func (self *yamlGenerator) loadFromFile(path string) error {
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

func (self *yamlGenerator) generateFromRelation(relation Relation) {
	if _, exists := self.generatedTable[relation.Table]; exists {
		return
	}

	for _, relate := range relation.Relate {
		if trelation, exists := self.yaml.RelationMap[relate.Table]; exists {
			self.generateFromRelation(trelation)
		} else {
			self.generatedTable[relate.Table] = self.gen.GenerateModel(relate.Table)
		}
	}

	//找到所有relate,生成模型
	opt := make([]gen.ModelOpt, len(relation.Relate))
	for i, table := range relation.Relate {
		var fieldType field.RelationshipType
		switch table.Type {
		case "has_one":
			fieldType = field.HasOne
		case "has_many":
			fieldType = field.HasMany
		case "many_many":
			fieldType = field.Many2Many
		case "belongs_to":
			fieldType = field.BelongsTo
		}
		opt[i] = gen.FieldRelate(fieldType, table.Table, self.generatedTable[table.Table],
			&field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{table.ForeignKey}},
			})
	}
	self.generatedTable[relation.Table] = self.gen.GenerateModel(relation.Table, opt...)
}

func (self *yamlGenerator) Generate(opt ...gen.ModelOpt) []interface{} {
	for _, relation := range self.yaml.Relation {
		self.generateFromRelation(relation)
	}
	tableModels := make([]interface{}, len(self.generatedTable))
	for _, item := range self.generatedTable {
		tableModels = append(tableModels, item)
	}
	return tableModels
}