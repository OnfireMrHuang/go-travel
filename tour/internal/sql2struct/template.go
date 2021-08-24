package sql2struct

import (
	"fmt"
	"go-travel/tour/internal/word"
	"html/template"
	"os"
)

const structTpl = `type {{ .TableName | ToCamelCase}} struct {
{{range .Columns}} {{ $length := len .Comment}} {{ if gt $length 0 }} //
{{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{ .Name| ToCamelCass}}
	{{.Type}} {{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | ToCamelCass}}) TableName() string {
	return "{{.TableName}}"
}`


type StructTemplate struct {
	structTpl string
}

type StructColumn struct {
	Name string
	Type string
	Tag string
	Comment string
}

type StructTemplateDB struct {
	TableName string
	Columns []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{structTpl: structTpl}
}

func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0 ,len(tbColumns))
	for _,column := range tbColumns {
		tag := fmt.Sprintf("`" + "json:" + "\"%s\"" + "`",column.ColumnName)
		tplColumns = append(tplColumns,&StructColumn{
			Name: column.ColumnName,
			Type: DBTypeToStructType[column.DataType],
			Tag: tag,
			Comment: column.ColumnComment,
		})
	}
	return tplColumns
}

func (t *StructTemplate) Generate(tableName string,tplColumns []*StructColumn) error {
	tpl := template.Must(template.New("sql2Struct").Funcs(template.FuncMap{
		"ToCamelCase": word.UnderscoreToLowerCamelCase,
	}).Parse(t.structTpl))
	tplDB := StructTemplateDB{
		TableName: tableName,
		Columns: tplColumns,
	}
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}
	return nil
}


