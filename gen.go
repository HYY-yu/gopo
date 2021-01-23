package gopo

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/HYY-yu/gopo/log"
)

const (
	// CamelCase indicates using CamelCase strategy for struct field.
	CamelCase = "camelcase"

	// PascalCase indicates using PascalCase strategy for struct field.
	PascalCase = "pascalcase"

	// SnakeCase indicates using SnakeCase strategy for struct field.
	SnakeCase = "snakecase"
)

type Gen struct {
	config *Config

	files map[string]*ast.File

	genFiles map[string]*GenFileInfo
}

func New() *Gen {
	g := &Gen{}
	g.files = make(map[string]*ast.File)
	g.genFiles = make(map[string]*GenFileInfo)

	return g
}

type Config struct {
	Dir          string
	FileName     string
	OutDir       string
	DBPrefix     string
	TablePrefix  string
	NameStrategy string
	UseTag       bool
}

func (g *Gen) Build(config *Config) error {
	if _, err := os.Stat(config.Dir); os.IsNotExist(err) {
		return fmt.Errorf("dir: %s is not exist", config.Dir)
	}
	if len(config.OutDir) > 0 {
		if _, err := os.Stat(config.OutDir); os.IsNotExist(err) {
			if err := os.MkdirAll(config.OutDir, os.ModePerm); err != nil {
				return err
			}
		}
	}

	g.config = config

	log.L.Infow("Generate start ...")
	if err := g.getAllGoFileInfo(); err != nil {
		return err
	}

	if len(g.files) == 0 {
		log.L.Infow("Not files found ,exiting...")
		return nil
	}

	if err := g.parseFile(); err != nil {
		return err
	}

	log.L.Infow("Write to file ...")

	log.L.Debugw("data ... ", "genFile", g.genFiles)
	for fAbsPath, fileInfo := range g.genFiles {
		fDir, fName := filepath.Split(fAbsPath)
		fName = strings.Replace(fName, ".go", "", -1)
		if len(config.OutDir) > 0 {
			fDir = config.OutDir
		}

		genFile, err := os.Create(path.Join(fDir, fName+"_var.go"))
		if err != nil {
			return err
		}

		if err := packageTemplate.Execute(genFile, fileInfo); err != nil {
			return err
		}

		genFile.Close()
	}
	return nil
}

func (g *Gen) parseStruct(typeSpec *ast.TypeSpec) (*GenStructInfo, error) {
	// 结构体名称
	typeName := typeSpec.Name.Name
	tableName := g.NameStrategy(typeName)
	if len(g.config.TablePrefix) > 0 {
		tableName = g.config.TablePrefix + tableName
	}
	if len(g.config.DBPrefix) > 0 {
		tableName = g.config.DBPrefix + "." + tableName
	}

	genStructInfo := &GenStructInfo{
		StructName:    typeName,
		FieldNames:    make([]string, 0),
		MapFieldNames: make([]string, 0),
		MapTableName:  tableName,
	}

	// 遍历结构体Field
	if typeSpecStructType, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range typeSpecStructType.Fields.List {
			// 如果是匿名结构体，取结构体中的字段
			if len(field.Names) == 0 {
				if fieldIdent, ok := field.Type.(*ast.Ident); ok {
					if embTypeSpec, ok := fieldIdent.Obj.Decl.(*ast.TypeSpec); ok {
						genStructInfoEmb, err := g.parseStruct(embTypeSpec)
						if err != nil {
							return nil, err
						}
						genStructInfo.FieldNames = append(genStructInfo.FieldNames, genStructInfoEmb.FieldNames...)
						genStructInfo.MapFieldNames = append(genStructInfo.MapFieldNames, genStructInfoEmb.MapFieldNames...)
					}
				}
			} else {
				// 如果是非匿名结构体
				// 1. 取字段名
				// 2. 取tag
				var useNameStrategyFlag = true
				fieldName := field.Names[0].Name
				jsonTagName := parseJsonTag2JsonName(field.Tag.Value)
				gormTagName := parseGormTag2ColumnName(field.Tag.Value)

				// 如果设置了useTag
				if g.config.UseTag {
					if len(jsonTagName) > 0 {
						fieldName = jsonTagName
						useNameStrategyFlag = false
					} else if len(gormTagName) > 0 {
						fieldName = gormTagName
						useNameStrategyFlag = false
					}
				}

				if useNameStrategyFlag && len(g.config.NameStrategy) > 0 {
					fieldName = g.NameStrategy(fieldName)
				}

				genStructInfo.FieldNames = append(genStructInfo.FieldNames, field.Names[0].Name)
				genStructInfo.MapFieldNames = append(genStructInfo.MapFieldNames, fieldName)
			}
		}
	}
	return genStructInfo, nil
}

