package goutils

import (
	"fmt"
	"goutils/data"
	"regexp"
	"testing"
	"time"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHello(t *testing.T) {

	fmt.Printf("%v \n", "Begin TestHello")

	want := regexp.MustCompile(`\b` + "hello from utils" + `\b`)
	msg := data.Hello()
	if want.MatchString(msg) == false {
		t.Fatalf(`%v does not match %v`, msg, want)
	}

	fmt.Printf("%v \n", "End TestHello")

}

type TestJSONObj struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func TestParseJSONObject(t *testing.T) {

	jsond := "{ \"name\" : \"namevalue\", \"id\" : 5 }"
	var testObj TestJSONObj
	err := data.ParseJSONObject(jsond, &testObj)
	if err != nil {
		t.Fatal(err)
	}
	//var testObj = obj.(TestJSONObj)

	if testObj.Id != 5 {
		t.Fatalf(`%v does not match %v`, testObj.Id, 5)
	}

	if testObj.Name != "namevalue" {
		t.Fatalf(`%v does not match %v`, testObj.Name, "namevalue")
	}

}

func TestParseJSON(t *testing.T) {

	jsond := "{ \"name\" : \"namevalue\", \"id\" : 5 }"

	obj, err := data.ParseJSON(jsond)
	if err != nil {
		t.Fatal(err)
	}

	testObj := obj.(map[string]interface{})

	//json numbers are parsed as float64
	if testObj["id"].(float64) != 5 {
		t.Fatalf(`%v does not match %v`, testObj["id"], 5)
	}

	if testObj["name"].(string) != "namevalue" {
		t.Fatalf(`%v does not match %v`, testObj["name"], "namevalue")
	}

}

func TestHashAndSalt(t *testing.T) {
	password := "secret1234"
	result, err := data.HashAndSalt(password)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Error("HashAndSalt did not generate a string")
	}
	fmt.Println(result)
}

func TestCompareHash(t *testing.T) {

	password := "secret1234"
	hash, err := data.HashAndSalt(password)
	if err != nil {
		t.Fatal(err)
	}

	err = data.CompareHash(hash, password)
	if err != nil {
		t.Fatal(err)
	}

	err = data.CompareHash("1111", password)
	if err == nil {
		t.Fatal("error should be thrown for invalid hash")
	}

	//generate new hash, and compare to old password should return false
	newpassword := "secret12345"
	hash, err = data.HashAndSalt(newpassword)
	if err != nil {
		t.Fatal(err)
	}

	err = data.CompareHash(hash, password)
	if err == nil {
		t.Fatal("error should be thrown on hash not matching")
	}

}

func TestCopyStruct(t *testing.T) {

	var struct1 TestJSONObj
	var struct2 TestJSONObj

	struct1.Name = "name"
	struct1.Id = 10

	err := data.CopyStruct(&struct1, &struct2)
	if err != nil {
		t.Fatal(err)
	}

	if struct1.Name != struct2.Name {
		t.Fatalf(`%v does not match %v`, struct1.Name, struct2.Name)
	}

	if struct1.Id != struct2.Id {
		t.Fatalf(`%v does not match %v`, struct1.Id, struct2.Id)
	}

}

func TestTimeIn(t *testing.T) {
	cdate := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("cdate:", cdate)

	ctime, err := data.TimeIn(cdate, "US/Central")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("ctime:", ctime)
	if ctime.Year() != 2023 || int(ctime.Month()) != 12 || ctime.Day() != 31 || ctime.Hour() != 18 {
		t.Fatalf("invalid date %v %v %v %v", ctime.Year(), int(ctime.Month()), ctime.Day(), ctime.Hour())
	}

	// DayLights saving, Sun, Mar 10, 2024 – Sun, Nov 3, 2024 2AM
	centralLocation, err := time.LoadLocation("US/Central")
	if err != nil {
		fmt.Println("Error loading timezone:", err)
	}

	//2:00 AM on daylights savings will resulting 1:00 AM
	cdate = time.Date(2024, time.March, 10, 1, 59, 59, 59, centralLocation)
	fmt.Println("cdate:", cdate)

	//2:00 AM on daylights savings will resulting 1:00 AM
	cdate = time.Date(2024, time.March, 10, 2, 0, 0, 0, centralLocation)
	fmt.Println("cdate:", cdate)

	//offset will change based on daylights savings
	cdate = time.Date(2024, time.March, 10, 3, 0, 0, 0, centralLocation)
	fmt.Println("cdate:", cdate)

}
