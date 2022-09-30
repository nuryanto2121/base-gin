package util

import (
	"app/pkg/logging"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// WIB :
	WIB string = "Asia/Jakarta"
	// UTC :
	UTC string = "UTC"
)

// GetTimeNow :
func GetTimeNow() time.Time {
	return time.Now().In(GetLocation())
}

// GetLocation - get location wib
func GetLocation() *time.Location {
	var logger = logging.Logger{}
	// return time.FixedZone(WIB, 7*3600)
	loc, err := time.LoadLocation(UTC)
	if err != nil {
		logger.Fatal(err)
	}
	return loc
}

// GetUnixTimeNow :
func UnixNow() int64 {
	return time.Now().UnixMilli()
}

//NameStruct :
func NameStruct(myvar interface{}) string {
	name := ""
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}
	return fmt.Sprintf("%s%s", strings.ToLower(name[:1]), name[1:])
}

// Stringify :
func Stringify(data interface{}) string {
	dataByte, _ := json.Marshal(data)
	return string(dataByte)
}

func CheckEmail(e string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 3 && len(e) > 254 {
		return false
	}

	return emailRegex.MatchString(e)
}

func GenerateNumber(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func GetSecond(startDate int64, endDate int64) float64 {
	stDate := time.Unix(startDate/1000, 0)
	edDate := time.Unix(endDate/1000, 0)
	diff := edDate.Sub(stDate)

	return diff.Seconds()
}
func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func IntToString(val int) string {
	return strconv.Itoa(val)
}
func Int64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}
func StringToInt(val string) int {
	res, _ := strconv.Atoi(val)
	return res
}

func Int64ToTime(val int64) time.Time {
	return time.Unix(val, 0)
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
