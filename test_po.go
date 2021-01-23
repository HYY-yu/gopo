package gopo

type GrowthRecord struct {
	Id       int     `gorm:"primary_key" json:"id"` //记录id
	BabyId   int     `json:"baby_id"`               //宝宝id
	Height   float32 `json:"height"`                //身高
	Weight   float32 `json:"weight"`                //体重
	Content  string  `json:"content"`               //记录内容
	Imgs     string  `json:"imgs"`                  //图片
	CreateAt int64   `json:"create_at"`             //添加时间
	UpdateAt int64   `json:"update_at"`             //更新时间
}

func (i *GrowthRecord) TableName() string {
	return "growth_record"
}

type GrowthBaby struct {
	Id       int     `gorm:"primary_key" json:"id"` //宝宝id
	Uid      int     `json:"uid"`                   //用户id
	BabyName string  `json:"baby_name"`             //宝宝名字
	Birthday int64   `json:"birthday"`              //宝宝生日
	Sex      int8    `json:"sex"`                   //宝宝性别
	Height   float32 `json:"height"`                //宝宝身高
	Weight   float32 `json:"weight"`                //宝宝体重
	CreateAt int64   `json:"create_at"`             //添加时间
}

func (i *GrowthBaby) TableName() string {
	return "growth_baby"
}
