package yamlgen

import (
	"errors"
	"fmt"
	"github.com/we7coreteam/gorm-gen-yaml/template"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/schema"
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
	Type       string                       `yaml:"type"`
	Serializer string                       `yaml:"serializer"`
	Tag        map[string]map[string]string `yaml:"tag"`
	Comment    string                       `yaml:"comment"`
	Rename     string                       `yaml:"rename"`
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

	self.SetCustomColumnSaveDir(g.Config.ModelPkgPath + "/../accessor")

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

func (self *yamlGenerator) generateCustomColumn(column Column) error {
	if column.Serializer == "" {
		column.Serializer = "common"
	}
	customColumnTemplate, exists := template.CustomColumnTemplate[column.Serializer]
	if !exists {
		return errors.New("serializer type not support")
	}
	customColumnTemplate = strings.Replace(customColumnTemplate, "{{Package}}", strings.TrimRight(path.Base(self.customColumnSaveDir), "/"), 1)
	customColumnTemplate = strings.Replace(customColumnTemplate, "{{CustomStructName}}", column.Type, -1)

	return os.WriteFile(self.customColumnSaveDir+"/"+schema.NamingStrategy{}.TableName(column.Type)+".gen.go", []byte(customColumnTemplate), 0640)
}

func (self *yamlGenerator) getTableRelateOpt(table Table) []gen.ModelOpt {
	opt := make([]gen.ModelOpt, len(table.Relate))
	for i, table := range table.Relate {
		relatePointer := false
		var fieldType field.RelationshipType
		switch table.Type {
		case "has_one":
			fieldType = field.HasOne
			relatePointer = true
		case "has_many":
			fieldType = field.HasMany
		case "many_many":
			fieldType = field.Many2Many
		case "belongs_to":
			fieldType = field.BelongsTo
			relatePointer = true
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
			GORMTag:       relateConfig,
			RelatePointer: relatePointer,
		})
	}

	return opt
}

func (self *yamlGenerator) getTableColumnOpt(table Table) ([]gen.ModelOpt, bool) {
	opt := make([]gen.ModelOpt, 0)
	//找到column生成自定义column类型
	hasOption := false
	for name, column := range table.Column {
		if column.Type != "" {
			if strings.Contains(strings.ToLower(column.Type), "option") {
				err := self.generateCustomColumn(column)
				if err != nil {
					panic(err)
				}
				hasOption = true
				opt = append(opt, gen.FieldType(name, strings.TrimRight(path.Base(self.customColumnSaveDir), "/")+"."+column.Type))
			} else {
				opt = append(opt, gen.FieldType(name, column.Type))
			}
		}
		if column.Tag != nil {
			for tagType, tags := range column.Tag {
				ttags := tags
				if tagType == "gorm" {
					opt = append(opt, gen.FieldGORMTag(name, func(tag field.GormTag) field.GormTag {
						for tagName, val := range ttags {
							tag = tag.Set(tagName, val)
						}
						return tag
					}))
				} else {
					tTagType := tagType
					opt = append(opt, gen.FieldTag(name, func(tag field.Tag) field.Tag {
						tagStr := ""
						for tagName, val := range ttags {
							if val == "" {
								tagStr += tagName + ";"
							} else {
								tagStr += tagName + ":" + val + ";"
							}
						}
						return tag.Set(tTagType, strings.TrimRight(tagStr, ";"))
					}))
				}
			}
		}
		if column.Rename != "" {
			opt = append(opt, gen.FieldRename(name, column.Rename))
		}
		if column.Comment != "" {
			opt = append(opt, gen.FieldComment(name, column.Comment))
		}
	}

	return opt, hasOption
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
	relateOpt := self.getTableRelateOpt(table)
	columnOpt, hasOption := self.getTableColumnOpt(table)
	opt := append(relateOpt, columnOpt...)
	relateMate := self.gen.GenerateModel(table.Name, opt...)
	if hasOption {
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
