package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	capitalLetterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numericRunes       = []rune("0123456789")
)

// GetUUID ... Generate UUID but on os just install uuidgen
func GetUUID() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}

	r := strings.NewReplacer("\n", "", "-", "")
	return r.Replace(string(out)), nil
}

// GetRefID ... Generate reference id not dupplicate
func GetRefID() (string, error) {
	uuid, err := GetUUID()
	if err != nil {
		return "", err
	}

	uuid = string(uuid)

	t := time.Now()
	out := t.Format("20060102150405") + uuid[len(uuid)-6:]

	return strings.ToUpper(out), nil
}

// BytesToString ... convert byte array to string
func BytesToString(data []byte) string {
	if data == nil {
		return ""
	}

	return string(data[:])
}

// StringToBytes ... convert string to byte array
func StringToBytes(data string) []byte {
	return []byte(data)
}

// StringToUnixTimestamp ... convert string to unix timestamp
func StringToUnixTimestamp(format, datetime string) (int64, error) {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	t, err := time.Parse(format, datetime)
	return t.Unix(), err
}

// StringToUnixTime ... convert string to time
func StringToUnixTime(format, datetime string) (time.Time, error) {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	t, err := time.ParseInLocation(format, datetime, time.Local)
	return t, err
}

// StringToUnixTime ... convert string to time
func StringUnixToTime(unix string) (time.Time, error) {
	i, err := strconv.ParseInt(unix, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}

// TimestampToString ... convert string to unix timestamp
func TimestampToString(format string, sec int64) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	t := time.Unix(sec, 0)
	return t.Format(format)
}

// TimestampToStringWithLocation ... convert string to unix timestamp
func TimestampToStringWithLocation(format string, sec int64, local string) (string, error) {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	t := time.Unix(sec, 0)
	loc, err := time.LoadLocation(local)
	if err != nil {
		return "", err
	}

	return t.In(loc).Format(format), nil
}

// ConvertYearBuddistToChrist ... convert string in year Buddhist to Christ only format d/m/y
// validateYear is check (year - new_year) <= 0 in Buddhist or Christ (year must not lessthan now - 543)
func ConvertYearBuddhistToChrist(strDate string, validateYear bool) (string, error) {
	sp := strings.Split(strings.TrimSpace(strDate), " ")
	if len(sp) > 0 {

		spDate := strings.Split(sp[0], "/")
		if len(spDate) >= 3 {
			year, err := strconv.Atoi(spDate[2])
			if err != nil {
				return strDate, err
			}

			month, err := strconv.Atoi(spDate[1])
			if err != nil {
				return strDate, err
			}

			day, err := strconv.Atoi(spDate[0])
			if err != nil {
				return strDate, err
			}

			if validateYear {
				if (year - time.Now().Year()) <= 0 {
					return strDate, nil //Year is Christ
				}
			}

			year -= 543
			return fmt.Sprintf("%02d/%02d/%d", day, month, year), nil
		}
	}

	return strDate, fmt.Errorf("ConvertYearBuddhistToChrist cannot convert %s", strDate)
}

// ConvertYearChristToBuddhist ... convert string in year Buddhist to Christ only format d/m/y
// validateYear is check (year - new_year) <= 0 in Buddhist or Christ (year must not lessthan now - 543)
func ConvertYearChristToBuddhist(strDate string, validateYear bool) (string, error) {
	sp := strings.Split(strings.TrimSpace(strDate), " ")
	if len(sp) > 0 {

		spDate := strings.Split(sp[0], "/")
		if len(spDate) >= 3 {
			year, err := strconv.Atoi(spDate[2])
			if err != nil {
				return strDate, err
			}

			month, err := strconv.Atoi(spDate[1])
			if err != nil {
				return strDate, err
			}

			day, err := strconv.Atoi(spDate[0])
			if err != nil {
				return strDate, err
			}

			if validateYear {
				if (year - time.Now().Year()) <= 0 {
					year += 543
				}
			} else {
				year += 543
			}

			return fmt.Sprintf("%02d/%02d/%d", day, month, year), nil
		}
	}

	return strDate, fmt.Errorf("ConvertYearChristToBuddhist cannot convert %s", strDate)
}

// ConvertMDYToYMD convert MM/dd/yyyy to dd/MM/yyy
func ConvertMDYToYMD(strDate string) (string, error) {
	sp := strings.Split(strings.TrimSpace(strDate), " ")
	if len(sp) > 0 {

		spDate := strings.Split(sp[0], "/")
		if len(spDate) >= 3 {
			year, err := strconv.Atoi(spDate[2])
			if err != nil {
				return strDate, err
			}

			month, err := strconv.Atoi(spDate[0])
			if err != nil {
				return strDate, err
			}

			day, err := strconv.Atoi(spDate[1])
			if err != nil {
				return strDate, err
			}

			return fmt.Sprintf("%02d/%02d/%d", day, month, year), nil
		}
	}

	return strDate, fmt.Errorf("ConvertMMDDYYYYToDDMMYYYY cannot convert %s", strDate)
}

