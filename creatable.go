package jobj

type CreatableSchema interface {
	CreateDescription() CreatableSchema
	CreateFields() CreatableSchema
	GetDescription() string
	GetFields() []*Field
}
