package publicfunc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEncry(t *testing.T) {
	f, err := os.Open("F:\\project.7z")
	if err != nil {
		return
	}
	defer f.Close()
	bgtm := time.Now()
	sha256str, _ := GetFileSha256f(f)
	fmt.Println(bgtm, "===", time.Now())
	bgtm = time.Now()
	sha1str, _ := GetFileSha1f(f)
	fmt.Println(bgtm, "===", time.Now())
	fmt.Println(sha256str, sha1str)

	f2, err := os.Open("./test/测试 - 副本.txt")
	if err != nil {
		return
	}
	defer f2.Close()
	sha256str2, _ := GetFileSha256f(f2)
	sha1str2, _ := GetFileSha1f(f2)

	fmt.Println(sha256str2, sha1str2)
}

func TestSearchFile(t *testing.T) {
	filenames, err := filepath.Glob("F:\\project_git\\dsp-fileplugin\\fileinfo_collecter\\*\\*")

	if err != nil {
		return
	}
	fmt.Println(filenames)
}
