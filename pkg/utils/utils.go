package utils

import (
    "os"
    "strconv"
    "time"
)

func UnixTimeAsFloat() float64 {
    return float64(time.Now().UnixNano()) / 1e9
}

func IntFromEnv(varname string, defaultValue int) int {
    v := os.Getenv(varname)
    if v == "" {
        return defaultValue
    }
    i, err := strconv.Atoi(v)
    if err != nil {
        panic("Env var " + varname + " is expected to be integer!")
    }
    return i
}
