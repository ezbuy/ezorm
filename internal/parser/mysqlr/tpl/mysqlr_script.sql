{{- define "mysqlr_script"}}{{- $obj := . -}}
{{- if ne $obj.DbTable ""}}

USE `{{$obj.DbName}}`;

CREATE TABLE `{{$obj.DbTable}}` (
	{{- range $i, $field := $obj.Fields}}
	{{$field.SQLColumn }},
	{{- end}}
	{{- if and (eq (len $obj.NotPrimaryUniques) 0) (eq (len $obj.NotPrimaryIndexes) 0)}}
	{{$obj.PrimaryKey.SQLColumn }}
	{{- else}}
	{{$obj.PrimaryKey.SQLColumn }},
	{{- end}}
	{{- range $i, $unique := $obj.NotPrimaryUniques}}
	{{- if and (eq (add $i 1) (len $obj.NotPrimaryUniques)) (eq (len $obj.NotPrimaryIndexes) 0)}}
	UNIQUE KEY `uniq_{{$unique.Name | camel2name}}` (
		{{- range $i, $f := $unique.Fields -}}
			{{- if eq (add $i 1) (len $unique.Fields) -}}
				`{{- $f.Name | camel2name -}}`
			{{- else -}}
				`{{- $f.Name | camel2name -}}`,
			{{- end -}}
		{{- end -}}
	)
	{{- else}}
	UNIQUE KEY `uniq_{{$unique.Name | camel2name}}` (
		{{- range $i, $f := $unique.Fields -}}
			{{- if eq (add $i 1) (len $unique.Fields) -}}
				`{{- $f.Name | camel2name -}}`
			{{- else -}}
				`{{- $f.Name | camel2name -}}`,
			{{- end -}}
		{{- end -}}
	),
	{{- end}}
	{{- end}}
	{{- range $i, $index := $obj.NotPrimaryIndexes}}
	{{- if eq (add $i 1) (len $obj.NotPrimaryIndexes) }}
	KEY `{{$index.Name | camel2name}}` (`
		{{- range $i, $f := $index.Fields -}}
			{{- if eq (add $i 1) (len $index.Fields) -}}
				{{- $f.Name | camel2name -}}
			{{- else -}}
				{{- $f.Name | camel2name -}}
			{{- end -}}
		{{- end -}}
	`)
	{{- else}}
	KEY `{{$index.Name | camel2name}}` (`
		{{- range $i, $f := $index.Fields -}}
			{{- if eq (add $i 1) (len $index.Fields) -}}
				{{- $f.Name | camel2name -}}
			{{- else -}}
				{{- $f.Name | camel2name -}}
			{{- end -}}
		{{- end -}}
	`),
	{{- end}}
	{{- end -}}
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '{{$obj.Comment}}';

{{- end}}

{{end}}
