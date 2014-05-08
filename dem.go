//Author:Centny
//Package DEM provide the testing sql driver.
package DEM

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
)

var G_Dn, G_Dsn string

var Evb = &EvBase{}

func init() {
	Register("DEM", Evb)
}

//register one drive to system by name.
func Register(n string, ev DbEv) {
	sql.Register(n, &STDriver{N: n, Ev: ev})
}

func OpenDem() *sql.DB {
	db, _ := sql.Open("DEM", G_Dsn)
	return db
}

//
var Error error = errors.New("text")

//the type of TDbErr
type STErr uint32

const (
	OPEN_ERR STErr = 1 << iota
	BEGIN_ERR
	CLOSE_ERR
	PREPARE_ERR
	TX_ROLLBACK_ERR
	TX_COMMIT_ERR
	STMT_CLOSE_ERR
	STMT_QUERY_ERR
	STMT_EXEC_ERR
	ROWS_CLOSE_ERR
	ROWS_NEXT_ERR
	EMPTY_DATA_ERR
	LAST_INSERT_ID_ERR
	ROWS_AFFECTED_ERR
)

func (t STErr) String() string {
	switch t {
	case OPEN_ERR:
		return "OPEN_ERR"
	case BEGIN_ERR:
		return "CONN_BEGIN_ERR"
	case CLOSE_ERR:
		return "CONN_CLOSE_ERR"
	case PREPARE_ERR:
		return "PREPARE_ERR"
	case TX_ROLLBACK_ERR:
		return "ROLLBACK_ERR"
	case TX_COMMIT_ERR:
		return "COMMIT_ERR"
	case STMT_CLOSE_ERR:
		return "STMT_CLOSE_ERR"
	case STMT_QUERY_ERR:
		return "STMT_QUERY_ERR"
	case STMT_EXEC_ERR:
		return "STMT_EXEC_ERR"
	case ROWS_CLOSE_ERR:
		return "ROWS_CLOSE_ERR"
	case ROWS_NEXT_ERR:
		return "ROWS_NEXT_ERR"
	case EMPTY_DATA_ERR:
		return "EMPTY_DATA_ERR"
	case LAST_INSERT_ID_ERR:
		return "LAST_INSERT_ID_ERR"
	case ROWS_AFFECTED_ERR:
		return "ROWS_AFFECTED_ERR"
	}
	return ""
}

//if error contain target error.
func (t STErr) Is(e STErr) bool {
	if (t & e) == e {
		return true
	} else {
		return false
	}
}

func (t STErr) IsErr(e STErr) error {
	if t.Is(e) {
		return errors.New(fmt.Sprintf("DEM %v", e.String()))
	} else {
		return nil
	}
}

type DbEv interface {
	//
	OnOpen(dsn string) (*sql.DB, error)
	//
	OnBegin(c *STConn) error
	OnPrepare(c *STConn, query string) error
	OnNumInput(c *STConn, query string, stm *sql.Stmt) int
	OnClose(c *STConn) error
	//
	OnTxCommit(tx *STTx) error
	OnTxRollback(tx *STTx) error
	//
	OnStmQuery(stm *STStmt, args []driver.Value) error
	OnStmExec(stm *STStmt, args []driver.Value) error
	OnStmClose(stm *STStmt) error
	//
	OnResLIId(res *STResult) error
	OnResARow(res *STResult) error
	//
	OnRowNext(row *STRows) error
	IsEmpty(row *STRows) bool
	OnRowClose(row *STRows) error
}
type Query struct {
	Q    *regexp.Regexp
	Args *regexp.Regexp
}

func (q *Query) Match(query string, args []driver.Value) bool {
	return q.Q.MatchString(query) && q.Args.MatchString(fmt.Sprintf("%v", args))
}

type EvBase struct {
	Errs STErr
	QErr []Query
	Dn   string
}

