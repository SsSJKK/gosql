package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

// ErrNotFound ...
var ErrNotFound = errors.New("item not found")

// ErrInternal ...
var ErrInternal = errors.New("internal error")

// Service ...
type Service struct {
	db *sql.DB
}

// NewService ...
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

//Customer ...
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

// ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

//All ...
func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, phone, active, created FROM customers`)

	if errors.Is(err, sql.ErrNoRows) {
		log.Print("No Rows")
		return nil, nil
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

//AllActive ...
func (s *Service) AllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, phone, active, created FROM customers WHERE active = TRUE`)

	if errors.Is(err, sql.ErrNoRows) {
		log.Print("No Rows")
		return nil, nil
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

//Save ...
func (s *Service) Save(ctx context.Context, id int64, name string, phone string) (*Customer, error) {
	item := &Customer{}
	if id == 0 {
		err := s.db.QueryRowContext(ctx, `
		INSERT INTO customers(name, phone)
		VALUES($1, $2) ON CONFLICT (phone) DO
		NOTHING
		RETURNING id,
    	name,
    	phone,
    	active,
		created;`, name, phone).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		return item, nil

	}

	_, err := s.ByID(ctx, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, `
		UPDATE customers SET
		name = $1,
		phone = $2
		WHERE id = $3
		RETURNING id,
    	name,
    	phone,
    	active,
		created;`, name, phone, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

//RemoveByID ...
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx, `
		DELETE FROM customers WHERE id = $1 RETURNING id, name, phone, active, created;
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

//BlockByID ...
func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx, `
		UPDATE customers
		SET
		active = FALSE
		WHERE id = $1
		RETURNING
		id, name, phone, active, created;
		`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

//UnBlockByID ...
func (s *Service) UnBlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.db.QueryRowContext(ctx, `
		UPDATE customers
		SET
		active = TRUE
		WHERE id = $1
		RETURNING
		id, name, phone, active, created;
		`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}