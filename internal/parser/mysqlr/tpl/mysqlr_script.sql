{{- define "mysqlr_script"}}{{- $obj := . -}}
{{- if ne $obj.DbTable ""}}

USE `{{$obj.DbName}}`;

CREATE TABLE `{{$obj.DbTable}}` (
	{{- range $i, $field := $obj.Fields}}
	{{$field.SQLColumn }},
	{{- end}}
	{{$obj.PrimaryKey.SQLColumn }},
	{{- range $i, $unique := $obj.Uniques}}
	{{- if not $unique.HasPrimaryKey}},
	{{- if and (eq (add $i 1) (len $obj.Uniques)) (eq (len $obj.Indexes) 0)}}
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
	{{- end}}
	{{- range $i, $index := $obj.Indexes}}
	{{- if not $index.HasPrimaryKey}}
	{{- if eq (add $i 1) (len $obj.Indexes) }}
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
	{{- end}}
	{{- end -}}
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '{{$obj.Comment}}';

{{- end}}

{{end}}
