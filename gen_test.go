package yamlgen

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"path/filepath"
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
	path, _ := filepath.Abs("./test.db")
	db, _ = gorm.Open(sqlite.Open(path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "ims_",
			SingularTable: true,
		},
		Logger: newLogger,
	})
}

func TestParse(t *testing.T) {
	os.RemoveAll("./output/dao")
	os.RemoveAll("entity")

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./output/dao",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		ModelPkgPath: "entity",
	})
	g.UseDB(db)
	fieldOpts := []gen.ModelOpt{}

	NewYamlGenerator("./gen.yaml").UseGormGenerator(g).Generate(fieldOpts...)
	//g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)

	//userModel := g.GenerateModel("user")
	//g.ApplyBasic(userModel)

	//userProfileModel := g.GenerateModel("user_oauth")
	//user := g.GenerateModel("user", gen.FieldRelate(field.HasOne, "UserProfile", userProfileModel, &field.RelateConfig{
	//	GORMTag: field.GormTag{"foreignKey": []string{"user_id"}},
	//}))
	//g.ApplyBasic(user, userProfileModel)

	g.Execute()
}

// 生成完成 dao & model 文件时，删掉注释测试
func TestSelect(t *testing.T) {
	//	dao.SetDefault(db)
	//	tester := assert.New(t)
	//	row, _ := dao.Q.Club.Preload(dao.Q.Club.User, dao.Q.Club.ClubUser).Last()
	//	fmt.Printf("%v \n", row)
	//	fmt.Printf("%v \n", row.User)
	//	tester.Equal(row.Name, "测试俱乐部")
	//	tester.Equal(row.ClubUser[0].ClubID, row.ID)
	//	tester.Equal(row.User.ID, row.ApplicantID)
	//
	//	list, _ := dao.Q.Formula.Preload(dao.Q.Formula.Tag).Find()
	//	for _, formula := range list {
	//		fmt.Printf("%v \n", formula)
	//	}
	//	tester.Equal(len(list[0].Tag), 3)
	//	tester.Equal(len(list[1].Tag), 2)
	//testDeleteClubUser := &entity.ClubUser{
	//	ClubID: 1,
	//	Name:   "测试删除",
	//}
	//dao.ClubUser.Create(testDeleteClubUser)
	//dao.ClubUser.Where(dao.ClubUser.ID.Eq(testDeleteClubUser.ID)).Delete()
}
