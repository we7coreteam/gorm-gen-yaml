# gorm-gen-yaml

根据 yaml 配置文件生成数据库表的 model 文件及 dao 文件。

程序会根据配置文件中声明的表之间的依赖关系，按照顺序进行生成。让你摆锐复杂的依赖关系。

output 目录为根本下面配置生成出的文件。

# 配置

```yaml
relation:
  - table: declaration_log  
    relate:
      - table: club # 关联的表名
        foreign_key: club_id # 关联的外键
        type: has_one # 关联类型 has_one has_many many_many belongs_to
      - table: declaration_gift
        foreign_key: declaration_gift_id
        type: has_one
  - table: club
    relate:
      - table: club_user
        foreign_key: club_id
        type: has_many
      - table: user
        foreign_key: applicant_id
        type: belongs_to
  - table: user
    relate:
      - table: user_oauth
        foreign_key: user_id
        type: has_one
```

# 调用

```go
g := gen.NewGenerator(gen.Config{
    OutPath:      "./output/dao",
    Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
    ModelPkgPath: "./output/entity",
})
g.UseDB(db)
fieldOpts := []gen.ModelOpt{}

NewYamlGenerator("./gen.yaml").UseGormGenerator(g).Generate(fieldOpts...)
//g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
//g.GenerateModel(tableName)
g.Execute()
```
