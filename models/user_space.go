package models

import (
	"errors"
	"fmt"
	"gopan/constants"
	"gopan/utils"
	"reflect"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type UserSpace struct {
	Id         int   `orm:"column(space_id);auto" description:"空间id"`
	UserId     int   `orm:"column(user_id)" description:"用户id"`
	UseSpace   int64 `orm:"column(use_space)" description:"已用空间"`
	TotalSpace int64 `orm:"column(total_space);null" description:"总空间"`
}

func (t *UserSpace) TableName() string {
	return "user_space"
}

func init() {
	orm.RegisterModel(new(UserSpace))
}

// AddUserSpace insert a new UserSpace into database and returns
// last inserted Id on success.
func AddUserSpace(m *UserSpace) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserSpaceByUserId retrieves UserSpace by UserId.
// Returns error if Id doesn't exist
func GetUserSpaceByUserId(user *UserSpace) error {
	spaceFromRedis, err := utils.GetFromRedis[UserSpace](constants.USER_SPACE_PREFIX_USERID + strconv.Itoa(user.UserId))
	if err == nil {
		*user = *spaceFromRedis
		return nil
	} else {
		o := orm.NewOrm()
		err = o.Read(user, "user_id")
		utils.SetToRedis[UserSpace](constants.USER_SPACE_PREFIX_USERID+strconv.Itoa(user.UserId), user, constants.EXPIRE_NOT)
		return err
	}
}

// UpdateUserSpace 更新用户空间
func UpdateUserSpace(user *UserSpace) error {
	o := orm.NewOrm()
	_, err := o.Update(user, "use_space", "total_space")
	if err == nil {
		utils.SetToRedis[UserSpace](constants.USER_SPACE_PREFIX_USERID+strconv.Itoa(user.UserId), user, constants.EXPIRE_NOT)
	}
	return err
}

// UpdateUserSpaceWithOrm 更新用户空间
func UpdateUserSpaceWithOrm(user *UserSpace, o orm.TxOrmer) error {
	_, err := o.Update(user, "use_space", "total_space")
	if err == nil {
		utils.SetToRedis[UserSpace](constants.USER_SPACE_PREFIX_USERID+strconv.Itoa(user.UserId), user, constants.EXPIRE_NOT)
	}
	return err
}

// GetUserSpaceById retrieves UserSpace by Id. Returns error if
// Id doesn't exist
func GetUserSpaceById(userId int) (v *UserSpace, err error) {
	o := orm.NewOrm()
	v = &UserSpace{UserId: userId}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserSpace retrieves all UserSpace matches certain condition. Returns empty list if
// no records exist
func GetAllUserSpace(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserSpace))
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

	var l []UserSpace
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

// UpdateUserSpace updates UserSpace by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserSpaceById(m *UserSpace) (err error) {
	o := orm.NewOrm()
	v := UserSpace{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserSpace deletes UserSpace by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserSpace(id int) (err error) {
	o := orm.NewOrm()
	v := UserSpace{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserSpace{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
