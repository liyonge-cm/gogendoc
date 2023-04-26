package markdown

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/liyonge-cm/gogendoc"
)

type SubTable struct {
	Title  string
	Fields []gogendoc.Field
}

type Markdown struct {
	doc         *gogendoc.Document
	id          int
	subReqList  []SubTable
	subRespList []SubTable
}

func New(doc *gogendoc.Document) *Markdown {
	return &Markdown{doc: doc, subReqList: make([]SubTable, 0), subRespList: make([]SubTable, 0)}
}

func (m *Markdown) Generate(file string) {
	err := deleteDir(file)
	if err != nil {
		panic(err)
	}
	err = createDir(file)
	if err != nil {
		panic(err)
	}

	readmeFile := path.Join(file, "README.md")
	readmePage := m.RenderReadmePage(m.doc)
	err = createFile(readmeFile, []byte(readmePage))
	if err != nil {
		panic(err)
	}

	m.generateList(file, m.doc.GetList())

	group := m.doc.GetGroup()
	for _, g := range group {
		groupFile := path.Join(file, g.Name)
		err = createDir(groupFile)
		if err != nil {
			panic(err)
		}
		m.generateList(groupFile, g.List)
	}
}

func (m *Markdown) generateList(baseFile string, list []gogendoc.DocItem) {
	for i, item := range list {
		pageFile := path.Join(baseFile, string(item.Url)+".md")
		if m.doc.FileNameKey == "title" {
			pageFile = path.Join(baseFile, item.Title+".md")
		}
		page := m.RenderPage(i, item)
		_ = os.Remove(pageFile)
		err := createFile(pageFile, []byte(page))
		if err != nil {
			panic(err)
		}
	}
}

func (m *Markdown) RenderReadmePage(v *gogendoc.Document) string {
	ts := TplPage
	ts = strings.Replace(ts, "{title}", v.Title, 1)
	ts = strings.Replace(ts, "{version}", v.Version, 1)
	ts = strings.Replace(ts, "{author}", v.Author, 1)
	ts = strings.Replace(ts, "{baseUrl}", v.BaseUrl, 1)
	return ts
}

func (m *Markdown) RenderPage(id int, v gogendoc.DocItem) string {
	m.subReqList = make([]SubTable, 0)
	m.subRespList = make([]SubTable, 0)
	m.id = id
	ts := TplBody
	ts = strings.Replace(ts, "{id}", fmt.Sprintf("%v", id), 1)
	ts = strings.Replace(ts, "{name}", v.Title, 1)
	ts = strings.Replace(ts, "{method}", string(v.Method), 1)
	ts = strings.Replace(ts, "{url}", string(v.Url), 1)
	if len(v.ReqFields) > 0 {
		reqTable := m.RenderReqTable("", v.ReqFields)
		if len(m.subReqList) > 0 {
			subTable := ""
			for _, item := range m.subReqList {
				tpl := m.RenderReqTable(item.Title, item.Fields)
				subTable = fmt.Sprintf("%s%s", subTable, tpl)
			}
			reqTable = fmt.Sprintf("%s%s", reqTable, subTable)
		}
		ts = strings.Replace(ts, "{reqTable}", reqTable, 1)
	} else {
		ts = strings.Replace(ts, "{reqTable}", "", 1)
	}
	if len(v.RespFields) > 0 {
		respTable := m.RenderRespTable("", v.RespFields)
		if len(m.subRespList) > 0 {
			subTable := ""
			for _, item := range m.subRespList {
				tpl := m.RenderRespTable(item.Title, item.Fields)
				subTable = fmt.Sprintf("%s%s", subTable, tpl)
			}
			respTable = fmt.Sprintf("%s%s", respTable, subTable)
		}
		ts = strings.Replace(ts, "{respTable}", respTable, 1)
	} else {
		ts = strings.Replace(ts, "{respTable}", "", 1)
	}
	if v.ReqParam != nil {
		reqParam, _ := json.MarshalIndent(v.ReqParam, "", "\t")
		ts = strings.Replace(ts, "{reqParam}", fmt.Sprintf("```json\n %s \n```", string(reqParam)), 1)
	} else {
		ts = strings.Replace(ts, "{reqParam}", "", 1)
	}
	if v.RespParam != nil {
		respParam, _ := json.MarshalIndent(v.RespParam, "", "\t")
		ts = strings.Replace(ts, "{respParam}", fmt.Sprintf("```json\n %s \n```", string(respParam)), 1)
	} else {
		ts = strings.Replace(ts, "{respParam}", "", 1)
	}
	return ts

}

func (m *Markdown) RenderReqTable(title string, fields []gogendoc.Field) string {
	ts := ""
	if title != "" {
		ts = fmt.Sprintf("\n<a id=\"%d.%s\"></a> \n##### %s \n %s ", m.id, title, title, TplReqTable)
	} else {
		ts = TplReqTable
	}

	params := ""
	for _, v := range fields {
		tpl := m.RenderReqParam(v)
		params = fmt.Sprintf("%s%s", params, tpl)
	}
	ts = strings.Replace(ts, "{params}", params, 1)
	return ts
}

func (m *Markdown) RenderReqParam(v gogendoc.Field) string {
	ts := TplReqParam
	required := "是"
	if !v.Required {
		required = "否"
	}
	ts = strings.Replace(ts, "{name}", v.Name, 1)
	ts = strings.Replace(ts, "{type}", v.Type, 1)
	ts = strings.Replace(ts, "{required}", required, 1)
	if len(v.List) > 0 {
		subTable := SubTable{
			Title:  v.Name,
			Fields: v.List,
		}
		m.subReqList = append(m.subReqList, subTable)
		v.Description = fmt.Sprintf("%s [go](#%d.%s)", v.Description, m.id, v.Name)
	}
	ts = strings.Replace(ts, "{description}", v.Description, 1)
	return ts
}

func (m *Markdown) RenderRespTable(title string, fields []gogendoc.Field) string {
	ts := ""
	if title != "" {
		ts = fmt.Sprintf("\n<a id=\"%d.%s\"></a> \n##### %s \n %s ", m.id, title, title, TplRespTable)
	} else {
		ts = TplRespTable
	}
	params := ""
	for _, v := range fields {
		tpl := m.RenderRespParam(v)
		params = fmt.Sprintf("%s%s", params, tpl)
	}
	ts = strings.Replace(ts, "{params}", params, 1)
	return ts
}

func (m *Markdown) RenderRespParam(v gogendoc.Field) string {
	ts := TplRespParam
	ts = strings.Replace(ts, "{name}", v.Name, 1)
	ts = strings.Replace(ts, "{type}", v.Type, 1)
	if len(v.List) > 0 {
		subTable := SubTable{
			Title:  v.Name,
			Fields: v.List,
		}
		m.subRespList = append(m.subRespList, subTable)
		v.Description = fmt.Sprintf("%s [go](#%d.%s)", v.Description, m.id, v.Name)
	}
	ts = strings.Replace(ts, "{description}", v.Description, 1)
	return ts
}

func createFile(fileName string, val []byte) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(val)
	if err != nil {
		return err
	}
	return err
}

func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func createDir(path string) error {
	if !fileExists(path) {
		err := os.Mkdir(path, os.ModePerm)
		return err
	}
	return nil
}

func deleteDir(path string) error {
	return os.RemoveAll(path)
}
