package gorm_gen_yaml

import (
	"fmt"
	"github.com/we7coreteam/gorm-gen-yaml/output/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"testing"
	"time"
)

var db *gorm.DB

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	dsn := "root:123456@tcp(172.16.1.198:3306)/paidashen?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "bd_",
			SingularTable: true,
		},
		Logger: newLogger,
	})
}

func TestParse(t *testing.T) {
	os.RemoveAll("./output/dao")
	os.RemoveAll("./output/entity")

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./output/dao",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		ModelPkgPath: "./output/entity",
	})
	g.UseDB(db)
	fieldOpts := []gen.ModelOpt{}
	fmt.Println(db.Config.NamingStrategy.SchemaName("abc_def"))

	NewYamlGenerator("./gen.yaml").UseGormGenerator(g).Generate(fieldOpts...)
	//g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
	//g.GenerateModel(tableName)
	g.Execute()
}

func TestSelect(t *testing.T) {
	dao.SetDefault(db)
	row, _ := dao.Q.DeclarationLog.Preload(dao.Q.DeclarationLog.Club, dao.Q.DeclarationLog.Club.ClubUser).Last()
	fmt.Printf("%+v", row.Club.ClubUser)

	row1, _ := dao.Q.Club.Preload(dao.Club.User).Where(dao.Club.ID.In(45)).First()
	fmt.Printf("%+v \n", row1)
	fmt.Printf("%+v \n", row1.User)
}
