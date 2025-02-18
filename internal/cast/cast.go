package cast

import "github.com/spf13/cast"

func GetString(s any) string {
	return cast.ToString(s)
}

func GetUint(s any) uint64 {
	return cast.ToUint64(s)
}

func GetFloat(s any) float64 {
	return cast.ToFloat64(s)
}

func GetInt64(s any) int64 {
	return cast.ToInt64(s)
}
