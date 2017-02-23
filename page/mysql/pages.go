package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mono83/tame/page"
	"net/http"
	"net/url"
	"time"
)

// Repository represents DAO, that binds to sql connection and sourceId
type Repository interface {
	Store(page.Page) (uint64, error)
	Get(uint64) (*page.Page, error)
	Find(url string) (*page.Page, error)
}

type repository struct {
	DB       *sql.DB
	Table    string
	SourceID *int

	stmtInsert, stmtGet, stmtFind, stmtFindWithSource string
}

// NewRepository builds and return DAO, bound to SQL connection and sourceId, that
// can be used for storing and retrieving pages.
func NewRepository(db *sql.DB, table string, sourceID *int) Repository {
	r := &repository{
		DB:       db,
		Table:    table,
		SourceID: sourceID,
	}

	// Preparing statements
	r.stmtInsert = "INSERT INTO `" + table + "` SET" +
		" `sourceID` = ?, `createdAt` = ?, `headers` = ?, `body` = ?," +
		" `timeDns` = ?, `timeConnection` = ?, `timeRequestSent` = ?, `timeTotal` = ?," +
		" `url` = ?, `hash` = CRC32(?)"
	sel := "SELECT `sourceID`, `createdAt`, `url`, `headers`, `body`," +
		" `timeDns`, `timeConnection`, `timeRequestSent`, `timeTotal`" +
		" FROM `" + table + "`"
	r.stmtGet = sel + " WHERE `id` = ? LIMIT 1"
	r.stmtFind = sel + " WHERE `hash` = CRC32(?) AND `url` = ? ORDER BY `createdAt` DESC LIMIT 1"
	r.stmtFindWithSource = sel + " WHERE `hash` = CRC32(?) AND `url` = ? AND `sourceID` = ?" +
		" ORDER BY `createdAt` DESC LIMIT 1"

	return r
}

func (r *repository) Store(p page.Page) (uint64, error) {
	if p.StatusCode != 200 {
		return 0, fmt.Errorf("Allowed to store only HTTP 200 OK, but got %d", p.StatusCode)
	}
	stmt, err := r.DB.Prepare(r.stmtInsert)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	createdAt := p.Trace.Init
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	headerBts, err := json.Marshal(p.Header)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(
		r.SourceID,
		createdAt.Unix(),
		headerBts,
		p.Body,
		packDuration(p.Trace.DNS),
		packDuration(p.Trace.Connection),
		packDuration(p.Trace.RequestSent),
		packDuration(p.Trace.Total),
		p.URL.String(),
		p.URL.String(),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *repository) Get(id uint64) (*page.Page, error) {
	if id < 1 {
		return nil, errors.New("Invalid id")
	}

	stmt, err := r.DB.Prepare(r.stmtGet)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// Reading result set
	pages, err := readPages(res)
	if err != nil {
		return nil, err
	}

	l := len(pages)
	if l == 0 {
		return nil, nil
	} else if l > 1 {
		return nil, fmt.Errorf("Expected 1 result, but got %d", l)
	}

	p := pages[0]

	return &p, nil
}

func (r *repository) Find(url string) (*page.Page, error) {
	if len(url) == 0 {
		return nil, errors.New("Invalid url")
	}
	var stmt *sql.Stmt
	var res *sql.Rows
	var err error
	if r.SourceID == nil {
		stmt, err = r.DB.Prepare(r.stmtFind)
	} else {
		stmt, err = r.DB.Prepare(r.stmtFindWithSource)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if r.SourceID == nil {
		res, err = stmt.Query(url, url)
	} else {
		res, err = stmt.Query(url, url, r.SourceID)
	}

	if err != nil {
		return nil, err
	}
	defer res.Close()

	// Reading result set
	pages, err := readPages(res)
	if err != nil {
		return nil, err
	}

	l := len(pages)
	if l == 0 {
		return nil, nil
	} else if l > 1 {
		return nil, fmt.Errorf("Expected 1 result, but got %d", l)
	}

	p := pages[0]

	return &p, nil
}

// readPages is utility method, that builds Page objects from SQL result set
func readPages(rows *sql.Rows) ([]page.Page, error) {
	var sourceID *int
	var createdAt int64
	var urlStr string
	var headers, body []byte
	var timeDNS, timeConnection, timeRequestSent, timeTotal uint16

	response := []page.Page{}
	for rows.Next() {
		// SELECT `sourceID`, `createdAt`, `url`, `headers`, `body`,
		// `timeDNS`, `timeConnection`, `timeRequestSent`, `timeTotal`
		err := rows.Scan(
			&sourceID,
			&createdAt,
			&urlStr,
			&headers,
			&body,
			&timeDNS,
			&timeConnection,
			&timeRequestSent,
			&timeTotal,
		)
		if err != nil {
			return nil, err
		}

		// Building page
		p := page.Page{
			Body:       body,
			StatusCode: 200,
			Trace: page.Trace{
				Init:        time.Unix(createdAt, 0),
				DNS:         unpackDuration(timeDNS),
				Connection:  unpackDuration(timeConnection),
				RequestSent: unpackDuration(timeRequestSent),
				Total:       unpackDuration(timeTotal),
			},
		}

		// Converting URL
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		p.URL = u

		// Converting headers
		var h map[string][]string
		err = json.Unmarshal(headers, &h)
		if err != nil {
			return nil, err
		}
		p.Header = http.Header(h)

		response = append(response, p)
	}

	return response, nil
}

// packDuration converts time.Duration into milliseconds and packs it into uint16
func packDuration(duration time.Duration) uint16 {
	millis := duration.Nanoseconds() / 1e6
	if millis > 65535 {
		millis = 65535
	}

	return uint16(millis)
}

// unpackDuration converts milliseconds into time.Duration
func unpackDuration(millis uint16) time.Duration {
	return time.Duration(int64(millis) * 1e6)
}
