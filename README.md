Database Error Mocker
======
library for mock the sql error

it will intercept all database operation and check if mock error

===
#Install
```
go get github.com/Centny/DEM
```
#Example
###Test code

```
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
}
```

###Code:

```

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

````