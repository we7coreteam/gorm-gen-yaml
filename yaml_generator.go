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

type YamlGenerator struct {
	yaml                *DbTable
	gen                 *gen.Generator
	generatedTable      map[string]string
	columnOptionSaveDir string
}

func NewYamlGenerator(path string) *YamlGenerator {
	obj := &YamlGenerator{}
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

func (y *YamlGenerator) UseGormGenerator(g *gen.Generator) *YamlGenerator {
	y.gen = g

	y.SetColumnOptionSaveDir(g.Config.OutPath + "/../accessor")

	return y
}

func (y *YamlGenerator) loadFromFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("%s file not found", path))
	}
	content, _ := io.ReadAll(file)
	y.yaml = &DbTable{}
	err = yaml.Unmarshal(content, y.yaml)
	if err != nil {
		return err
	}
	y.yaml.TableMap = make(map[string]*Table)

	for _, table := range y.yaml.Table {
		t := table
		y.yaml.TableMap[table.Name] = &t
	}

	return nil
}

func (y *YamlGenerator) SetColumnOptionSaveDir(columnOptionSaveDir string) {
	y.columnOptionSaveDir = columnOptionSaveDir

	if err := os.MkdirAll(y.columnOptionSaveDir, os.ModePerm); err != nil {
		panic(err)
	}
}

func (y *YamlGenerator) generateColumnOption(column Column) error {
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
	columnOptionTemplate = strings.Replace(columnOptionTemplate, "{{Package}}", strings.TrimRight(path.Base(y.columnOptionSaveDir), "/"), 1)
	columnOptionTemplate = strings.Replace(columnOptionTemplate, "{{OptionStructName}}", column.Type, -1)

	p := y.columnOptionSaveDir + "/" + CamelCaseToUnderscore(column.Type) + ".go"
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return os.WriteFile(p, []byte(columnOptionTemplate), 0640)
	}
	return nil
}

func (y *YamlGenerator) getTableRelateOpt(table *Table) []gen.ModelOpt {
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
		opt[i] = gen.FieldRelate(fieldType, y.generatedTable[table.Table], y.gen.Data[y.generatedTable[table.Table]].QueryStructMeta, &field.RelateConfig{
			GORMTag:       relateConfig,
			RelatePointer: relatePointer,
		})
	}

	return opt
}

func (y *YamlGenerator) getTableColumnOpt(table *Table) ([]gen.ModelOpt, bool) {
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
					err := y.generateColumnOption(column)
					if err != nil {
						panic(err)
					}
				} else {
					// 自定义生成 Scan Value
					err := y.generateColumnOption(column)
					if err != nil {
						panic(err)
					}
				}
				hasOption = true
				opt = append(opt, gen.FieldType(name, "*"+strings.TrimRight(path.Base(y.columnOptionSaveDir), "/")+"."+column.Type))
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
		tag.Set(field.TagKeyJson, UnderscoreToCamelCase(name, y.yaml.Config.TagJsonCamel == "upper"))
		opt = append(opt, gen.FieldNew(UnderscoreToCamelCase(name, true), "*"+strings.TrimRight(path.Base(y.columnOptionSaveDir), "/")+"."+column.Type, tag))
		// 自定义生成 Scan Value
		column.Serializer = "common"
		err := y.generateColumnOption(column)
		if err != nil {
			panic(err)
		}
		hasOption = true
	}
	return opt, hasOption
}

func (y *YamlGenerator) generateFromTable(table *Table, opt ...gen.ModelOpt) {
	if _, exists := y.generatedTable[table.Name]; exists {
		return
	}
	table.Flag = 1

	for _, relate := range table.Relate {
		if tTable, exists := y.yaml.TableMap[relate.Table]; exists && tTable.Flag == 0 {
			y.generateFromTable(tTable, opt...)
		} else {
			relateMate := y.gen.GenerateModel(relate.Table, opt...)
			y.gen.ApplyBasic(relateMate)
			y.generatedTable[relate.Table] = relateMate.ModelStructName
		}
	}

	//找到所有relate,生成模型
	relateOpt := y.getTableRelateOpt(table)
	columnOpt, hasOption := y.getTableColumnOpt(table)
	if opt == nil {
		opt = make([]gen.ModelOpt, 0)
	}
	opt = append(opt, relateOpt...)
	opt = append(opt, columnOpt...)
	relateMate := y.gen.GenerateModel(table.Name, opt...)
	if hasOption {
		pkgs, err := packages.Load(&packages.Config{
			Mode: packages.NeedName,
			Dir:  y.columnOptionSaveDir,
		})
		if err != nil {
			panic(err)
		}
		relateMate.ImportPkgPaths = append(relateMate.ImportPkgPaths, "\""+pkgs[0].PkgPath+"\"")
	}
	if _, exists := y.generatedTable[table.Name]; exists {
		delete(y.gen.Data, y.generatedTable[table.Name])
	}
	y.gen.ApplyBasic(relateMate)
	y.generatedTable[table.Name] = relateMate.ModelStructName
}

func (y *YamlGenerator) Generate(opt ...gen.ModelOpt) {
	if y.yaml.Config.TagJsonCamel != "" {
		y.gen.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
			return UnderscoreToCamelCase(columnName, y.yaml.Config.TagJsonCamel == "upper")
		})
	}
	for _, table := range y.yaml.TableMap {
		y.generateFromTable(table, opt...)
	}
}
