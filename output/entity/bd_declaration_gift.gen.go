// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

const TableNameDeclarationGift = "bd_declaration_gift"

// DeclarationGift mapped from table <bd_declaration_gift>
type DeclarationGift struct {
	ID                    int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ClubID                int32   `gorm:"column:club_id;not null;comment:俱乐部ID" json:"club_id"`                                                 // 俱乐部ID
	DeclarationTypeID     int32   `gorm:"column:declaration_type_id;not null;comment:报单类型ID" json:"declaration_type_id"`                        // 报单类型ID
	Name                  string  `gorm:"column:name;not null;comment:名称" json:"name"`                                                          // 名称
	Icon                  string  `gorm:"column:icon;not null;comment:图标" json:"icon"`                                                          // 图标
	Price                 float64 `gorm:"column:price;not null;default:0.00;comment:价格" json:"price"`                                           // 价格
	Commission            float64 `gorm:"column:commission;not null;default:0.00;comment:抽成" json:"commission"`                                 // 抽成
	Rebate                float64 `gorm:"column:rebate;not null;default:0.00;comment:返利" json:"rebate"`                                         // 返利
	Weigh                 int32   `gorm:"column:weigh;not null;comment:权重" json:"weigh"`                                                        // 权重
	Description           string  `gorm:"column:description;not null;comment:描述" json:"description"`                                            // 描述
	Status                string  `gorm:"column:status;not null;comment:状态" json:"status"`                                                      // 状态
	ConvertType           string  `gorm:"column:convert_type;not null;default:1;comment:折算方式:1=不折算,2=固定单价,3=陪玩单价,4=折算成0" json:"convert_type"`   // 折算方式:1=不折算,2=固定单价,3=陪玩单价,4=折算成0
	ConvertPrice          float64 `gorm:"column:convert_price;not null;default:0.00;comment:固定折算价格" json:"convert_price"`                       // 固定折算价格
	ConvertRemoveTailType string  `gorm:"column:convert_remove_tail_type;not null;comment:去尾方式:1=忽略小数点,2=四舍五入" json:"convert_remove_tail_type"` // 去尾方式:1=忽略小数点,2=四舍五入
	CreateTime            int32   `gorm:"column:create_time;not null;comment:添加时间" json:"create_time"`                                          // 添加时间
	UpdateTime            int32   `gorm:"column:update_time;not null;comment:修改时间" json:"update_time"`                                          // 修改时间
	DeleteTime            int32   `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                                                   // 删除时间
}

// TableName DeclarationGift's table name
func (*DeclarationGift) TableName() string {
	return TableNameDeclarationGift
}