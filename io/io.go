package io

import (
	"errors"
	"fmt"
	"goutils/data"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func CopyFileFast(srcPath string, dstPath string) error {

	sourceFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dstPath)
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

func GetDirectoryEntries(dir string) ([]fs.DirEntry, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entries, nil
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

func ExecCMD(workingdir string, args ...string) (string, error) {

	log.Printf("ExecCMD %v", args)

	cwOut := data.CustomWriter{}
	cwOut.Builder = strings.Builder{}
	cwErr := data.CustomWriter{}
	cwErr.Builder = strings.Builder{}

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
	output := fmt.Sprintf("%s %s", cwOut.Builder.String(), cwErr.Builder.String())
	return output, nil
}
