// Code generated by "gocqlx/cmd/schemagen"; DO NOT EDIT.

package scyllat

import (
	"github.com/scylladb/gocqlx/v2/table"
)

// Table models.
var (
	Records = table.New(table.Metadata{
		Name: "records",
		Columns: []string{
			"id",
			"record",
		},
		PartKey: []string{
			"id",
		},
		SortKey: []string{},
	})
)

type RecordsStruct struct {
	Id     [16]byte
	Record string
}
