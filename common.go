package gorm_gen_yaml

import "gorm.io/gen"

type parser interface {
	loadFromFile(path string) error
	Generate(opt ...gen.ModelOpt) []interface{}
}
