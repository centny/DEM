package example

import (
	"database/sql"
	"errors"
	"github.com/Centny/Cny4go/dbutil"
	"github.com/Centny/Cny4go/log"
	"time"
)

func DefaultCon() *sql.DB {
	panic("db not init")
}

var DbCon = DefaultCon

type TSt struct {
	Tid    int64     `m2s:"TID"`
	Tname  string    `m2s:"TNAME"`
	Titem  string    `m2s:"TITEM"`
	Tval   string    `m2s:"TVAL"`
	Status string    `m2s:"STATUS"`
	Time   time.Time `m2s:"TIME"`
	T      int64     `m2s:"TIME" it:"Y"`
	Fval   float64   `m2s:"FVAL"`
	Uival  int64     `m2s:"UIVAL"`
	Add1   string    `m2s:"ADD1"`
	Add2   string    `m2s:"Add2"`
}

func ListData() error {
	var ts []TSt
	err := dbutil.DbQueryS(DbCon(), &ts, "select * from ttable where tid>?", 1)
	if err != nil {
		log.D("Err:%v", err.Error())
		return err
	}
	if len(ts) < 1 {
		log.D("Err:%v", "NOT DATA")
		return errors.New("NOT DATA")
	}
	iid, err := dbutil.DbInsert(DbCon(), "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	if err != nil {
		log.D("Err:%v", err.Error())
		return err
	}
	_, err = dbutil.DbUpdate(DbCon(), "delete from ttable where tid=?", iid)
	if err != nil {
		log.D("Err:%v", err.Error())
		return err
	}
	return nil
}
