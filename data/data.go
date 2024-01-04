package data

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func Hello() string {
	fmt.Println("hello from utils")
	return "hello from utils"
}

func ParseJSONObjectFromFile(filename string, obj interface{}) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	byteValue, _ := io.ReadAll(jsonFile)
	err = ParseJSONObject(string(byteValue), obj)
	return err
}

func ParseJSONFromFile(filename string) (interface{}, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var result interface{}
	err = ParseJSONObject(string(byteValue), &result)
	return result, nil
}

func ParseJSONObject(jsond string, obj interface{}) error {
	byteValue := []byte(jsond)
	err := json.Unmarshal([]byte(byteValue), &obj)
	return err
}

func ParseJSON(jsond string) (interface{}, error) {
	byteValue := []byte(jsond)
	var result interface{}
	json.Unmarshal([]byte(byteValue), &result)
	return result, nil
}

func HashAndSalt(password string) (string, error) {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHash(hashedPwd string, password string) error {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(password))
	return err
}

// pass source and dest as reference
func CopyStruct(source interface{}, dest interface{}) error {
	err := copier.Copy(dest, source)
	if err != nil {
		// log.Fatal(err)
		return err
	}
	return nil
}

// timezone = from the IANA Time Zone database
// US/Pacific, US/Central, US/Mountain, US/Eastern
func TimeIn(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return t, err
	}
	t = t.In(loc)
	return t, err
}

func ToString(arg interface{}) string {
	if arg == nil {
		return ""
	}
	return fmt.Sprintf("%v", arg)
}

func RandomInt(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	if max > 0 {
		return r1.Intn(max)
	}
	return r1.Int()
}

func GetMD5sum(filename string) (string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new MD5 hash object
	hash := md5.New()

	// Read the file contents into the MD5 hash object
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	// Calculate the MD5 checksum of the file contents
	checksum := hash.Sum(nil)

	// Convert the MD5 checksum to a string
	checksumString := hex.EncodeToString(checksum)

	return checksumString, nil
}

type CustomWriter struct {
	Builder strings.Builder
}

func (cw *CustomWriter) Write(p []byte) (n int, err error) {
	l, err := cw.Builder.WriteString(string(p))
	//fmt.Println(l, err)
	fmt.Print(string(p))
	return l, err
}

func EndsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}

	return s[len(s)-len(suffix):] == suffix
}