func NewEvBase(dn string) *EvBase {
	return &EvBase{
		Dn: dn,
	}
}
func (e *EvBase) SetErrs(err STErr) {
	e.Errs = err
}
func (e *EvBase) AddErrs(err STErr) *EvBase {
	e.Errs = e.Errs | err
	return e
}
func (e *EvBase) ClsQErr() {
	e.QErr = []Query{}
}
func (e *EvBase) AddQErr(err Query) *EvBase {
	e.QErr = append(e.QErr, err)
	return e
}
func (e *EvBase) AddQErr2(qreg string, areg string) *EvBase {
	return e.AddQErr(Query{
		Q:    regexp.MustCompile(qreg),
		Args: regexp.MustCompile(areg),
	})
}
func (e *EvBase) AddQErr3(qreg string) *EvBase {
	return e.AddQErr2(qreg, ".*")
}
func (e *EvBase) Match(query string, args []driver.Value) bool {
	for _, q := range e.QErr {
		if q.Match(query, args) {
			return true
		}
	}
	return false
}
func (e *EvBase) OnOpen(dsn string) (*sql.DB, error) {
	err := e.Errs.IsErr(OPEN_ERR)
	if err != nil {
		return nil, err
	}
	dn := ""
	if len(e.Dn) < 1 {
		dn = G_Dn
	}
	if len(dn) < 1 {
		return nil, errors.New("dbname is not initial for event handler")
	}
	return sql.Open(dn, dsn)
}
func (e *EvBase) OnBegin(c *STConn) error {
	return e.Errs.IsErr(BEGIN_ERR)
}
func (e *EvBase) OnPrepare(c *STConn, query string) error {
	return e.Errs.IsErr(PREPARE_ERR)
}
func (e *EvBase) OnNumInput(c *STConn, query string, stm *sql.Stmt) int {
	return strings.Count(query, "?")
}
func (e *EvBase) OnClose(c *STConn) error {
	return e.Errs.IsErr(CLOSE_ERR)
}
func (e *EvBase) OnTxCommit(tx *STTx) error {
	return e.Errs.IsErr(TX_COMMIT_ERR)
}
func (e *EvBase) OnTxRollback(tx *STTx) error {
	return e.Errs.IsErr(TX_ROLLBACK_ERR)
}
func (e *EvBase) OnStmQuery(stm *STStmt, args []driver.Value) error {
	err := e.Errs.IsErr(STMT_QUERY_ERR)
	if err != nil {
		return err
	}
	if e.Match(stm.Q, args) {
		return errors.New("DEM query matched error")
	}
	return nil
}
func (e *EvBase) OnStmExec(stm *STStmt, args []driver.Value) error {
	err := e.Errs.IsErr(STMT_EXEC_ERR)
	if err != nil {
		return err
	}
	if e.Match(stm.Q, args) {
		return errors.New("DEM query matched error")
	}
	return nil
}
func (e *EvBase) OnStmClose(stm *STStmt) error {
	return e.Errs.IsErr(STMT_CLOSE_ERR)
}
func (e *EvBase) OnResLIId(res *STResult) error {
	return e.Errs.IsErr(LAST_INSERT_ID_ERR)
}
func (e *EvBase) OnResARow(res *STResult) error {
	return e.Errs.IsErr(ROWS_AFFECTED_ERR)
}
func (e *EvBase) OnRowNext(row *STRows) error {
	return e.Errs.IsErr(ROWS_NEXT_ERR)
}
func (e *EvBase) IsEmpty(row *STRows) bool {
	return e.Errs.Is(EMPTY_DATA_ERR)
}
func (e *EvBase) OnRowClose(row *STRows) error {
	return e.Errs.IsErr(ROWS_CLOSE_ERR)
}

type STDriver struct {
	N  string //driver name.
	Ev DbEv   //database evernt
}

func (d *STDriver) Open(dsn string) (driver.Conn, error) {
	con, err := d.Ev.OnOpen(dsn)
	return &STConn{
		Db: con,
		Dr: d,
		Ev: d.Ev,
	}, err
}

type STConn struct {
	Db *sql.DB
	Dr *STDriver
	Ev DbEv //database evernt
}

