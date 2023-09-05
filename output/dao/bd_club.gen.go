// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/we7coreteam/gorm-gen-yaml/output/entity"
)

func newClub(db *gorm.DB, opts ...gen.DOOption) club {
	_club := club{}

	_club.clubDo.UseDB(db, opts...)
	_club.clubDo.UseModel(&entity.Club{})

	tableName := _club.clubDo.TableName()
	_club.ALL = field.NewAsterisk(tableName)
	_club.ID = field.NewInt32(tableName, "id")
	_club.Name = field.NewString(tableName, "name")
	_club.Logo = field.NewString(tableName, "logo")
	_club.Intro = field.NewString(tableName, "intro")
	_club.Content = field.NewString(tableName, "content")
	_club.WxNickname = field.NewString(tableName, "wx_nickname")
	_club.RealName = field.NewString(tableName, "real_name")
	_club.Mobile = field.NewString(tableName, "mobile")
	_club.ApplicantID = field.NewInt32(tableName, "applicant_id")
	_club.Status = field.NewInt32(tableName, "status")
	_club.AdminID = field.NewInt32(tableName, "admin_id")
	_club.AdminMemo = field.NewString(tableName, "admin_memo")
	_club.Version = field.NewInt32(tableName, "version")
	_club.UserNumber = field.NewInt32(tableName, "user_number")
	_club.NumberLimit = field.NewInt32(tableName, "number_limit")
	_club.BillFee = field.NewFloat64(tableName, "bill_fee")
	_club.NextSettlementDate = field.NewTime(tableName, "next_settlement_date")
	_club.CreateTime = field.NewInt32(tableName, "create_time")
	_club.UpdateTime = field.NewInt32(tableName, "update_time")
	_club.DeleteTime = field.NewInt32(tableName, "delete_time")
	_club.ClubUser = clubHasManyClubUser{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("ClubUser", "entity.ClubUser"),
	}

	_club.User = clubBelongsToUser{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("User", "entity.User"),
		UserOauth: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("User.UserOauth", "entity.UserOauth"),
		},
	}

	_club.fillFieldMap()

	return _club
}

type club struct {
	clubDo

	ALL                field.Asterisk
	ID                 field.Int32
	Name               field.String // 俱乐部名称
	Logo               field.String // LOGO
	Intro              field.String // 简介
	Content            field.String
	WxNickname         field.String  // 微信昵称
	RealName           field.String  // 真实姓名
	Mobile             field.String  // 联系人手机号
	ApplicantID        field.Int32   // 申请人ID
	Status             field.Int32   // 状态:0=待审批,1=已通过,2=已驳回
	AdminID            field.Int32   // 审批人ID
	AdminMemo          field.String  // 审批备注
	Version            field.Int32   // 当前版本：0=免费版 1=付费版
	UserNumber         field.Int32   // 当前成员数量
	NumberLimit        field.Int32   // 人数限额
	BillFee            field.Float64 // 本次账单费用
	NextSettlementDate field.Time    // 下次结算日期
	CreateTime         field.Int32   // 添加时间
	UpdateTime         field.Int32   // 修改时间
	DeleteTime         field.Int32   // 删除时间
	ClubUser           clubHasManyClubUser

	User clubBelongsToUser

	fieldMap map[string]field.Expr
}

