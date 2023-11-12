package coalescing

import (
	"fmt"
	"strconv"
)

func ToBool(val interface{}) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return len(v) > 0, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot coalesce %T into bool", val)
	}
}

func ToFloat64(val interface{}) (float64, error) {
	var result float64

	switch v := val.(type) {
	case bool:
		if v {
			result = 1
		} else {
			result = 0
		}
	case int64:
		result = float64(v)
	case float64:
		result = v
	case nil:
		result = 0
	default:
		return 0, fmt.Errorf("cannot coalesce %T into float64", val)
	}

	return result, nil
}

func Int64Compatible(val interface{}) bool {
	switch val.(type) {
	case bool:
		return true
	case int64:
		return true
	case nil:
		return true
	default:
		return false
	}
}

func ToInt64(val interface{}) (int64, error) {
	var result int64

	switch v := val.(type) {
	case bool:
		if v {
			result = 1
		} else {
			result = 0
		}
	case int64:
		result = v
	case nil:
		result = 0
	default:
		return 0, fmt.Errorf("cannot coalesce %T into int64", val)
	}

	return result, nil
}

func ToString(val interface{}) (string, error) {
	var result string

	switch v := val.(type) {
	case bool:
		result = strconv.FormatBool(v)
	case int64:
		result = strconv.FormatInt(v, 10)
	case float64:
		result = fmt.Sprintf("%f", v)
	case nil:
		result = "null"
	case string:
		result = v
	default:
		return "", fmt.Errorf("cannot coalesce %T into string", val)
	}

	return result, nil
}