// make JSONMarshal for don't want unicode when convert
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// InterfaceToStruct ... convert interface{] to struct
func InterfaceToStruct(src, des interface{}) error {
	if src == nil {
		return fmt.Errorf("InterfaceToStruct src is nil")
	}

	byteData, err := JSONMarshal(src)
	if err != nil {
		return fmt.Errorf("InterfaceToStruct json Marshal error: %v", err)
	}

	err = json.Unmarshal(byteData, des)
	if err != nil {
		return fmt.Errorf("InterfaceToStruct json Unmarshal error: %v", err)
	}

	return nil
}

// InterfaceToString ... convert interface{] to string
func InterfaceToString(src interface{}) (string, error) {
	if src == nil {
		return "", fmt.Errorf("InterfaceToString src is nil")
	}

	byteData, err := JSONMarshal(src)
	if err != nil {
		return "", fmt.Errorf("InterfaceToString json Marshal error: %v", err)
	}

	return BytesToString(byteData), nil
}

// StringToStruct ... convert String to struct
func StringToStruct(src string, des interface{}) error {
	return BytesToStruct(StringToBytes(src), des)
}

// BytesToStruct ... convert []Byte to struct
func BytesToStruct(src []byte, des interface{}) error {
	err := json.Unmarshal(src, des)
	if err != nil {
		return fmt.Errorf("BytesToStruct json Unmarshal error: %v", err)
	}

	return nil
}

// StructToBytes ... convert struct to []Byte
func StructToBytes(src interface{}) ([]byte, error) {
	if src == nil {
		return nil, fmt.Errorf("StructToBytes src is nil")
	}

	jsonByte, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("StructToBytes json Marshal error: %v", err)
	}

	return jsonByte, nil
}

// GetSizeFromInterface ... get size of byte from interface{}
func GetSizeFromInterface(src interface{}) int {
	if src == nil {
		return 0
	}

	jsonByte, err := StructToBytes(src)
	if err != nil {
		return 0
	}

	return len(jsonByte)
}

// GetStringFromSQL ... get string when query sqldb
func GetStringFromSQL(val sql.NullString) string {
	if val.Valid {
		return val.String
	}

	return ""
}

// GetIntFromSQL ... get int when query sqldb
func GetIntFromSQL(val sql.NullInt32) int {
	if val.Valid {
		return int(val.Int32)
	}

	return 0
}

// GetTimeFromSQL ... get time when query sqldb
func GetTimeFromSQL(val sql.NullTime) *time.Time {
	if val.Valid {
		strTime, _ := TimestampToStringWithLocation("", val.Time.Unix(), "UTC")
		// loc, _ := time.LoadLocation("Asia/Bangkok")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local)
		return &t
	}

	return nil
}

// GetBoolFromSQL ... get boolean when query sqldb
func GetBoolFromSQL(val sql.NullBool) bool {
	if val.Valid {
		return val.Bool
	}

	return false
}

func GetCurrencyFromString(strCurrency string) *big.Float {
	strCurrency = strings.ReplaceAll(strCurrency, ",", "")
	n := new(big.Float)
	n, ok := n.SetString(strCurrency)
	if !ok {
		return big.NewFloat(0)
	}

	return n
}

// stringInSlice ... check has value in array
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IntInSlice ... check has value in array
func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// RandCapitalString ... random capital letter n characters
func RandCapitalString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = capitalLetterRunes[rand.Intn(len(capitalLetterRunes))]
	}
	return string(b)
}

// RandNumberString ... random capital letter n characters
func RandNumberString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numericRunes[rand.Intn(len(numericRunes))]
	}
	return string(b)
}

// HMACSha1 encrypt hmac sha1
func HMACSha1(payload string, key string) string {
	keyForSign := []byte(key)
	data := []byte(payload)
	h := hmac.New(sha1.New, keyForSign)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HMACSha256 encrypt hmac sha256
func HMACSha256(payload string, key string) string {
	keyForSign := []byte(key)
	data := []byte(payload)
	h := hmac.New(sha256.New, keyForSign)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// WeekRange get start/end time from week of year
func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

// WeekStart get start time from week of year
func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

// MonthInterval get firstday/lastday time from month of year
func MonthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}
