# gorm-gen-yaml

根据 yaml 配置文件生成数据库表的 model 文件及 dao 文件。

程序会根据配置文件中声明的表之间的依赖关系，按照顺序进行生成。让你摆锐复杂的依赖关系。

output 目录为根据下面配置生成出的文件。

# 配置

<a href="https://github.com/we7coreteam/gorm-gen-yaml/blob/main/gen.yaml">查看配置文件</a>

# 调用

```go
g := gen.NewGenerator(gen.Config{
    OutPath:      "./output/dao",
    Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
    ModelPkgPath: "./output/entity",
})
g.UseDB(db)
fieldOpts := []gen.ModelOpt{}
yamlgen.NewYamlGenerator("./gen.yaml").UseGormGenerator(g).Generate(fieldOpts...)
//g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
//g.GenerateModel(tableName)
g.Execute()
```

# 测试

运行 gen_test.go 文件中的 TestParse 方法。生成完成 dao 和 model 文件后，可以删除 TestSelect 注释进行测试

# 交流群

QQ群：364768550

<img src="https://s2.loli.net/2024/06/03/uMjYwCWmVPaRSUt.png" width="300" >
