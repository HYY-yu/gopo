Gopo
---

[![Build Status](https://travis-ci.com/HYY-yu/gopo.svg?branch=master)](https://travis-ci.com/HYY-yu/gopo)


我们有时候使用gorm或其他orm的时候，一般都会为数据库表建立对应的结构体。
或者使用(https://github.com/xxjwxc/gormt) 这个工具包，自动生成结构体。

```
type SomeTable struct{
        Id       int     `gorm:"primary_key" json:"id"` //记录id
    	Height   float32 `json:"height"`                //身高
    	Weight   float32 `json:"weight"`                //体重
}
```

但是我们要调用表名、列名时候，还是需要手动写下来
```
func FindHeight(id int){
    db.Select("id","height").Find(&SomeTable{},id)
}
func FindWeight(id int){
    db.Select("id","weight").Find(&SomeTable{},id)
}
```

这样即不方便，也不好管理，这个工具就是为了解决这个问题：
```
type SomeTable struct{
        Id       int     `gorm:"primary_key" json:"id"` //记录id
    	Height   float32 `json:"height"`                //身高
    	Weight   float32 `json:"weight"`                //体重
}

// --自动生成
const (
	SomeTableTableName = "cms.some_table"
	
	SomeTableId = "id"
	SomeTableHeight = "height"
	SomeTableWeight = "weight"
)
// --自动生成

// 想要什么字段，直接传进去就是，无需FindHeight、FindWeight
func FindBy(id int, columns ...string){
    db.Select(colums).Find(&SomeTable{},id)
}
```

### 使用方法

本工具支持：`homebrew`、`scoop`直接下载

```
brew install HYY-yu/tap/gopo
```

```
scoop bucket add gopo https://github.com/HYY-yu/scoop-bucket
scoop install gopo/gopo
```

也可在`github release`中下载打好的包

```
 gopo -h
NAME:
   gopo - Automatically generate go struct'field to variable

USAGE:
   gopo [global options] command [command options] [arguments...]

VERSION:
   0.1.8

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir value, -d value           Directory you want to parse (default: "./")
   --filename value, -f value      Which file you want to parse, default is all file in directory
   --out value, -o value           Output directory for all the generated files
   --db_prefix value, -D value     Output database schema
   --table_prefix value, -T value  Output table prefix
   --field value                   Field Naming Strategy like snakecase,camelcase,pascalcase, if use_tag is true, use tag name (default: "snakecase")
   --use_tag                       use tag name by parse go struct, if true, gorm tag is priority than json tag (default: true)
   --debug                         Print debug log or not (default: false)
   --help, -h                      show help (default: false)
   --version, -v                   print the version (default: false)

```
