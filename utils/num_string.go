package utils

import "github.com/jackc/pgx/v5/pgtype"

func numericToString(numeric pgtype.Numeric) string {
	// use the String method of *big.Int to convert to string
	strValue := numeric.Int.String()

	// if there's a scale, insert the decimal point
	if numeric.Exp < 0 {
		decimalPosition := len(strValue) + int(numeric.Exp)
		strValue = strValue[:decimalPosition] + "." + strValue[decimalPosition:]
	}

	return strValue
}
