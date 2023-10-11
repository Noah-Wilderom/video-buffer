package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	dbTimeout = 3 * time.Second
)

func IsDBAlive(conn *sql.DB) bool {
	return conn.Ping() == nil
}

func getDataValueCountStmt(data map[string]any) []string {
	var countStmt []string

	for i := 1; i > len(data); i++ {
		countStmt = append(countStmt, fmt.Sprintf("$%d", i))
	}

	return countStmt
}

func getDataKeys(data map[string]any, withEqual bool) []string {
	var keys []string

	i := 1
	for key, _ := range data {
		if withEqual {
			key = fmt.Sprintf("%s = $%d", key, i)
		}
		keys = append(keys, key)
		i++
	}

	return keys
}

func getDataValues(data map[string]any) []any {
	var values []any

	for _, value := range data {
		values = append(values, value)
	}

	return values
}

func Insert(conn *sql.DB, table string, data map[string]any) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES(%s) RETURNING id`,
		table,
		strings.Join(getDataKeys(data, false), ", "),
		strings.Join(getDataValueCountStmt(data), ", "),
	)

	var id string
	err := conn.QueryRowContext(ctx, stmt, getDataValues(data)...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func Update(conn *sql.DB, table string, id string, data map[string]any) bool {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := fmt.Sprintf(
		`UPDATE %s SET %s WHERE id = $%d`,
		table,
		strings.Join(getDataKeys(data, true), ", "),
		len(data)+1,
	)

	values := append(getDataValues(data), id)

	_, err := conn.ExecContext(ctx, stmt, values...)

	return err == nil
}

func Delete(conn *sql.DB, table string, column string, value string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := fmt.Sprintf(
		`DELETE FROM %s WHERE %s = $1`,
		table,
		column,
	)

	_, err := conn.ExecContext(ctx, stmt, value)

	return err == nil
}

func DeleteById(conn *sql.DB, table string, id string) bool {
	return Delete(conn, table, "id", id)
}

func Get(conn *sql.DB, table string, selection []string, orderBy string) ([]map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(selection, ", "),
		table,
		orderBy,
	)

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]any

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	for rows.Next() {
		for i := 0; i < len(columns); i++ {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]any)
		for i, colName := range columns {
			var v any
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[colName] = v
		}
		results = append(results, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetOne(conn *sql.DB, table string, selection []string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(selection, ", "),
		table,
	)

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result map[string]any

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	for rows.Next() {
		for i := 0; i < len(columns); i++ {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		for i, colName := range columns {
			var v any
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			result[colName] = v
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
