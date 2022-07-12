package config

import (
	"os"
)

type TypeEnv int

const (
	TypeEnvDev TypeEnv = iota
	TypeEnvTest
	TypeEnvProd
	TypeEnvUnknown
)

func (te TypeEnv) String() string {
	switch te {
	case TypeEnvDev:
		return "dev"
	case TypeEnvTest:
		return "test"
	case TypeEnvProd:
		return "prod"
	default:
		return "unknown"
	}
}

func GetEnv() TypeEnv {
	envStr, ok := os.LookupEnv("ENV")
	if !ok {
		return TypeEnvUnknown
	}

	return convertStringToTypeEnv(envStr)
}

func SetEnv(e TypeEnv) {
	os.Setenv("ENV", e.String())
}

func convertStringToTypeEnv(str string) TypeEnv {
	switch str {
	case "dev":
		return TypeEnvDev
	case "test":
		return TypeEnvTest
	case "prod":
		return TypeEnvProd
	default:
		return TypeEnvUnknown
	}
}