func (c club) Table(newTableName string) *club {
	c.clubDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c club) As(alias string) *club {
	c.clubDo.DO = *(c.clubDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *club) updateTableName(table string) *club {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewInt32(table, "id")
	c.Name = field.NewString(table, "name")
	c.Logo = field.NewString(table, "logo")
	c.Intro = field.NewString(table, "intro")
	c.Content = field.NewString(table, "content")
	c.WxNickname = field.NewString(table, "wx_nickname")
	c.RealName = field.NewString(table, "real_name")
	c.Mobile = field.NewString(table, "mobile")
	c.ApplicantID = field.NewInt32(table, "applicant_id")
	c.Status = field.NewInt32(table, "status")
	c.AdminID = field.NewInt32(table, "admin_id")
	c.AdminMemo = field.NewString(table, "admin_memo")
	c.Version = field.NewInt32(table, "version")
	c.UserNumber = field.NewInt32(table, "user_number")
	c.NumberLimit = field.NewInt32(table, "number_limit")
	c.BillFee = field.NewFloat64(table, "bill_fee")
	c.NextSettlementDate = field.NewTime(table, "next_settlement_date")
	c.CreateTime = field.NewInt32(table, "create_time")
	c.UpdateTime = field.NewInt32(table, "update_time")
	c.DeleteTime = field.NewInt32(table, "delete_time")

	c.fillFieldMap()

	return c
}

func (c *club) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *club) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 22)
	c.fieldMap["id"] = c.ID
	c.fieldMap["name"] = c.Name
	c.fieldMap["logo"] = c.Logo
	c.fieldMap["intro"] = c.Intro
	c.fieldMap["content"] = c.Content
	c.fieldMap["wx_nickname"] = c.WxNickname
	c.fieldMap["real_name"] = c.RealName
	c.fieldMap["mobile"] = c.Mobile
	c.fieldMap["applicant_id"] = c.ApplicantID
	c.fieldMap["status"] = c.Status
	c.fieldMap["admin_id"] = c.AdminID
	c.fieldMap["admin_memo"] = c.AdminMemo
	c.fieldMap["version"] = c.Version
	c.fieldMap["user_number"] = c.UserNumber
	c.fieldMap["number_limit"] = c.NumberLimit
	c.fieldMap["bill_fee"] = c.BillFee
	c.fieldMap["next_settlement_date"] = c.NextSettlementDate
	c.fieldMap["create_time"] = c.CreateTime
	c.fieldMap["update_time"] = c.UpdateTime
	c.fieldMap["delete_time"] = c.DeleteTime

}

func (c club) clone(db *gorm.DB) club {
	c.clubDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c club) replaceDB(db *gorm.DB) club {
	c.clubDo.ReplaceDB(db)
	return c
}

type clubHasManyClubUser struct {
	db *gorm.DB

	field.RelationField
}

func (a clubHasManyClubUser) Where(conds ...field.Expr) *clubHasManyClubUser {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a clubHasManyClubUser) WithContext(ctx context.Context) *clubHasManyClubUser {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a clubHasManyClubUser) Session(session *gorm.Session) *clubHasManyClubUser {
	a.db = a.db.Session(session)
	return &a
}

func (a clubHasManyClubUser) Model(m *entity.Club) *clubHasManyClubUserTx {
	return &clubHasManyClubUserTx{a.db.Model(m).Association(a.Name())}
}

type clubHasManyClubUserTx struct{ tx *gorm.Association }

func (a clubHasManyClubUserTx) Find() (result []*entity.ClubUser, err error) {
	return result, a.tx.Find(&result)
}

func (a clubHasManyClubUserTx) Append(values ...*entity.ClubUser) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a clubHasManyClubUserTx) Replace(values ...*entity.ClubUser) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a clubHasManyClubUserTx) Delete(values ...*entity.ClubUser) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a clubHasManyClubUserTx) Clear() error {
	return a.tx.Clear()
}

func (a clubHasManyClubUserTx) Count() int64 {
	return a.tx.Count()
}

type clubBelongsToUser struct {
	db *gorm.DB

	field.RelationField

	UserOauth struct {
		field.RelationField
	}
}

func (a clubBelongsToUser) Where(conds ...field.Expr) *clubBelongsToUser {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a clubBelongsToUser) WithContext(ctx context.Context) *clubBelongsToUser {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a clubBelongsToUser) Session(session *gorm.Session) *clubBelongsToUser {
	a.db = a.db.Session(session)
	return &a
}

func (a clubBelongsToUser) Model(m *entity.Club) *clubBelongsToUserTx {
	return &clubBelongsToUserTx{a.db.Model(m).Association(a.Name())}
}

type clubBelongsToUserTx struct{ tx *gorm.Association }

func (a clubBelongsToUserTx) Find() (result *entity.User, err error) {
	return result, a.tx.Find(&result)
}

func (a clubBelongsToUserTx) Append(values ...*entity.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a clubBelongsToUserTx) Replace(values ...*entity.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a clubBelongsToUserTx) Delete(values ...*entity.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a clubBelongsToUserTx) Clear() error {
	return a.tx.Clear()
}

func (a clubBelongsToUserTx) Count() int64 {
	return a.tx.Count()
}

type clubDo struct{ gen.DO }

