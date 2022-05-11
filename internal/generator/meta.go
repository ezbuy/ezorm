package generator

type IField interface {
	GetName() string
	GetGoType() string
}

type IObject interface {
	Read(string, Schema) error
	FieldsMap() map[string]IField
	GetTable() string
}
