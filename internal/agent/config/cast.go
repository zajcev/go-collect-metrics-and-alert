package config

import "github.com/spf13/cast"

func GetString(s any) string {
	return cast.ToString(s)
}

func GetUint(s any) uint64 {
	return cast.ToUint64(s)
}
