# gogendoc

文档生成工具

## MarkDown生成说明

### 示例

```go

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/liyonge-cm/gogendoc"
	"github.com/liyonge-cm/gogendoc/markdown"

)

type ChannelUser struct {
	UserID   int    `json:"UserID" v:"required" comment:"渠道用户ID"`
	UserName string `json:"UserName" v:"required" comment:"渠道用户名称"`
	Channel  string `json:"Channel" v:"required" comment:"渠道"`
}
type CreateUserInfo struct {
	UserAge     int         `json:"UserAge" v:"required" comment:"用户年龄"`
	UserName    string      `json:"UserName" v:"required" comment:"用户名称"`
	Type        string      `json:"Type" v:"required" comment:"用户类型"`
	ChannelUser ChannelUser `json:"ChannelUser" comment:"来源渠道用户"`
	ValidTime   int         `json:"ValidTime" v:"required" comment:"有效期"`
}

type CreateUserInfoResponse struct {
	Code    int         `json:"RetCode"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func TestGenDoc(t *testing.T) {
	// 实例化文档
	doc := gogendoc.NewDocument(&gogendoc.Document{
		Title:   "用户接口文档",
		Author:  gogendoc.Author,
		BaseUrl: "http://xxx",
	})
	doc.AddItem("创建用户信息", "CreateUserInfo", gogendoc.POST, &CreateUserInfo{}, &CreateUserInfoResponse{Data: &CreateUserInfo{}})
	doc.GenerateFields()
	list := doc.GetList()
	b, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestGenMarkDown(t *testing.T) {
	// 实例化文档
	doc := gogendoc.NewDocument(&gogendoc.Document{
		Title:       "用户接口文档1111",
		Version:     "1.0.0",
		Author:      gogendoc.Author,
		BaseUrl:     "http://xxx",
		FileNameKey: "title", //  title、url
	})
	// 自定义字段名称
	// doc.GetFieldName(func(field reflect.StructField) string { return field.Tag.Get("json") })

	// 自定义字段必填
	// doc.GetFieldRequired(func(field reflect.StructField) bool { return true })

	// 自定义字段说明
	// doc.GetFieldDescription(func(field reflect.StructField) string { return field.Tag.Get("desc") })

	// 添加接口
	doc.AddItem("创建用户信息", "/createUserInfo", gogendoc.POST, &CreateUserInfo{}, &CreateUserInfoResponse{Data: &CreateUserInfo{}})

	// 生成字段
	doc.GenerateFields()
	// 实例化文档生成器
	md := markdown.New(doc)
	// 生成文档
	md.Generate("./docs")
}
```