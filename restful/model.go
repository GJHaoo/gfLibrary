package restful

import (
	"fmt"
	"reflect"

	"strings"

	"codeup.aliyun.com/sevenfifteen/sevenfifteenBoilerplate/go-library/utils"
	"github.com/gogf/gf/util/gconv"
	"gorm.io/gorm"
)

func GenSqlByRes(sql *gorm.DB, res interface{}) *gorm.DB {
	var (
		selects []interface{}
	)
	resType := reflect.TypeOf(res).Elem() //通过反射获取type定义
	for i := 0; i < resType.NumField(); i++ {
		var (
			item                 = resType.Field(i)
			selectTag            = item.Tag.Get("select")
			selectItem, relation = getColumnNameAndRelation(sql, item.Name, selectTag)
		)

		if relation != "" {
			genJoinByRelation(sql, relation)
			selectItem = selectTag + " as " + item.Name
		}

		selects = append(selects, selectItem)
	}
	return sql.Debug().Select(selects[0], selects[1:]...)
}

// 默认是相等的条件，_代表不筛选此字段
func GenSqlByParam(sql *gorm.DB, param interface{}) *gorm.DB {
	var (
		dtoType  = reflect.TypeOf(param).Elem() //通过反射获取type定义
		dtoValue = reflect.ValueOf(param).Elem()
	)

	for i := 0; i < dtoType.NumField(); i++ {
		var (
			itemType             = dtoType.Field(i)
			itemValue            = dtoValue.Field(i).Interface()
			operator             = "=" // 默认值
			whereTag             = itemType.Tag.Get("where")
			whereTagArr          = strings.Split(whereTag, ";")
			columnName, relation = getColumnNameAndRelation(sql, itemType.Name, "") // 默认值
		)
		// 此处是默认所有的Dto都是指针类型或者数组
		if utils.IsZeroOfUnderlyingType(itemValue) {
			continue
		}

		if whereTagArr[0] != "" {
			operator = whereTagArr[0]
		} else {
			continue
		}

		if len(whereTagArr) == 2 {
			columnName, relation = getColumnNameAndRelation(sql, itemType.Name, whereTagArr[1])
			if relation != "" {
				genJoinByRelation(sql, relation)
			}
		}

		genCondition(sql, columnName, operator, itemValue)
	}
	return sql
}

func genCondition(sql *gorm.DB, name, operator string, itemValue interface{}) {
	switch operator {
	case "_":
	case "In":
	case "=":
		sql.Where(name+" = ?", itemValue)
	case ">":
	case "<":
	case "like":
		sql.Where(name+" LIKE ?", "%"+*itemValue.(*string)+"%")
	}
}

func getColumnNameAndRelation(sql *gorm.DB, fieldName string, tag string) (columnName string, relation string) {
	var (
		tableName = sql.NamingStrategy.TableName(reflect.TypeOf(sql.Statement.Model).Elem().Name())
		arr       = strings.Split(tag, ".")
		len       = len(arr)
	)

	// 为空代表没有tag，默认值是结构体的字段名
	if tag == "" {
		columnName = tableName + "." + sql.NamingStrategy.ColumnName("", fieldName)
		return
	}
	// _代表不需要处理这个字段
	if tag == "_" {
		return
	}

	// 长度1代表是自定义字段名
	if len == 1 {
		columnName = tableName + "." + sql.NamingStrategy.ColumnName("", arr[0])
	}
	// 长度2代表是连表字段
	if len == 2 {
		columnName = tag
		relation = arr[0]
	}
	return
}

func genJoinByRelation(sql *gorm.DB, relation string) {
	tableName := sql.NamingStrategy.TableName(reflect.TypeOf(sql.Statement.Model).Elem().Name())
	relationTableName := sql.NamingStrategy.TableName(relation)
	joinName := fmt.Sprintf("LEFT JOIN `%s` `%s` ON `%s`.`%s_id` = `%s`.`id` AND `%s`.`deleted_at` IS NULL", relationTableName, relation, tableName, relationTableName, relation, relation)
	isContains := false
	for _, join := range sql.Statement.Joins {
		// TODO: 支持排除多种join模式
		if join.Name == relation || join.Name == joinName {
			isContains = true
			break
		}
	}
	if !isContains {
		sql.Joins(joinName)
	}
}

func Paginate(req interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var page Page
		gconv.Struct(req, &page)
		if page.Current > 0 && page.PageSize > 0 {
			offset := (page.Current - 1) * page.PageSize
			return db.Offset(offset).Limit(page.PageSize)
		}
		return db
	}
}