func (g *Gen) NameStrategy(in string) (out string) {
	// 切换 fieldName 的命名策略
	switch g.config.NameStrategy {
	case CamelCase:
		out = toLowerCamelCase(in)
	case SnakeCase:
		out = toSnakeCase(in)
	case PascalCase:
		// use struct field name
		out = in
	default:
		out = toSnakeCase(in)
	}
	return
}

func (g *Gen) parseFile() error {
	for path, astFile := range g.files {
		genFileInfo := &GenFileInfo{
			PkgName:     astFile.Name.Name,
			StructInfos: make([]GenStructInfo, 0),
		}

		for _, astDeclaration := range astFile.Decls {
			if generalDeclaration, ok := astDeclaration.(*ast.GenDecl); ok && generalDeclaration.Tok == token.TYPE {
				for _, astSpec := range generalDeclaration.Specs {
					if typeSpec, ok := astSpec.(*ast.TypeSpec); ok {
						structInfo, err := g.parseStruct(typeSpec)
						if err != nil {
							return err
						}
						genFileInfo.StructInfos = append(genFileInfo.StructInfos, *structInfo)
					}
				}
			}
		}

		g.genFiles[path] = genFileInfo
	}
	return nil
}

type GenFileInfo struct {
	PkgName     string
	StructInfos []GenStructInfo
}

type GenStructInfo struct {
	StructName    string
	FieldNames    []string
	MapFieldNames []string

	MapTableName string
}

var packageTemplate = template.Must(template.New("").Parse(`// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by gopo 

package {{.PkgName}}

{{range $index, $struct := .StructInfos}}
const (
	{{$struct.StructName}}TableName = "{{.MapTableName}}"
	{{range $index2, $field := $struct.FieldNames }}
	{{$struct.StructName}}{{$field}} = "{{index $struct.MapFieldNames $index2}}"{{end}}
)
{{end}}
`))

func (g *Gen) getAllGoFileInfo() error {
	return filepath.Walk(g.config.Dir, g.visit)
}

func (g *Gen) visit(path string, f os.FileInfo, err error) error {
	if g.skip(path, f) {
		return nil
	}

	// 解析file
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("ParseFile error: %+v", err)
	}
	absFilePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("filepath Abs error: %+v", err)
	}

	log.L.Debug(absFilePath)
	g.files[absFilePath] = astFile
	return nil
}

func (g *Gen) skip(path string, f os.FileInfo) bool {
	// exclude vendor
	if f.IsDir() && f.Name() == "vendor" {
		return true
	}
	// exclude all hidden folder
	if f.IsDir() && len(f.Name()) > 1 && f.Name()[0] == '.' {
		return true
	}

	// only .go
	if ext := filepath.Ext(path); ext != ".go" {
		return true
	}

	// only filename if it's set
	if len(g.config.FileName) > 0 {
		if f.Name() != g.config.FileName {
			log.L.Debugf("skip file:%s by filename: %s", f.Name(), g.config.FileName)
			return true
		}
	}
	return false
}
