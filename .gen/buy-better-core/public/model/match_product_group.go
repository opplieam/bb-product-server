//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type MatchProductGroup struct {
	ID        int32 `sql:"primary_key"`
	GroupID   int32
	ProductID int32
}