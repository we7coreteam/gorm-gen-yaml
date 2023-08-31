package yamlgen

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

func TestSelect(t *testing.T) {
	dao.SetDefault(db)
	row, _ := dao.Q.DeclarationLog.Preload(dao.Q.DeclarationLog.Club, dao.Q.DeclarationLog.Club.ClubUser).Last()
	fmt.Printf("%+v", row.Club.ClubUser)

	row1, _ := dao.Q.Club.Preload(dao.Club.User).Where(dao.Club.ID.In(45)).First()
	fmt.Printf("%+v \n", row1)
	fmt.Printf("%+v \n", row1.User)
}

func TestSql(t *testing.T) {
	result := map[string]interface{}{}
	db.Table("bd_user").Take(&result)
	fmt.Printf("%+v \n", result)

	type User struct {
		ID       int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
		Nickname string `gorm:"column:nickname;not null" json:"nickname"` // 昵称
	}

	var resultUser User
	db.First(&resultUser)
	fmt.Printf("%+v \n", resultUser)

	dao.SetDefault(db)
	resultDao, _ := dao.Q.User.First()
	fmt.Printf("%+v \n", resultDao)

	resultDao1, _ := dao.Q.User.Preload(dao.Q.User.UserOauth).First()
	fmt.Printf("%+v \n", resultDao1.UserOauth.Openid)

	dUser := dao.Q.User
	resultDao2, _ := dUser.Select(dUser.Mobile, dUser.Nickname).Where(dUser.ID.In(1, 20, 30)).Find()
	for _, user := range resultDao2 {
		fmt.Printf("%+v \n", user)
	}

}
