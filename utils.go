package goutils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func Hello() string {
	fmt.Println("hello from utils")
	return "hello from utils"
}

func ParseJSONObject(filename string) (interface{}, error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	byteValue, _ := io.ReadAll(jsonFile)
	var result interface{}
	json.Unmarshal([]byte(byteValue), &result)
	jsonFile.Close()

	return result, nil
}

func ParseJSONFromFile(filename string) (map[string]interface{}, error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	byteValue, _ := io.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	jsonFile.Close()

	return result, nil
}

func ParseJSON(jsond string) (map[string]interface{}, error) {
	byteValue := []byte(jsond)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	return result, nil
}

func HashAndSalt(password string) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		// log.Fatal(err)
		return "", nil
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func CompareHash(hashedPwd string, password string) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(password))
	if err != nil {
		// log.Fatal(err)
		return false, nil
	}

	return true, nil
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
func TimeIn(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// log.Fatal(err)
		return t, err
	}
	t = t.In(loc)
	return t, err
}

func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	// Read the contents of the file
	contents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(contents), nil
}

func FileOrDirExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		//permission denied or some other error
		// log.Fatal(err)
		return false, err
	}
}

func WriteFile(filename string, content string, overwrite bool) error {
	exists, err := FileOrDirExists(filename)
	if overwrite == false && exists == true {
		return errors.New("file " + filename + " exists")
	}

	f, err := os.Create(filename)

	if err != nil {
		// log.Fatal(err)
		return err
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Fatal(err2)
		return err2
	}

	return nil
}

func IsFile(filename string) (bool, error) {

	fi, err := os.Stat(filename)
	if err != nil {
		// log.Fatal(err)
		return false, err
	}
	return fi.Mode().IsRegular(), nil
}

func IsDirectory(dirname string) (bool, error) {

	fi, err := os.Stat(dirname)
	if err != nil {
		// log.Fatal(err)
		return false, err
	}
	return fi.Mode().IsDir(), nil
}

func DeleteFile(filename string) error {

	isf, err := IsFile(filename)
	if err != nil {
		// log.Fatal(err)
		return err
	}

	if isf {
		err = os.Remove(filename)
		if err != nil {
			// log.Fatal(err)
			return err
		}
		return nil
	}
	return errors.New("not a file")
}

func DeleteDirectory(dirname string) error {
	isDir, err := IsDirectory(dirname)
	if err != nil {
		// log.Fatal(err)
		return err
	}

	if isDir {
		err = os.Remove(dirname)
		if err != nil {
			// log.Fatal(err)
			return err
		}
		return nil
	}
	return errors.New("not a directory")
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

func CopyFileFast(src string, dst string) error {

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile1 copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir1 recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	if exists, err := FileOrDirExists(dst); err == nil {
		if !exists {
			//create directory
			err = os.MkdirAll(dst, si.Mode())
			if err != nil {
				return err
			}
		}
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Type().IsRegular() {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyDir1(source string, destination string) error {

	//fix the seperator
	source = filepath.FromSlash(source)
	destination = filepath.FromSlash(destination)
	if exists, err := FileOrDirExists(destination); err == nil {
		if !exists {
			//create directory
			os.Mkdir(destination, os.ModePerm)
		}
	}

	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}

		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {

			err1 := CopyFileFast(filepath.Join(source, relPath), filepath.Join(destination, relPath))
			if err1 != nil {
				return err1
			}
		}
		return nil
	})
	return err
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

func GetFileName(path string) string {
	filename := filepath.Base(path)
	return filename
}

func GetDirectoryName(path string) string {
	dir := filepath.Dir(path)
	if dir == "" {
		dir = filepath.Base(path)
	}
	return dir
}

func GetCurrentDirectory() (string, error) {
	exePath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// Get the directory containing the Go executable.
	return exePath, nil
}

type customWriter struct {
	builder strings.Builder
}

func (cw *customWriter) Write(p []byte) (n int, err error) {
	l, err := cw.builder.WriteString(string(p))
	//fmt.Println(l, err)
	fmt.Print(string(p))
	return l, err
}

func ExecCMD(workingdir string, args ...string) (string, error) {

	log.Printf("ExecCMD %v", args)

	cwOut := customWriter{}
	cwOut.builder = strings.Builder{}
	cwErr := customWriter{}
	cwErr.builder = strings.Builder{}

	baseCmd := args[0]
	cmdArgs := args[1:]

	var cmd = exec.Command(baseCmd, cmdArgs...)
	cmd.Dir = workingdir
	cmd.Stdout = &cwOut
	cmd.Stderr = &cwErr

	err := cmd.Run()
	if err != nil {
		log.Printf("ExecCMD failed with %s\n", err)
	}
	output := fmt.Sprintf("%s %s", cwOut.builder.String(), cwErr.builder.String())
	return output, nil
}

func HTTPGetBody(url string) (string, error) {

	// Send an HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Convert the response body to a string
	return string(body), nil

}

func HTTPGetCode(url string) (int, error) {

	response, err := http.Head(url)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	status := response.Status
	statusCode := response.StatusCode

	fmt.Printf("URL: %s\n", url)
	fmt.Printf("HTTP Status: %s\n", status)
	fmt.Printf("Status Code: %d\n", statusCode)
	return statusCode, nil
}
func HTTPostJson(url string, jsond string) (string, error) {

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(jsond)))
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}
	return string(body), nil
}

func GetDirectoryEntries(dir string) ([]fs.DirEntry, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entries, nil
}

func EndsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}

	return s[len(s)-len(suffix):] == suffix
}

func GetPathSeperator() string {
	path, err := GetCurrentDirectory()

	if err != nil {
		return string(filepath.Separator)
	}

	if strings.Index(path, "/") > -1 {
		return "/"
	}

	if strings.Index(path, "\\") > -1 {
		return "\\"
	}
	return string(filepath.Separator)
}

func CreateDirectory(dirPath string) error {

	return os.Mkdir(dirPath, 0755)
}