func (c *STConn) Begin() (driver.Tx, error) {
	if e := c.Ev.OnBegin(c); e != nil {
		return nil, e
	}
	tx, err := c.Db.Begin()
	return &STTx{
		Tx:   tx,
		Conn: c,
		Ev:   c.Ev,
	}, err
}

func (c *STConn) Prepare(query string) (driver.Stmt, error) {
	if e := c.Ev.OnPrepare(c, query); e != nil {
		return nil, e
	}
	stm, err := c.Db.Prepare(query)
	return &STStmt{
		Q:    query,
		Conn: c,
		Stmt: stm,
		Ev:   c.Ev,
		Num:  c.Ev.OnNumInput(c, query, stm),
	}, err
}

func (c *STConn) Close() error {
	if e := c.Ev.OnClose(c); e != nil {
		return e
	}
	return c.Db.Close()
}

type STTx struct {
	Conn *STConn
	Tx   *sql.Tx
	Ev   DbEv //database evernt
}

func (tx *STTx) Commit() error {
	if e := tx.Ev.OnTxCommit(tx); e != nil {
		return e
	}
	return tx.Tx.Commit()
}

func (tx *STTx) Rollback() error {
	if e := tx.Ev.OnTxRollback(tx); e != nil {
		return e
	}
	return tx.Tx.Rollback()
}

type STStmt struct {
	Q    string
	Conn *STConn
	Stmt *sql.Stmt
	Ev   DbEv //database evernt
	Num  int
}

func (s *STStmt) NumInput() int {
	return s.Num
}

func (s *STStmt) Query(args []driver.Value) (driver.Rows, error) {
	if e := s.Ev.OnStmQuery(s, args); e != nil {
		return nil, e
	}
	targs := []interface{}{}
	for _, v := range args {
		targs = append(targs, v)
	}
	rows, e := s.Stmt.Query(targs...)
	return &STRows{
		Stmt: s,
		Rows: rows,
		Args: args,
		Ev:   s.Ev,
	}, e
}

func (s *STStmt) Exec(args []driver.Value) (driver.Result, error) {
	if e := s.Ev.OnStmExec(s, args); e != nil {
		return nil, e
	}
	targs := []interface{}{}
	for _, v := range args {
		targs = append(targs, v)
	}
	res, e := s.Stmt.Exec(targs...)
	return &STResult{
		Stmt: s,
		Res:  res,
		Args: args,
		Ev:   s.Ev,
	}, e
}

func (s *STStmt) Close() error {
	if e := s.Ev.OnStmClose(s); e != nil {
		return e
	}
	return s.Stmt.Close()
}

type STResult struct {
	Stmt *STStmt
	Res  sql.Result
	Args []driver.Value
	Ev   DbEv //database evernt
}

func (r *STResult) LastInsertId() (int64, error) {
	if e := r.Ev.OnResLIId(r); e != nil {
		return 0, e
	}
	return r.Res.LastInsertId()
}

func (r *STResult) RowsAffected() (int64, error) {
	if e := r.Ev.OnResARow(r); e != nil {
		return 0, e
	}
	return r.Res.RowsAffected()
}

type STRows struct {
	Stmt *STStmt
	Args []driver.Value
	Rows *sql.Rows
	Ev   DbEv //database evernt
}

func (rc *STRows) Columns() []string {
	cls, _ := rc.Rows.Columns()
	return cls
}

func (rc *STRows) Next(dest []driver.Value) error {
	if e := rc.Ev.OnRowNext(rc); e != nil {
		return e
	}
	if rc.Ev.IsEmpty(rc) {
		return io.EOF
	}
	if rc.Rows.Next() {
		l := len(rc.Columns())
		sary := make([]interface{}, l) //scan array.
		for i := 0; i < l; i++ {
			var a interface{}
			sary[i] = &a
		}
		e := rc.Rows.Scan(sary...)
		for i := 0; i < l; i++ {
			dest[i] = reflect.Indirect(reflect.ValueOf(sary[i])).Interface()
		}
		return e
	} else {
		return io.EOF
	}
}

func (rc *STRows) Close() error {
	if e := rc.Ev.OnRowClose(rc); e != nil {
		return e
	}
	return rc.Rows.Close()
}
