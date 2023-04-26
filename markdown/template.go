package markdown

const TplPage = `# {title}
> 版本号：{version}<br>
> 作者：{author}<br>
> BaseUrl: {baseUrl}
`

const TplBody = `
# {name}

### 请求说明
> 请求方式：{method}<br>
请求URL ：{url}

### 请求参数
{reqTable}

### 请求示例
{reqParam}

### 返回参数
{respTable}

### 响应示例
{respParam}
`

const TplReqTable = `
| 字段      | 字段类型       | 必填     | 字段说明    |
|---------|--------------|--------|-----------|
{params}`

const TplReqParam = `| {name} | {type} | {required} | {description} |
`

const TplRespTable = `
| 字段      | 字段类型       | 字段说明    |
|---------|--------------|-----------|
{params}`

const TplRespParam = `| {name} | {type} | {description} |
`
