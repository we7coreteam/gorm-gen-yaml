package yamlgen

import (
	"errors"
	"fmt"
	"github.com/we7coreteam/gorm-gen-yaml/template"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"io"
	"os"
	"path"
	"strings"
)

type yamlGenerator struct {
	yaml                *DbTable
	gen                 *gen.Generator
	generatedTable      map[string]string
	customColumnSaveDir string
}

func NewYamlGenerator(path string) *yamlGenerator {
	obj := &yamlGenerator{}
	err := obj.loadFromFile(path)
	if err != nil {
		panic(err)
		return nil
	}
	obj.generatedTable = make(map[string]string)
	return obj
}

type DbTable struct {
	Table    []Table `yaml:"table"`
	TableMap map[string]Table
}

type Table struct {
	Name   string            `yaml:"table"`
	Relate []Relate          `yaml:"relate"`
	Column map[string]Column `yaml:"column"`
}

type Column struct {
	CustomColumnType string `yaml:"custom_type"`
}

type Relate struct {
	Table          string `yaml:"table"`
	ForeignKey     string `yaml:"foreign_key"`
	References     string `yaml:"references"`
	JoinForeignKey string `yaml:"join_foreign_key"`
	JoinReferences string `yaml:"join_references"`
	Many2many      string `yaml:"many_2_many"`
	Type           string `yaml:"type"`
}

func (self *yamlGenerator) UseGormGenerator(g *gen.Generator) *yamlGenerator {
	self.gen = g

	self.SetCustomColumnSaveDir(g.Config.ModelPkgPath + "/../custom")

	return self
}

func (self *yamlGenerator) loadFromFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("%s file not found", path))
	}
	content, _ := io.ReadAll(file)
	self.yaml = &DbTable{}
	err = yaml.Unmarshal(content, self.yaml)
	if err != nil {
		return err
	}
	fmt.Printf("%+v \n", self.yaml)
	self.yaml.TableMap = make(map[string]Table)

	for _, table := range self.yaml.Table {
		self.yaml.TableMap[table.Name] = table
	}

	return nil
}

func (self *yamlGenerator) SetCustomColumnSaveDir(customColumnSaveDir string) {
	self.customColumnSaveDir = customColumnSaveDir

	if err := os.MkdirAll(self.customColumnSaveDir, os.ModePerm); err != nil {
		panic(err)
	}
}

func (self *yamlGenerator) generateCustomColumn(customColumnType string) error {
	customColumnTemplate := template.CustomColumnTemplate
	customColumnTemplate = strings.Replace(customColumnTemplate, "{{Package}}", strings.TrimRight(path.Base(self.customColumnSaveDir), "/"), 1)
	customColumnTemplate = strings.Replace(customColumnTemplate, "{{CustomStructName}}", customColumnType, -1)

	return os.WriteFile(self.customColumnSaveDir+"/"+strings.ToLower(customColumnType)+".go", []byte(customColumnTemplate), 0640)
}

func (self *yamlGenerator) generateFromTable(table Table) {
	if _, exists := self.generatedTable[table.Name]; exists {
		return
	}

	for _, relate := range table.Relate {
		if tTable, exists := self.yaml.TableMap[relate.Table]; exists {
			self.generateFromTable(tTable)
		} else {
			relateMate := self.gen.GenerateModel(relate.Table)
			self.gen.ApplyBasic(relateMate)
			self.generatedTable[relate.Table] = relateMate.ModelStructName
		}
	}

	//找到所有relate,生成模型
	opt := make([]gen.ModelOpt, len(table.Relate)+len(table.Column))
	for i, table := range table.Relate {
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
		relateConfig := make(field.GormTag)

		if table.ForeignKey != "" {
			relateConfig.Append("foreignKey", table.ForeignKey)
		}
		if table.JoinForeignKey != "" {
			relateConfig.Append("joinForeignKey", table.JoinForeignKey)
		}
		if table.References != "" {
			relateConfig.Append("references", table.References)
		}
		if table.JoinReferences != "" {
			relateConfig.Append("joinReferences", table.JoinReferences)
		}
		if table.Many2many != "" {
			relateConfig.Append("many2many", table.Many2many)
		}
		opt[i] = gen.FieldRelate(fieldType, self.generatedTable[table.Table], self.gen.Data[self.generatedTable[table.Table]].QueryStructMeta, &field.RelateConfig{
			GORMTag: relateConfig,
		})
	}
	//找到column生成自定义column类型
	i := len(table.Relate)
	for name, column := range table.Column {
		err := self.generateCustomColumn(column.CustomColumnType)
		if err != nil {
			panic(err)
		}

		opt[i] = gen.FieldType(name, strings.TrimRight(path.Base(self.customColumnSaveDir), "/")+"."+column.CustomColumnType)
		i++
	}

	relateMate := self.gen.GenerateModel(table.Name, opt...)
	if i > len(table.Relate) {
		pkgs, err := packages.Load(&packages.Config{
			Mode: packages.NeedName,
			Dir:  self.customColumnSaveDir,
		})
		if err != nil {
			panic(err)
		}
		relateMate.ImportPkgPaths = append(relateMate.ImportPkgPaths, "\""+pkgs[0].PkgPath+"\"")
	}
	self.gen.ApplyBasic(relateMate)
	self.generatedTable[table.Name] = relateMate.ModelStructName
}

func (self *yamlGenerator) Generate(opt ...gen.ModelOpt) {
	for _, table := range self.yaml.Table {
		self.generateFromTable(table)
	}
}
