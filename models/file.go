package models

import (
	"errors"
	"fmt"
	"gopan/dto"
	"gopan/enums"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type File struct {
	Id             int       `orm:"column(file_id);auto" description:"文件id"`
	UserId         int       `orm:"column(user_id)" description:"用户id"`
	FileMd5        string    `orm:"column(file_md5);size(32)" description:"文件md5值"`
	FilePid        int       `orm:"column(file_pid)" description:"文件父id"`
	FileSize       int64     `orm:"column(file_size)" description:"文件大小"`
	FileName       string    `orm:"column(file_name);size(200)" description:"文件名"`
	FileCover      string    `orm:"column(file_cover);size(100)" description:"文件封面"`
	FilePath       string    `orm:"column(file_path);size(100)" description:"文件路径"`
	CreateTime     time.Time `orm:"column(create_time);type(datetime)" description:"创建时间"`
	LastUpdateTime time.Time `orm:"column(last_update_time);type(datetime)" description:"上次更新时间"`
	FolderType     int8      `orm:"column(folder_type)" description:"文件夹类型"`
	FileCategory   int8      `orm:"column(file_category)" description:"文件分类"`
	FileType       int8      `orm:"column(file_type)" description:"文件类型"`
	Status         int8      `orm:"column(status)" description:"文件状态"`
	RecoveryTime   time.Time `orm:"column(recovery_time);type(datetime)" description:"文件恢复时间"`
}

func (t *File) TableName() string {
	return "file"
}

func init() {
	orm.RegisterModel(new(File))
}

// GetFileList 根据条件获取文件列表
func GetFileList(cond *dto.LoadFileList, userId int) (error, []*File) {
	if cond.Category == "all" {
		return GetFileListByPidAndUserIdAndFileName(cond.FilePid, userId, cond.FileName)
	} else {
		categoryCode := enums.GetCodeByCategoryName(cond.Category)
		return GetFileListByPidAndUserIdAndFileNameAndCategory(userId, cond.FileName, categoryCode)
	}
}

// GetFileListByPidAndUserIdAndFileName 根据父目录id和用户id获取文件列表
func GetFileListByPidAndUserIdAndFileName(pid int, userId int, fileName string) (error, []*File) {
	o := orm.NewOrm()
	var fileList []*File
	sql := "select * from file where file_pid = ? and user_id = ? and file_name like ?"
	_, err := o.Raw(sql, pid, userId, "%"+fileName+"%").QueryRows(&fileList)
	return err, fileList
}

// GetFileListByPidAndUserIdAndFileNameAndCategory 根据父目录id和用户id和文件分类获取文件列表
func GetFileListByPidAndUserIdAndFileNameAndCategory(userId int, fileName string, category int8) (error, []*File) {
	o := orm.NewOrm()
	var fileList []*File
	sql := "select * from file where user_id = ? and file_name like ? and file_category = ?"
	_, err := o.Raw(sql, userId, "%"+fileName+"%", category).QueryRows(&fileList)
	return err, fileList
}

// AddFile insert a new File into database and returns
// last inserted Id on success.
func AddFile(m *File) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetFilesByPid 根据父目录id获取文件列表
func GetFilesByPid(userId int, pid int) (error, []*File) {
	o := orm.NewOrm()
	var fileList []*File
	sql := "select file_name, file_id from file where file_pid = ? and user_id = ?"
	_, err := o.Raw(sql, pid, userId).QueryRows(&fileList)
	return err, fileList
}

// GetFileById retrieves File by Id. Returns error if
// Id doesn't exist
func GetFileById(id int) (v *File, err error) {
	o := orm.NewOrm()
	v = &File{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateFileName 更新文件名
func UpdateFileName(file *File) error {
	o := orm.NewOrm()
	_, err := o.Update(file, "file_name", "last_update_time")
	return err
}

// GetFoldersByPid 根据父目录id获取文件夹列表
func GetFoldersByPid(pid int) (error, []*File) {
	o := orm.NewOrm()
	var fileList []*File
	sql := "select * from file where file_pid = ? and folder_type = 1"
	_, err := o.Raw(sql, pid).QueryRows(&fileList)
	return err, fileList
}

// UpdatePidById 更新文件父id
func UpdatePidById(file *File) error {
	o := orm.NewOrm()
	_, err := o.Update(file, "file_pid", "last_update_time")
	return err
}

// GetFilePathById 根据文件id获取文件路径
func GetFilePathById(id int) (error, string) {
	o := orm.NewOrm()
	var file File
	sql := "select file_path from file where file_id = ?"
	err := o.Raw(sql, id).QueryRow(&file)
	return err, file.FilePath
}

// GetAllFile retrieves all File matches certain condition. Returns empty list if
// no records exist
func GetAllFile(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(File))
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

	var l []File
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

// UpdateFile updates File by Id and returns error if
// the record to be updated doesn't exist
func UpdateFileById(m *File) (err error) {
	o := orm.NewOrm()
	v := File{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFile deletes File by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFile(id int) (err error) {
	o := orm.NewOrm()
	v := File{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&File{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