type IClubDo interface {
	gen.SubQuery
	Debug() IClubDo
	WithContext(ctx context.Context) IClubDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IClubDo
	WriteDB() IClubDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IClubDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IClubDo
	Not(conds ...gen.Condition) IClubDo
	Or(conds ...gen.Condition) IClubDo
	Select(conds ...field.Expr) IClubDo
	Where(conds ...gen.Condition) IClubDo
	Order(conds ...field.Expr) IClubDo
	Distinct(cols ...field.Expr) IClubDo
	Omit(cols ...field.Expr) IClubDo
	Join(table schema.Tabler, on ...field.Expr) IClubDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IClubDo
	RightJoin(table schema.Tabler, on ...field.Expr) IClubDo
	Group(cols ...field.Expr) IClubDo
	Having(conds ...gen.Condition) IClubDo
	Limit(limit int) IClubDo
	Offset(offset int) IClubDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IClubDo
	Unscoped() IClubDo
	Create(values ...*entity.Club) error
	CreateInBatches(values []*entity.Club, batchSize int) error
	Save(values ...*entity.Club) error
	First() (*entity.Club, error)
	Take() (*entity.Club, error)
	Last() (*entity.Club, error)
	Find() ([]*entity.Club, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Club, err error)
	FindInBatches(result *[]*entity.Club, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*entity.Club) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IClubDo
	Assign(attrs ...field.AssignExpr) IClubDo
	Joins(fields ...field.RelationField) IClubDo
	Preload(fields ...field.RelationField) IClubDo
	FirstOrInit() (*entity.Club, error)
	FirstOrCreate() (*entity.Club, error)
	FindByPage(offset int, limit int) (result []*entity.Club, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IClubDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c clubDo) Debug() IClubDo {
	return c.withDO(c.DO.Debug())
}

func (c clubDo) WithContext(ctx context.Context) IClubDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c clubDo) ReadDB() IClubDo {
	return c.Clauses(dbresolver.Read)
}

func (c clubDo) WriteDB() IClubDo {
	return c.Clauses(dbresolver.Write)
}

func (c clubDo) Session(config *gorm.Session) IClubDo {
	return c.withDO(c.DO.Session(config))
}

func (c clubDo) Clauses(conds ...clause.Expression) IClubDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c clubDo) Returning(value interface{}, columns ...string) IClubDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c clubDo) Not(conds ...gen.Condition) IClubDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c clubDo) Or(conds ...gen.Condition) IClubDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c clubDo) Select(conds ...field.Expr) IClubDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c clubDo) Where(conds ...gen.Condition) IClubDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c clubDo) Order(conds ...field.Expr) IClubDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c clubDo) Distinct(cols ...field.Expr) IClubDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c clubDo) Omit(cols ...field.Expr) IClubDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c clubDo) Join(table schema.Tabler, on ...field.Expr) IClubDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c clubDo) LeftJoin(table schema.Tabler, on ...field.Expr) IClubDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c clubDo) RightJoin(table schema.Tabler, on ...field.Expr) IClubDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c clubDo) Group(cols ...field.Expr) IClubDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c clubDo) Having(conds ...gen.Condition) IClubDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c clubDo) Limit(limit int) IClubDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c clubDo) Offset(offset int) IClubDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c clubDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IClubDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c clubDo) Unscoped() IClubDo {
	return c.withDO(c.DO.Unscoped())
}

func (c clubDo) Create(values ...*entity.Club) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c clubDo) CreateInBatches(values []*entity.Club, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c clubDo) Save(values ...*entity.Club) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c clubDo) First() (*entity.Club, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Club), nil
	}
}

func (c clubDo) Take() (*entity.Club, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Club), nil
	}
}

func (c clubDo) Last() (*entity.Club, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Club), nil
	}
}

func (c clubDo) Find() ([]*entity.Club, error) {
	result, err := c.DO.Find()
	return result.([]*entity.Club), err
}

func (c clubDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Club, err error) {
	buf := make([]*entity.Club, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c clubDo) FindInBatches(result *[]*entity.Club, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c clubDo) Attrs(attrs ...field.AssignExpr) IClubDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c clubDo) Assign(attrs ...field.AssignExpr) IClubDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c clubDo) Joins(fields ...field.RelationField) IClubDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c clubDo) Preload(fields ...field.RelationField) IClubDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c clubDo) FirstOrInit() (*entity.Club, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Club), nil
	}
}

func (c clubDo) FirstOrCreate() (*entity.Club, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Club), nil
	}
}

func (c clubDo) FindByPage(offset int, limit int) (result []*entity.Club, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c clubDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c clubDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c clubDo) Delete(models ...*entity.Club) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *clubDo) withDO(do gen.Dao) *clubDo {
	c.DO = *do.(*gen.DO)
	return c
}