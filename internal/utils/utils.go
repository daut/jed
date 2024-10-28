package utils

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

// converts a float64 to a pgtype.Numeric
func ConvertToPGNumeric(number float64) (*pgtype.Numeric, error) {
	value := &pgtype.Numeric{}
	parse := strconv.FormatFloat(number, 'f', -1, 64)
	if err := value.Scan(parse); err != nil {
		return nil, err
	}
	return value, nil
}

func StrPtr(s string) *string {
	return &s
}
