package DEM

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/Centny/Cny4go/dbutil"
	"github.com/Centny/Cny4go/test"
	"github.com/Centny/Cny4go/util"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

const TDbCon string = "cny:123@tcp(127.0.0.1:3306)/cny?charset=utf8"

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

func TestDem(t *testing.T) {
	db, _ := sql.Open("DEM", test.TDbCon)
	db.Begin()
	NewEvBase("")
	//
	G_Dn = "mysql"
	G_Dsn = TDbCon
	ev := Evb
	db = OpenDem()
	T(db, t)
	db.Close()
	//
	db = OpenDem()
	dbutil.DbExecF(db, "ttable.sql")
	//
	ev.AddQErr3(".*ttable.*")
	T2(db, t)
	ev.ClsQErr()
	//
	ev.AddQErr3(".*t__able.*")
	//
	ev.SetErrs(BEGIN_ERR)
	T2(db, t)
	ev.SetErrs(PREPARE_ERR)
	T2(db, t)
	ev.SetErrs(TX_ROLLBACK_ERR)
	T2(db, t)
	ev.SetErrs(TX_COMMIT_ERR)
	T2(db, t)
	//
	ev.SetErrs(STMT_QUERY_ERR)
	T2(db, t)
	ev.SetErrs(STMT_EXEC_ERR)
	T2(db, t)
	ev.SetErrs(STMT_CLOSE_ERR)
	T2(db, t)
	//
	ev.SetErrs(ROWS_NEXT_ERR)
	T2(db, t)
	ev.SetErrs(ROWS_AFFECTED_ERR)
	T2(db, t)
	ev.SetErrs(ROWS_CLOSE_ERR)
	T2(db, t)
	ev.SetErrs(LAST_INSERT_ID_ERR)
	T2(db, t)
	//
	ev.SetErrs(EMPTY_DATA_ERR)
	T2(db, t)
	//
	ev.SetErrs(CLOSE_ERR)
	db.Close()
	//
	ev.SetErrs(0)
	db.Close()
	//
	ev.SetErrs(OPEN_ERR)
	db = OpenDem()
	db.Begin()
	//
	ev.AddErrs(OPEN_ERR)
	//
	var ee STErr = 100
	fmt.Println(ee.String())
	fmt.Println(OPEN_ERR.String())
	fmt.Println(EMPTY_DATA_ERR.String())
}
func T(db *sql.DB, t *testing.T) {
	err := dbutil.DbExecF(db, "ttable.sql")
	if err != nil {
		t.Error(err.Error())
	}
	res, err := dbutil.DbQuery(db, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(res) < 1 {
		t.Error("not data")
		return
	}
	if len(res[0]) < 1 {
		t.Error("data is empty")
		return
	}
	bys, err := json.Marshal(res)
	fmt.Println(string(bys))
	//
	var mres []TSt
	err = dbutil.DbQueryS(db, &mres, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(mres) < 1 {
		t.Error("not data")
		return
	}
	fmt.Println("...", mres[0].T, util.Timestamp(mres[0].Time), util.Timestamp(time.Now()))
	fmt.Println(mres, mres[0].Add1)
	//
	ivs, err := dbutil.DbQueryInt(db, "select * from ttable where tid")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(ivs) < 1 {
		t.Error("not data")
		return
	}
	//
	svs, err := dbutil.DbQueryString(db, "select tname from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(svs) < 1 {
		t.Error("not data")
		return
	}
	//
	iid, err := dbutil.DbInsert(db, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(iid)
	//
	tx, _ := db.Begin()
	iid2, err := dbutil.DbInsert2(tx, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(iid2)
	tx.Commit()
	//
	erow, err := dbutil.DbUpdate(db, "delete from ttable where tid=?", iid)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(erow, "-----")
	//
	tx, _ = db.Begin()
	erow, err = dbutil.DbUpdate2(tx, "delete from ttable where tid=?", iid2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(erow, "-----")
	tx.Commit()
	//
	_, err = dbutil.DbQuery(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = dbutil.DbQueryInt(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = dbutil.DbQueryString(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = dbutil.DbInsert(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx, _ = db.Begin()
	_, err = dbutil.DbInsert2(tx, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = dbutil.DbUpdate(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = dbutil.DbUpdate2(tx, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = dbutil.DbQuery(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = dbutil.DbQueryInt(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = dbutil.DbQueryString(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = dbutil.DbInsert(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = dbutil.DbInsert2(tx, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = dbutil.DbUpdate(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = dbutil.DbUpdate2(tx, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	err = dbutil.DbQueryS(nil, nil, "select * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	dbutil.DbQueryInt(nil, "select * from ttable where tid>?", 1, 2)
	dbutil.DbQueryString(nil, "select * from ttable where tid>?", 1, 2)
	dbutil.DbInsert(nil, "select * from ttable where tid>?", 1, 2)
	dbutil.DbUpdate(nil, "select * from ttable where tid>?", 1, 2)
	dbutil.DbInsert2(nil, "select * from ttable where tid>?", 1, 2)
	dbutil.DbUpdate2(nil, "select * from ttable where tid>?", 1, 2)
}
func T2(db *sql.DB, t *testing.T) {
	dbutil.DbQuery(db, "select * from ttable where tid>?", 1)
	iid, _ := dbutil.DbInsert(db, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	dbutil.DbUpdate(db, "delete from ttable where tid=?", iid)

	if tx, err := db.Begin(); err == nil {
		dbutil.DbInsert2(tx, "selectt * from ttable where tid>?", 1)
		tx.Rollback()
	}
	if tx, err := db.Begin(); err == nil {
		dbutil.DbInsert2(tx, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
		tx.Commit()
	}
}
func Map2Val2(columns []string, row map[string]interface{}, dest []driver.Value) {
	for i, c := range columns {
		if v, ok := row[c]; ok {
			switch c {
			case "INT":
				dest[i] = int(v.(float64))
			case "UINT":
				dest[i] = uint32(v.(float64))
			case "FLOAT":
				dest[i] = float32(v.(float64))
			case "SLICE":
				dest[i] = []byte(v.(string))
			case "STRING":
				dest[i] = v.(string)
			case "STRUCT":
				dest[i] = time.Now()
			case "BOOL":
				dest[i] = true
			}
		} else {
			dest[i] = nil
		}
	}
}

// func TestDbUtil2(t *testing.T) {
// 	TDb.Map2Val = Map2Val2
// 	db, _ := sql.Open("TDb", "td@tdata.json")
// 	defer db.Close()
// 	res, err := dbutil.DbQuery(db, "SELECT * FROM TESTING WHERE INT=? AND STRING=?", 1, "cny")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	fmt.Println(res)
// }

// func TestDbExecF(t *testing.T) {
// 	db, _ := sql.Open("mysql", test.TDbCon)
// 	defer db.Close()
// 	err := dbutil.DbExecF(db, "ttable.sql")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	dbutil.DbExecF(nil, "ttable.sql")
// 	dbutil.DbExecF(db, "ttables.sql")
// 	db.Close()
// 	dbutil.DbExecF(db, "ttable.sql")
// }
