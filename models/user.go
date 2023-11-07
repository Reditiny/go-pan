package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id            int       `orm:"column(user_id);auto" description:"用户id"`
	NickName      string    `orm:"column(nick_name);size(20)" description:"昵称"`
	Email         string    `orm:"column(email);size(50)" description:"邮箱"`
	QqOpenId      string    `orm:"column(qq_open_id);size(35);null" description:"qq登陆id"`
	QqAvatar      string    `orm:"column(qq_avatar);size(150);null" description:"qq头像"`
	Password      string    `orm:"column(password);size(50)" description:"密码"`
	JoinTime      time.Time `orm:"column(join_time);type(datetime)" description:"创建时间"`
	LastLoginTime time.Time `orm:"column(last_login_time);type(datetime)" description:"上次登陆时间"`
	Status        int8      `orm:"column(status)" description:"账号状态 0:禁用 1:正常"`
}

func (t *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserByEmailAndPassword retrieves User by email and password.
// Returns error if Id doesn't exist
func GetUserByEmailAndPassword(user *User) error {
	o := orm.NewOrm()
	return o.Read(user, "email", "password")
}

// GetUserByEmail retrieves User by email.
// Returns error if Id doesn't exist
func GetUserByEmail(user *User) error {
	o := orm.NewOrm()
	return o.Read(user, "email")
}

// UpdateUserPassword updates User password by Id
func UpdateUserPassword(user *User) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(user, "password")
	return
}

// GetUserByNickname retrieves User by nickname.
// Returns error if Id doesn't exist
func GetUserByNickname(user *User) error {
	o := orm.NewOrm()
	return o.Read(user, "nick_name")
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateUserAvatarById updates User avatar by Id
func UpdateUserAvatarById(user *User) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(user, "qq_avatar")
	return
}

// GetAllUser retrieves all User matches certain condition. Returns empty list if
// no records exist
func GetAllUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
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

	var l []User
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

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
