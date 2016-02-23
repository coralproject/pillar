package main

import (
	"gopkg.in/mgo.v2"
)

type CollectionStats struct {
	Name            string
	Collection      *mgo.Collection
	ForeignKeyField string
}
