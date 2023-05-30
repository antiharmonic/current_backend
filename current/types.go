package current

import (
	"database/sql"
	"encoding/json"
	"strconv"
)

type NullInt struct {
	sql.NullInt64
}

type NullString struct {
	sql.NullString
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(`null`), nil
}

func (n NullInt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return []byte(`null`), nil
}

func ParamToNullString(s string) (NullString, error) {
	var ns NullString
	if s != "" {
		ns.String = s
		ns.Valid = true
	}
	return ns, nil
}

func StringParamToNullInt(p string, is_zero_nil bool) (NullInt, error) {
	var ni NullInt
	if p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return ni, err
		}
		ni, err = IntParamToNullInt(n, is_zero_nil)
		if err != nil {
			return ni, err
		}
	}
	return ni, nil
}

func IntParamToNullInt(p int, is_zero_nil bool) (NullInt, error) {
	var ni NullInt
	if is_zero_nil && p == 0 {
		ni.Valid = false
	} else {
		ni.Int64 = int64(p)
		ni.Valid = true
	}
	return ni, nil
}