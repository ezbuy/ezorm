{{- define "mysql_script"}}{{$objs := .}}
{{- range $obj := $objs}}
-- DDL for object {{$obj.Name}}.
CREATE TABLE `{{$obj.Table}}` (
{{- range $field := $obj.Fields}}
  {{$field.MysqlCreation}},
{{- end}}
  PRIMARY KEY (`{{camel2name $obj.GetPrimaryKeyName}}`)
  {{- if gt (len $obj.Indexes) 0 -}}
  ,
  {{- end}}
{{- range $i, $index := $obj.Indexes}}
  {{- if eq (add $i 1) (len $obj.Indexes)}}
  {{$index.MysqlCreation $obj}}
  {{- else}}
  {{$index.MysqlCreation $obj}},
  {{- end}}
{{- end}}) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '{{$obj.Comment}}';

{{end}}

{{- end}}
