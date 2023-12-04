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
	columnOptionSaveDir string
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
	Config   Config  `yaml:"config"`
	Table    []Table `yaml:"relation"`
	TableMap map[string]*Table
}

type Table struct {
	Flag   uint
	Name   string            `yaml:"table"`
	Relate []Relate          `yaml:"relate"`
	Column map[string]Column `yaml:"column"`
	Props  map[string]Column `yaml:"props"`
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

type Config struct {
	TagJsonCamel string `yaml:"tag_json_camel"`
}

func (self *yamlGenerator) UseGormGenerator(g *gen.Generator) *yamlGenerator {
	self.gen = g

	self.SetColumnOptionSaveDir(g.Config.ModelPkgPath + "/../accessor")

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
	self.yaml.TableMap = make(map[string]*Table)

	for _, table := range self.yaml.Table {
		ttable := table
		self.yaml.TableMap[table.Name] = &ttable
	}

	return nil
}

func (self *yamlGenerator) SetColumnOptionSaveDir(columnOptionSaveDir string) {
	self.columnOptionSaveDir = columnOptionSaveDir

	if err := os.MkdirAll(self.columnOptionSaveDir, os.ModePerm); err != nil {
		panic(err)
	}
}

func (self *yamlGenerator) generateColumnOption(column Column) error {
	var columnOptionTemplate string
	var exists bool

	if column.Serializer == "json" || column.Serializer == "gob" || column.Serializer == "unixtime" {
		columnOptionTemplate, exists = template.ColumnOptionTemplate["json"]
	} else {
		columnOptionTemplate, exists = template.ColumnOptionTemplate["common"]
	}

	if !exists {
		return errors.New("serializer type not support")
	}
	columnOptionTemplate = strings.Replace(columnOptionTemplate, "{{Package}}", strings.TrimRight(path.Base(self.columnOptionSaveDir), "/"), 1)
	columnOptionTemplate = strings.Replace(columnOptionTemplate, "{{OptionStructName}}", column.Type, -1)

	path := self.columnOptionSaveDir + "/" + CamelCaseToUnderscore(column.Type) + ".go"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.WriteFile(path, []byte(columnOptionTemplate), 0640)
	}
	return nil
}

func (self *yamlGenerator) getTableRelateOpt(table *Table) []gen.ModelOpt {
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

func (self *yamlGenerator) getTableColumnOpt(table *Table) ([]gen.ModelOpt, bool) {
	opt := make([]gen.ModelOpt, 0)
	//找到column生成自定义column类型
	hasOption := false
	for name, column := range table.Column {
		if column.Type != "" {
			if strings.Contains(strings.ToLower(column.Type), "option") {
				if column.Serializer == "json" || column.Serializer == "gob" || column.Serializer == "unixtime" {
					if column.Tag == nil {
						column.Tag = map[string]map[string]string{
							"gorm": {
								"serializer": column.Serializer,
							},
						}
					} else {
						column.Tag["gorm"]["serializer"] = column.Serializer
					}
					// 生成对应的类型文件
					err := self.generateColumnOption(column)
					if err != nil {
						panic(err)
					}
				} else {
					// 自定义生成 Scan Value
					err := self.generateColumnOption(column)
					if err != nil {
						panic(err)
					}
				}
				hasOption = true
				opt = append(opt, gen.FieldType(name, "*"+strings.TrimRight(path.Base(self.columnOptionSaveDir), "/")+"."+column.Type))
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

	for name, column := range table.Props {
		tag := field.Tag{}
		tag.Set(field.TagKeyJson, UnderscoreToCamelCase(name, self.yaml.Config.TagJsonCamel == "upper"))
		opt = append(opt, gen.FieldNew(UnderscoreToCamelCase(name, true), "*"+strings.TrimRight(path.Base(self.columnOptionSaveDir), "/")+"."+column.Type, tag))
		// 自定义生成 Scan Value
		column.Serializer = "common"
		err := self.generateColumnOption(column)
		if err != nil {
			panic(err)
		}
		hasOption = true
	}
	return opt, hasOption
}

func (self *yamlGenerator) generateFromTable(table *Table) {
	if _, exists := self.generatedTable[table.Name]; exists {
		return
	}
	table.Flag = 1

	for _, relate := range table.Relate {
		if tTable, exists := self.yaml.TableMap[relate.Table]; exists && tTable.Flag == 0 {
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
			Dir:  self.columnOptionSaveDir,
		})
		if err != nil {
			panic(err)
		}
		relateMate.ImportPkgPaths = append(relateMate.ImportPkgPaths, "\""+pkgs[0].PkgPath+"\"")
	}
	if _, exists := self.generatedTable[table.Name]; exists {
		delete(self.gen.Data, self.generatedTable[table.Name])
	}
	self.gen.ApplyBasic(relateMate)
	self.generatedTable[table.Name] = relateMate.ModelStructName
}

func (self *yamlGenerator) Generate(opt ...gen.ModelOpt) {
	if self.yaml.Config.TagJsonCamel != "" {
		self.gen.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
			return UnderscoreToCamelCase(columnName, self.yaml.Config.TagJsonCamel == "upper")
		})
	}
	for _, table := range self.yaml.TableMap {
		self.generateFromTable(table)
	}
}
