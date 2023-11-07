package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Share struct {
	Id         int       `orm:"column(share_id);auto"`
	FileId     int       `orm:"column(file_id)" description:"文件id"`
	UserId     int       `orm:"column(user_id)" description:"用户id"`
	ValidType  int8      `orm:"column(valid_type)"`
	Code       string    `orm:"column(code);size(5)"`
	ShowCount  int       `orm:"column(show_count)"`
	ExpireTime time.Time `orm:"column(expire_time);type(datetime)" description:"过期时间"`
}

func (t *Share) TableName() string {
	return "share"
}

func init() {
	orm.RegisterModel(new(Share))
}

// AddShare insert a new Share into database and returns
// last inserted Id on success.
func AddShare(m *Share) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetShareByUserId 根据用户id获取分享列表
func GetShareByUserId(userId int) (error, []*Share) {
	o := orm.NewOrm()
	var shareList []*Share
	sql := "select * from share where user_id = ?"
	_, err := o.Raw(sql, userId).QueryRows(&shareList)
	return err, shareList
}

// GetShareById retrieves Share by Id. Returns error if
// Id doesn't exist
func GetShareById(id int) (v *Share, err error) {
	o := orm.NewOrm()
	v = &Share{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllShare retrieves all Share matches certain condition. Returns empty list if
// no records exist
func GetAllShare(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Share))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Share
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateShare updates Share by Id and returns error if
// the record to be updated doesn't exist
func UpdateShareById(m *Share) (err error) {
	o := orm.NewOrm()
	v := Share{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteShare deletes Share by Id and returns error if
// the record to be deleted doesn't exist
func DeleteShare(id int) (err error) {
	o := orm.NewOrm()
	v := Share{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Share{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
