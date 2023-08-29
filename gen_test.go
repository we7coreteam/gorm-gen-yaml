package gorm_gen_yaml

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestParse(t *testing.T) {
	at := assert.New(t)

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./output/dao",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		ModelPkgPath: "./output/entity",
	})
	dsn := "root:123456@tcp(127.0.0.1:3306)/paidashen?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   "bd_",
		//	SingularTable: true,
		//},
	})
	g.UseDB(db)

	parser, err := NewYamlGenerate(g, "./gen.yaml")
	at.NoError(err)

	fieldOpts := []gen.ModelOpt{}
	//g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
	g.ApplyBasic(parser.Generate(fieldOpts...)...)
	//g.GenerateModel(tableName)
	g.Execute()
}
