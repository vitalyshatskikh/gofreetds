package freetds

import (
	"errors"
	"fmt"
	"reflect"
)

type Result struct {
	Columns      []*ResultColumn
	Rows         [][]interface{}
	ReturnValue  int
	RowsAffected int
	Message      string
	currentRow   int
	scanCount     int
}

func NewResult() *Result {
	return &Result{
		Columns:    make([]*ResultColumn, 0),
		Rows:       nil,
		currentRow: -1,
	}
}

func (r *Result) addColumn(name string, dbSize, dbType int) {
	c := new(ResultColumn)
	c.Name = name
	c.DbSize = dbSize
	c.DbType = dbType
	r.Columns = append(r.Columns, c)
}

func (r *Result) addValue(row, col int, value interface{}) {
	if r.Rows == nil {
		r.Rows = make([][]interface{}, 1)
		r.Rows[0] = make([]interface{}, len(r.Columns))
	}
	for rc := len(r.Rows) - 1; rc < row; rc++ {
		r.Rows = append(r.Rows, make([]interface{}, len(r.Columns)))
	}
	r.Rows[row][col] = value
}

type ResultColumn struct {
	Name   string
	DbSize int
	DbType int
	Type   string
}

func (r *Result) Next() bool {
	if len(r.Rows) == 0 {
		return false
	}
	if r.currentRow >= len(r.Rows)-1 {
		return false
	}
	r.currentRow++
	return true
}

//Scan copies the columns in the current row into the values pointed at by dest.
func (r *Result) Scan(dest ...interface{}) error {
	r.scanCount = 0
	if r.currentRow == -1 {
		return errors.New("Scan called without calling Next.")
	}
	for _, d := range dest {
		if !isPointer(d) {
			return errors.New("Destination not a pointer.")
		}
	}
	if len(dest) == 1 {
		if s := asStructPointer(dest[0]); s != nil {
			return r.scanStruct(s)
		}
	}
	err := assignValues(r.Rows[r.currentRow], dest)
	if err == nil {
		r.scanCount = len(dest)
	}
	return err
}

//Must Scan exactly cnt number of values from result.
//Useful when scanning into structure, to know whether are all expected fields filled with values.
//cnt - number of values assigned to fields
func (r *Result) MustScan(cnt int, dest ...interface{}) error {
	if err := r.Scan(dest...); err != nil {
		return err
	}
	if cnt != r.scanCount {
		return errors.New(fmt.Sprintf("Worng scan count, expected %d, actual %d.", cnt, r.scanCount))
	}
	return nil
}

//Copies values for the current row to the structure.
//Struct filed name must match database column name.
func (r *Result) scanStruct(s *reflect.Value) error {
	for i, col := range r.Columns {
		f := s.FieldByName(camelize(col.Name))
		if f.IsValid() {
			if f.CanSet() {
				if err := convertAssign(f.Addr().Interface(), r.Rows[r.currentRow][i]); err != nil {
				//if err := assignValue(r.Rows[r.currentRow][i], f.Addr().Interface());  err != nil {
					return err
				}
				r.scanCount++
			}
		}
	}
	return nil
}

func asStructPointer(p interface{}) *reflect.Value {
	sp := reflect.ValueOf(p)
	if sp.Kind() == reflect.Ptr {
		s := sp.Elem() 
		if s.Kind() == reflect.Struct {
			return &s
		}
	}
	return nil
}

func isPointer(p interface{}) bool {
	sp := reflect.ValueOf(p)
	return sp.Kind() == reflect.Ptr 
}

//assignValues copies to dest values in src
//dest should be a pointer type
//error is returned if types don't match and conversion failed
func assignValues(src, dest []interface{}) error {
	if len(dest) > len(src) {
		return errors.New(fmt.Sprintf("More dest values %d than src values %d.", len(dest), len(src)))
	}
	for i, d := range dest {
		err := convertAssign(d, src[i])
		if err != nil {
			return err
		}
	}
	return nil
}
