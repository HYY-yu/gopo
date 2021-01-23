package gopo

import (
	"testing"

	"github.com/HYY-yu/gopo/log"
	"go.uber.org/zap"
)

func TestGen_Build(t *testing.T) {
	// init log
	logger, _ := zap.NewDevelopment()
	log.L = logger.Sugar()
	defer logger.Sync()

	g := New()

	err := g.Build(&Config{
		Dir:          "./",
		FileName:     "test_po.go",
		OutDir:       "./log",
		DBPrefix:     "",
		TablePrefix:  "",
		NameStrategy: SnakeCase,
		UseTag:       true,
	})

	if err != nil {
		log.L.Error(err)
	}
}
