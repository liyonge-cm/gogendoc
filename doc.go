package gogendoc

import (
	"fmt"
	"reflect"
	"strings"
)

type MethodType string
type UrlType string
type RequestType interface{}
type ResponseType interface{}

const (
	POST   MethodType = "POST"
	GET    MethodType = "GET"
	PUT    MethodType = "PUT"
	DELETE MethodType = "DELETE"
)

const (
	Author string = "Yonge"
)

type DocGroup struct {
	Name string
	List []DocItem
}

type Document struct {
	Title          string               `json:"title"`   // 文档标题
	Version        string               `json:"version"` // 版本号
	BaseUrl        string               `json:"baseUrl"` // BaseUrl
	Author         string               `json:"author"`  // 作者
	Group          map[string]*DocGroup `json:"group"`   // 分组
	List           []DocItem            `json:"list"`    // 文档列表
	getName        func(field reflect.StructField) string
	getRequired    func(field reflect.StructField) bool
	getDescription func(field reflect.StructField) string
	FileNameKey    string
}

// 字段
type Field struct {
	Name        string  `json:"name"`        // 字段名称
	Type        string  `json:"type"`        // 字段类型
	Kind        string  `json:"kind"`        // 字段类型
	Required    bool    `json:"required"`    // 是否必填
	Description string  `json:"description"` // 字段说明
	List        []Field `json:"list"`        // 字段列表
}

// 文档对象
type DocItem struct {
	Title      string       `json:"title"`     // 标题
	Url        UrlType      `json:"url"`       // 接口地址
	Method     MethodType   `json:"method"`    // 请求类型
	ReqParam   RequestType  `json:"reqParam"`  // 请求参数
	RespParam  ResponseType `json:"respParam"` // 返回参数
	ReqFields  []Field      `json:"fields"`    // 字段列表
	RespFields []Field      `json:"fields"`    // 字段列表
}

func NewDocument(doc *Document) *Document {
	doc.getName = getName
	doc.getRequired = getRequired
	doc.getDescription = getDescription
	doc.Group = make(map[string]*DocGroup)
	return doc
}

/*
添加文档对象
title 接口名称
url 相对路径
method http请求方法
author 作者
req 请求结构体（需要实例化）
resp 返回结果结构体（需要实例化）
*/
func (d *Document) AddItem(title string, url UrlType, method MethodType, req, resp interface{}) {
	v := DocItem{
		Title:     title,
		Url:       url,
		Method:    method,
		ReqParam:  req,
		RespParam: resp,
	}
	if len(d.List) == 0 {
		d.List = make([]DocItem, 0)
		d.List = append(d.List, v)
	} else {
		d.List = append(d.List, v)
	}
}

func (d *Document) NewGroup(name string) *DocGroup {
	d.Group[name] = &DocGroup{
		Name: name,
		List: []DocItem{},
	}
	return d.Group[name]
}
func (g *DocGroup) AddGroupItem(title string, url UrlType, method MethodType, req, resp interface{}) {
	v := DocItem{
		Title:     title,
		Url:       url,
		Method:    method,
		ReqParam:  req,
		RespParam: resp,
	}

	if len(g.List) == 0 {
		g.List = make([]DocItem, 0)
		g.List = append(g.List, v)
	} else {
		g.List = append(g.List, v)
	}
	//d.Group[g.name] = g
}

// 生成接口列表
func (d *Document) GenerateFields() {
	if len(d.List) > 0 {
		for i, docItem := range d.List {
			docItem.ReqFields = d.createFields(docItem.ReqParam)
			docItem.RespFields = d.createFields(docItem.RespParam)
			d.List[i] = docItem
		}
	}
	if len(d.Group) > 0 {
		for _, g := range d.Group {
			for i, docItem := range g.List {
				docItem.ReqFields = d.createFields(docItem.ReqParam)
				docItem.RespFields = d.createFields(docItem.RespParam)
				g.List[i] = docItem
			}
		}
	}
}

// 获取接口列表
func (d *Document) GetList() []DocItem {
	return d.List
}
func (d *Document) GetGroup() map[string]*DocGroup {
	return d.Group
}

func (d *Document) GetFieldName(getNameFunc func(field reflect.StructField) string) {
	d.getName = getNameFunc
}
func (d *Document) GetFieldRequired(getRequiredFunc func(field reflect.StructField) bool) {
	d.getRequired = getRequiredFunc
}
func (d *Document) GetFieldDescription(getDescriptionFunc func(field reflect.StructField) string) {
	d.getDescription = getDescriptionFunc
}

// 创建字段
func (d *Document) createFields(param interface{}) []Field {
	if param == nil {
		return nil
	}
	fields := make([]Field, 0)
	val := reflect.ValueOf(param)
	if !val.IsValid() {
		panic("not valid")
	}
	// fmt.Println(val.Kind())
	if val.Kind() == reflect.Slice {
		if val.Len() > 0 {
			val = val.Index(0)
		} else {
			return nil
		}
	}
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	typ := val.Type()
	cnt := val.NumField()
	for i := 0; i < cnt; i++ {
		fd := val.Field(i)
		kd := fd.Kind()
		ty := typ.Field(i)
		// fmt.Println("ty Name", ty.Name)
		// fmt.Println("ty PkgPath", ty.PkgPath)
		// fmt.Println("ty Type", ty.Type)
		// fmt.Println("ty Tag", ty.Tag)
		// fmt.Println("ty Offset", ty.Offset)
		// fmt.Println("ty Index", ty.Index)
		// fmt.Println("ty Anonymous", ty.Anonymous)
		// fmt.Println("kd Kind", kd.String())
		// fmt.Println("fd", fd)

		if d.getName(ty) == "-" {
			continue
		}
		field := Field{
			Name:        d.getName(ty),
			Type:        fmt.Sprint(ty.Type),
			Kind:        kd.String(),
			Required:    d.getRequired(ty),
			Description: d.getDescription(ty),
			List:        nil,
		}
		if field.Kind == "interface" || field.Kind == "struct" {
			subFields := d.createFields(fd.Interface())
			field.List = subFields
		}
		// if field.Kind == "slice" {
		// 	switch fd.Index(0).Interface().(type) {
		// 	case string:
		// 		field.Kind = "字符串数组"
		// 	case int:
		// 		field.Kind = "整型数组"
		// 	case bool:
		// 		field.Kind = "布尔数组"
		// 	case float64:
		// 		field.Kind = "浮点型数组"
		// 	case float32:
		// 		field.Kind = "浮点型数组"
		// 	default: //结构体
		// 		field.Kind = "结构体数组"
		// 		subFields := createFields(fd.Interface())
		// 		field.List = subFields
		// 	}
		// }
		// 如果是数字型字符串 例 Id int `json:"id,string"`
		// if field.Kind == "int" && strings.Contains(ty.Tag.Get("json"), ",string") {
		// 	field.Kind = "string"
		// }
		//如果是内嵌结构体
		if ty.Anonymous {
			subFields := d.createFields(fd.Interface())
			fields = append(fields, subFields...)
		} else {
			fields = append(fields, field)
		}
	}
	return fields
}

func getName(field reflect.StructField) string {
	name := field.Tag.Get("json")
	if name == "" {
		name = field.Name
	}
	return name
}

// 获取tag
func getDescription(field reflect.StructField) string {
	desc := field.Tag.Get("comment")
	return desc
}

// 判断是否必填
func getRequired(field reflect.StructField) bool {
	tag := field.Tag.Get("v")
	return strings.Contains(tag, "required")
}
