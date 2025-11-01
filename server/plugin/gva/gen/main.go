package main

import (
	"gorm.io/gen"
	"path/filepath"
)

func main() {
	g := gen.NewGenerator(gen.Config{OutPath: filepath.Join("..", "model", "dao"), Mode: gen.WithoutContext | gen.WithDefaultQuery})
	g.ApplyBasic()
	g.Execute()
}
