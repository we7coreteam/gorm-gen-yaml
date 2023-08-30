package gorm_gen_yaml

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	os.RemoveAll("./output/dao")
	os.RemoveAll("./output/entity")

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./output/dao",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		ModelPkgPath: "./output/entity",
	})
	dsn := "root:123456@tcp(172.16.1.198:3306)/paidashen?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   "bd_",
		//	SingularTable: true,
		//},
	})
	g.UseDB(db)

	yamlGen := NewYamlGenerator("./gen.yaml").UseGormGenerator(g)
	fieldOpts := []gen.ModelOpt{}
	g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
	g.ApplyBasic(yamlGen.Generate(fieldOpts...)...)
	//g.GenerateModel(tableName)
	g.Execute()
}
