package models

import (
	"database/sql"
	"time"

	"<%= myAppPath %>/utils"

	_ "github.com/SAP/go-hdb/driver"
	"gopkg.in/mgo.v2/bson"
	//"fmt"
)

type User struct {
	Id         bson.ObjectId `bson:"_id"`
	EmployeeNo string        `bson:"ID"`
	Name       string        `bson:"eName"`
	// Info       LoginerInfo   `bson:"Info"`
}

type Reader struct {
	Id         bson.ObjectId `bson:"_id"`
	EmployeeNo string        `bson:"ID"`
	Name       string        `bson:"eName"`
	DepList    []ReaderDep   `bson:"mfDepList"`
}

type ReaderDep struct {
	DepId   string `bson:"mfDepId"`
	OrderNo int32  `bson:"orderByNum"`
}

//总部人事结构(用于显示)
type MFStruct struct {
	//自动生成ID
	Id         bson.ObjectId `bson:"_id"`
	Pid        string        `bson:"pId"`
	Rank       string        `bson:"rank"`
	ElementId  string        `bson:"elementId"`
	Name       string
	SystemName string
	Children   []*MFStruct
}

type MFElement struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

//片区人事结构(用于显示)
type AStructs struct {
	//自动生成ID
	Id         bson.ObjectId `bson:"_id"`
	Pid        string        `bson:"Pid"`
	Rank       string        `bson:"rank"`
	ElementId  string        `bson:"ElementId"`
	Name       string
	SystemName string
	Children   []*AStructs
}

type AElement struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"Name"`
}

type MFUserTag struct {
	//自动生成ID
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

func GetBrandList() (list []*ShowBrand, err error) {

	conn, _ := sql.Open(utils.HANA_DRIVER, utils.HANA_DNS)
	defer conn.Close()

	var sql = `SELECT BEHVO,BVTXT FROM SLTREP.T144T WHERE MANDT='810' AND SPRAS='1'`
	rows, err := conn.Query(sql)

	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var temp ShowBrand
		arr := []interface{}{
			&temp.Code, &temp.Name,
		}
		if err = rows.Scan(arr...); err != nil {
			return
		} else {
			list = append(list, &temp)
		}
	}

	return
}

//用户权限
type UserPermission struct {
	Id             bson.ObjectId `bson:"_id"`
	EmployeeNo     string        `bson:"employeeNo"`
	PermissionCode string        `bson:"permissionCode"`
}
type UserPermissionDB struct {
	EmployeeNo     string `bson:"employeeNo"`
	PermissionCode string `bson:"permissionCode"`
}

type UserCompany struct {
	Id          bson.ObjectId `bson:"_id"`
	EmployeeNo  string        `bson:"employeeNo"`
	CompanyCode string        `bson:"companyCode"`
}
type UserCompanyDB struct {
	EmployeeNo  string `bson:"employeeNo"`
	CompanyCode string `bson:"companyCode"`
}

type UserBrand struct {
	Id         bson.ObjectId `bson:"_id"`
	EmployeeNo string        `bson:"employeeNo"`
	BrandCode  string        `bson:"brandCode"`
}
type UserBrandDB struct {
	EmployeeNo string `bson:"employeeNo"`
	BrandCode  string `bson:"brandCode"`
}

type UserPublisher struct {
	Id            bson.ObjectId `bson:"_id"`
	EmployeeNo    string        `bson:"employeeNo"`
	PublisherCode string        `bson:"publisherCode"`
}
type UserPublisherDB struct {
	EmployeeNo    string `bson:"employeeNo"`
	PublisherCode string `bson:"publisherCode"`
}

type UserLargeAndShop struct {
	Id           bson.ObjectId       `bson:"_id"`
	EmployeeNo   string              `bson:"employeeNo"`
	LargeAndShop []*ShowLargeAndShop `bson:"largeAndShop"`
}
type UserLargeAndShopDB struct {
	EmployeeNo   string              `bson:"employeeNo"`
	LargeAndShop []*ShowLargeAndShop `bson:"largeAndShop"`
}

type ChangePermissionAndRange struct {
	ChangePermission   []*UserPermissionDB
	ChangeCompany      []*UserCompanyDB
	ChangeBrand        []*UserBrandDB
	ChangeLargeAndShop []*ShowLargeAndShop
	ChangePublisher    []*UserPublisherDB
}

type CopyPermissionData struct {
	EmployeeNo   string   //来源账号
	Range        []int    //权限的来源范围
	EmployeeList []string //复制到哪些工号
}

type ShowPermission struct {
	Id       bson.ObjectId `bson:"_id"`
	Pid      string        `bson:"pid"`
	Name     string        `bson:"name"`
	Code     string        `bson:"code"`
	Orderno  int           `bson:"orderno"`
	Children []*ShowPermission
	Checked  bool
}

type ShowCompany struct {
	Code    string
	Name    string
	Checked bool
}

type ShowBrand struct {
	Code    string
	Name    string
	Checked bool
}

type ShowLargeAndShop struct {
	Code       string
	Name       string
	Checked    bool //界面对应勾选
	AllChecked bool //是否全选
	Selected   bool //是否拥有权限
	SalesBrand string
	Children   []*ShowLargeAndShop
}

type ShowPublisher struct {
	Code    string
	Name    string
	Checked bool
}

//我的权限范围,展示到界面
type MyPermissionRange struct {
	MyPermission   []*ShowPermission
	MyCompany      []*ShowCompany
	MyBrand        []*ShowBrand
	MyLargeAndShop []*ShowLargeAndShop
	MyPublisher    []*ShowPublisher
}

type LargeAndShop struct {
	LARGE      string
	CITYNAME   string
	WERKS      string
	SHOPNAME   string
	SalesBrand string
}
type StoreToBrand struct {
	shop_code  string
	brand_code string
}

//用户协议-同意表
type UserAgree struct {
	Id         bson.ObjectId `bson:"_id"`
	EmployeeNo string        `bson:"employeeno"`
	IsAgree    bool          `bson:"isagree"`
	CreateTime time.Time     `bson:"createtime"`
}
type UserAgreeDB struct {
	EmployeeNo string    `bson:"employeeno"`
	IsAgree    bool      `bson:"isagree"`
	CreateTime time.Time `bson:"createtime"`
}
