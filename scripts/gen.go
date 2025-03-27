package main

import (
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./gen",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	// https://github.com/go-gorm/rawsql/blob/master/tests/gen_test.go
	gormdb, _ := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: []string{
			"./scripts//sql",
		},
	}))
	g.UseDB(gormdb) // reuse your gorm db

	g.ApplyBasic(
		// Generate structs from all tables of current database
		g.GenerateAllTable()...,
	)
	// Generate the code
	g.Execute()
}
