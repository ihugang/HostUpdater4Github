package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const url = "https://raw.hellogithub.com/hosts"

var regex = regexp.MustCompile(`# GitHub520 Host Start\n\r.*# GitHub520 Host End\n\r`)

var hostFile = strings.Replace(os.Getenv("windir"), "\\", "/", -1) + "/system32/drivers/etc/hosts"
var hostFile1 = "hosts"

func main() {
	fmt.Println("Auto update dns hosts about github for rapid visit.")
	fmt.Println("Coded by shrek (390652@qq.com)")
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	hostFile1 = filepath.Join(filepath.Dir(path), "hosts.txt")
	hostFile1 = strings.Replace(hostFile, "\\", "/", -1)

	hostsContent := ""
	dat, err := os.ReadFile(hostFile1)
	if err != nil {
		fmt.Println(err)
	} else {
		hostsContent = string(dat)
		fmt.Println(hostsContent)
	}

	updateHosts(hostsContent)
	fmt.Println("Press the Enter Key to exit")
	fmt.Scanln()
}

// 更新hosts
func updateHosts(hostsContent string) {
	content := fetch(url)
	newContent := regex.ReplaceAllString(hostsContent, content)
	if !regex.MatchString(content) {
		newContent = hostsContent + "\n\r\n\r" + content
	}
	fmt.Println(newContent)
	ptr, err := syscall.UTF16PtrFromString(hostFile)
	if err != nil {
		panic(err)
	}

	if attr, err := syscall.GetFileAttributes(ptr); err == nil {
		if attr&syscall.FILE_ATTRIBUTE_READONLY != 0 {
			fmt.Println("Host file is read only.")
			return
		}
	}

	fd, err := os.OpenFile(hostFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModeAppend|os.ModePerm)

	if err != nil {
		panic(err)
	}

	defer fd.Close()
	fd.WriteString(newContent)

}

func fetch(url string) string {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	secs := time.Since(start).Seconds()
	fmt.Printf("\n\r%.2fs %7d %s\n\r", secs, len(b), url)
	return string(b)
}
