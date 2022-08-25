package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (repo *SQLiteRepository) Migrate() error {

	query := `
		CREATE TABLE IF NOT EXISTS ride(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rideid TEXT NOT NULL UNIQUE,
			driverid TEXT UNIQUE,
			clientid TEXT NOT NULL UNIQUE,
			status TEXT
		);
	`
	_, err := repo.db.Exec(query)
	return err
}

func (repo *SQLiteRepository) Create(ride Ride) (*Ride, error) {
	res, err := repo.db.Exec("INSERT INTO ride(rideid, driverid, clientid, status) values(?,?,?,?)",
		ride.RideId, ride.DriverId, ride.Clientid, ride.Status)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	fmt.Println(id)
	return &ride, nil
}

func (repo *SQLiteRepository) All() ([]Ride, error) {
	rows, err := repo.db.Query("SELECT * FROM ride WHERE status = ?", "initiated")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []Ride

	for rows.Next() {
		var ride Ride
		if err := rows.Scan(&ride.Id, &ride.RideId, &ride.Clientid, &ride.DriverId,
			&ride.Status); err != nil {
			return nil, err
		}
		all = append(all, ride)
	}
	return all, nil
}

func (repo *SQLiteRepository) GetByName(name string) (*Ride, error) {
	row := repo.db.QueryRow("SELECT * FROM contact WHERE name = ?", name)

	var ride Ride
	if err := row.Scan(&ride.RideId, &ride.Clientid, &ride.DriverId,
		&ride.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}

		return nil, err
	}
	return &ride, nil
}

func (repo *SQLiteRepository) Update(id int64, newride Ride) (*Ride, error) {

	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	res, err := repo.db.Exec("UPDATE contact SET driverid = ?, status = ? WHERE id = ? ",
		newride.DriverId, newride.Status)

	if err != nil {
		return nil, err
	}

	rowsaffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsaffected == 0 {
		return nil, ErrUpdateFailed
	}
	return &newride, nil
}

func (repo *SQLiteRepository) Delete(rideid int64) error {
	res, err := repo.db.Exec("DELETE FROM websites WHERE rideid = ?", rideid)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
