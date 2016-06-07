package behavior

/*
An Item is a type that has an ID and a type

The traditional usage of the ID/Type is to map items into
 stores such as dbs and caches, define relationships and
 kindle the imagination
*/

type Item interface {
	GetId() interface{}
	GetType() string
}
