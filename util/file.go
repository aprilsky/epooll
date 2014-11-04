package util

import (
	"bufio"
	"errors"
	"io"
    "io/ioutil"
    "log"
	"os"
	"path/filepath"
	"regexp"
)

var PATH_ROOT = SelfDir()
var PATH_DATA = SelfDir()+"/data"

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

/**
 *  使用方法
    if ok, err := util.WriteLog("Just a test\n"); !ok {
        log.Print(err)
    }
 */
func WriteLog(filename string, format string) (bool, error) {

    logfile := PATH_DATA+"/log/"+filename+".log";
    f, err := os.OpenFile(logfile, os.O_RDWR | os.O_APPEND |  os.O_CREATE, 0777)
    if err != nil {
        return false, err
    }
    defer f.Close() 
    logger := log.New(f, "", log.Ldate | log.Ltime | log.Lshortfile)
    logger.Print(format)
    return true, err
}

/**
 *  使用方法
    if ok, err := util.PutFile("/data/golang/log/go.txt", "Just a test\n", 1); !ok {
        log.Print(err)
    }
 */
func PutFile(file string, format string, args ...interface{}) (bool, error) {

    f, err := os.OpenFile(file, os.O_RDWR | os.O_APPEND |  os.O_CREATE, 0777)
    // 上面的0777并不起作用
    os.Chmod(file, 0777)
    // 如果没有传参数，重新新建文件
    if args == nil {
        f, err = os.Create(file)
    }
    for _, arg := range args {
        // 参数为0，也重新创建文件
        if arg == 0 {
            f, err = os.Create(file)
        }
    }

    if err != nil {
        return false, err
    }

    // 要先检查nil是否为空，才能关闭打开的文件，否则报错
    // http://stackoverflow.com/questions/16280176/go-panic-runtime-error-invalid-memory-address-or-nil-pointer-dereference
    defer f.Close()

    f.WriteString(format)
    return true, err
    //f.Write([]byte("Just a test!\r\n"))
}

func GetFile(file string) (string, error) {
    
    f, err := os.Open(file)
    if err != nil {
        // 抛出异常
        //panic(err)
        return "", err
    }
    defer f.Close() 
    // 这里不用处理错误了，如果是文件不存在或者没有读权限，上面都直接抛异常了，这里还可能有错误么？
    fd, _  := ioutil.ReadAll(f)
    return string(fd), err
}

// Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	isLongLine := false
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}