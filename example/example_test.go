package example

import (
	"database/sql"
	"github.com/Centny/DEM"
	"github.com/Centny/gwf/dbutil"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestPanic(t *testing.T) {
	go func() {
		defer func() {
			recover()
		}()
		DefaultCon()
	}()
}

const TDbCon string = "cny:123@tcp(127.0.0.1:3306)/cny?charset=utf8"

func OpenDb() *sql.DB {
	db, _ := sql.Open("DEM", TDbCon)
	return db
}
func TestList(t *testing.T) {
	//setting database name
	DEM.G_Dn = "mysql"
	//setting database connection
	DEM.G_Dsn = TDbCon
	//default open database function.
	DbCon = DEM.OpenDem
	///
	dbutil.DbExecF(OpenDb(), "../ttable.sql")
	//
	err := ListData()
	if err != nil {
		t.Error(err.Error())
		return
	}
	//
	//mock open error
	DEM.Evb.SetErrs(DEM.OPEN_ERR | DEM.BEGIN_ERR)
	err = ListData()
	if err == nil {
		t.Error("not error")
	}
	//
	//mock  empty data error
	DEM.Evb.SetErrs(DEM.EMPTY_DATA_ERR)
	err = ListData()
	if err == nil {
		t.Error("not error")
	}
	//
	//clear error
	DEM.Evb.SetErrs(0)
	//add matched query error by regex
	DEM.Evb.AddQErr3(".*insert.*ttable.*")
	err = ListData()
	if err == nil {
		t.Error("not error")
	}
	//clear query error
	DEM.Evb.ClsQErr()
	//add matched query error by regex
	DEM.Evb.AddQErr3(".*delete.*ttable.*")
	err = ListData()
	if err == nil {
		t.Error("not error")
	}

	//error count match
	DEM.Evb.AddEC(DEM.OPEN_ERR, 2)
}
