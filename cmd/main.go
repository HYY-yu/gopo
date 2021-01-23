package main

import (
	"os"

	"github.com/HYY-yu/gopo"
	"github.com/HYY-yu/gopo/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	version = ""
)

// gopo
// gopo -d .../po.go --out ../po_var.go --db_prefix cms --table_prefix t_ --field camelcase
func main() {
	// init log
	logger, _ := zap.NewProduction()
	log.L = logger.Sugar()
	defer logger.Sync()

	app := cli.NewApp()
	app.Name = "gopo"
	app.Version = version
	app.Usage = "Automatically generate go struct'field to variable"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "dir",
			Aliases: []string{"d"},
			Value:   "./",
			Usage:   "Directory you want to parse",
		},
		&cli.StringFlag{
			Name:    "filename",
			Aliases: []string{"f"},
			Usage:   "Which file you want to parse, default is all file in directory",
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Usage:   "Output directory for all the generated files",
		},
		&cli.StringFlag{
			Name:    "db_prefix",
			Aliases: []string{"D"},
			Usage:   "Output database schema",
		},
		&cli.StringFlag{
			Name:    "table_prefix",
			Aliases: []string{"T"},
			Usage:   "Output table prefix",
		},
		&cli.StringFlag{
			Name:  "field",
			Value: "snakecase",
			Usage: "Field Naming Strategy like snakecase,camelcase,pascalcase, if use_tag is true, use tag name",
		},
		&cli.BoolFlag{
			Name:  "use_tag",
			Value: true,
			Usage: "use tag name by parse go struct, if true, gorm tag is priority than json tag",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Value: false,
			Usage: "Print debug log or not",
		},
	}

	app.Action = func(c *cli.Context) error {
		dir := c.String("dir")
		filename := c.String("filename")
		out := c.String("out")
		dbPrefix := c.String("db_prefix")
		tablePrefix := c.String("table_prefix")
		field := c.String("field")
		useTag := c.Bool("use_tag")
		debug := c.Bool("debug")
		if debug {
			logger, _ = zap.NewDevelopment()
			log.L = logger.Sugar()
		}

		switch field {
		case gopo.CamelCase, gopo.SnakeCase, gopo.PascalCase:
		default:
			return errors.Errorf("not supported %s propertyStrategy", field)
		}

		err := gopo.New().Build(&gopo.Config{
			Dir:          dir,
			FileName:     filename,
			OutDir:       out,
			DBPrefix:     dbPrefix,
			TablePrefix:  tablePrefix,
			NameStrategy: field,
			UseTag:       useTag,
		})
		if err != nil {
			log.L.Errorw(err.Error())
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.L.Fatal(err)
	}
}
