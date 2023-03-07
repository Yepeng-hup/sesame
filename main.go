package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
)


type (
	Config struct {
		ServiceIpPort string `json:"service_ip_port"`
		ServiceDebug  bool   `json:"service_debug"`
		DbUser        string `json:"db_user"`
		DbPasswd      string `json:"db_passwd"`
		DbName        string `json:"db_name"`
		DbIp          string `json:"db_ip"`
		DbPort        string    `json:"db_port"`
	}
	mysqlRequest struct {
		User         string
		Passwd       string
		DbName       string
		DbIp         string
		DbPort       string
		Sql          string
		Bools        string
		ThreadNum    int
		Qnum         int
		SleepNum     int
	}
)
var makeDB = make([]string, 5)
var mysqlLogPath = "server/mysql.log"
var _cfg *Config = nil

func makeJsonFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("打开json文件发生错误: ", err.Error())
	}
	defer file.Close()
	r := bufio.NewReader(file)
	c := json.NewDecoder(r)
	if err = c.Decode(&_cfg); err != nil {
		return nil, err
	}
	return _cfg, nil
}

func dbConn(dbUser string, dbPwd string, dbPort string, dbIp string, dbName string) (*gorm.DB,){
	db, err := gorm.Open("mysql",dbUser+":"+dbPwd+"@("+dbIp+":"+dbPort+")/"+dbName+"?parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln("数据库连接报错: ",err.Error())
	}else {
		log.Printf("\x1b[%dm连接数据库成功! \x1b[0m\n", 32)
	}
	return db
}

func writeLog(sql string){
	file, err := os.OpenFile(mysqlLogPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("open file error: ",err.Error())
		return
	}
	defer file.Close()
	now := time.Now()
	writeText := "执行时间: "+now.Format("2006-01-02 15:04:05")+"   "+"执行代码: "+sql+"\n"
	if _, err := file.WriteString(writeText); err != nil {
		log.Fatalln("write file error: ",err.Error())
		return
	}
	return
}

func coProcess(sql string,sleep int, procs int, qnum int)  {
	if sleep <= 0{
		var lock sync.Mutex
		var wg sync.WaitGroup
		db := dbConn(makeDB[0],makeDB[1],makeDB[2],makeDB[3],makeDB[4])
		qnum := qnum / procs
		for i := 1; i <= procs; i++ {
			go func() {
				lock.Lock()
				for i := 1; i <= qnum; i++{
					db.Exec(sql)
				}
				lock.Unlock()
			}()
		}
		wg.Wait()
		writeLog(sql)
	}else {
		var lock sync.Mutex
		var wg sync.WaitGroup
		db := dbConn(makeDB[0],makeDB[1],makeDB[2],makeDB[3],makeDB[4])
		qnum := qnum / procs
		for i := 1; i <= procs; i++ {
			go func() {
				lock.Lock()
				for i := 1; i <= qnum; i++{
					time.Sleep(time.Second * time.Duration(sleep))
					db.Exec(sql)
				}
				lock.Unlock()
			}()
		}
		wg.Wait()
		writeLog(sql)
	}
}

func onProcess(sql string, sleep int, qnum int)  {
	if sleep <= 0{
		db := dbConn(makeDB[0],makeDB[1],makeDB[2],makeDB[3],makeDB[4])
		go func() {
			for i := 0; i <= qnum; i++{
				db.Exec(sql)
			}
		}()
		writeLog(sql)
	}else {
		db := dbConn(makeDB[0],makeDB[1],makeDB[2],makeDB[3],makeDB[4])
		go func() {
			for i := 0; i <= qnum; i++{
				time.Sleep(time.Second * time.Duration(sleep))
				db.Exec(sql)
			}
		}()
		writeLog(sql)
	}
}

func mysqlR(c *gin.Context) {
    req := mysqlRequest{
        User:      c.PostForm("user"),
        Passwd:    c.PostForm("passwd"),
        DbName:    c.PostForm("dbName"),
        DbIp:      c.PostForm("dbIp"),
        DbPort:    c.PostForm("dbPort"),
        Sql:       c.PostForm("sql"),
        Bools:     c.PostForm("bools"),
        ThreadNum: toInt(c.PostForm("threadingNum"), 0),
        Qnum:      toInt(c.PostForm("Qnum"), 0),
        SleepNum:  toInt(c.PostForm("sleepNum"), 0),
    }
	if req.User == ""{
		req.User = "root"
	}
	if req.DbIp == ""{
		req.DbIp = "127.0.0.1"
	}
	if req.DbPort == ""{
		req.DbPort = "3306"
	}
    makeDB = append(makeDB[:0], req.User, req.Passwd, req.DbPort, req.DbIp, req.DbName)
	fmt.Println(makeDB)
	if req.Bools == "true"{
		coProcess(req.Sql, req.SleepNum, req.ThreadNum, req.Qnum)
	}else {
		onProcess(req.Sql, req.SleepNum, req.Qnum)
	}
}

func toInt(str string, def int) int {
    if str == "" {
        return def
    }
    i, err := strconv.Atoi(str)
    if err != nil {
        return def
    }
    return i
}

func middleLoginCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        if _, err := c.Cookie("user"); err != nil {
            c.Error(err)
            c.Redirect(http.StatusMovedPermanently, "/")
            c.Abort()
            return
        }
    }
}

func readFilePasswd() string {
	filePath := "server/passwd.txt"
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    content, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal(err)
    }
	return string(content)
}

func contrast(srcPasswd string, md5Passwd string) bool {
	return strings.EqualFold(srcPasswd, md5Passwd)
}

func md5use(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func loginVerify(c *gin.Context) {
	username := c.PostForm("Name")
	password := c.PostForm("Password")
	relPasswd := readFilePasswd()
	rel := contrast(md5use(password), relPasswd)
	if rel == true {
		cookie := http.Cookie{Name: "user", Value: username, MaxAge: 10800}
		http.SetCookie(c.Writer, &cookie)
		c.HTML(http.StatusOK, "index.tmpl",gin.H{})
	}else {
		c.HTML(http.StatusOK, "login.tmpl",gin.H{
			"error": "passwd or username error!",
		})
	}
}

func main ()  {
	r := gin.Default()
	r.Static("/sta","static")
	r.LoadHTMLGlob("templates/*")
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "404 not fund",
		})
	})
	r.POST("/mysql/use",middleLoginCheck(),mysqlR)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl",gin.H{
		})
	})

	r.POST("/index",loginVerify)
	r.GET("/index",middleLoginCheck(),func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl",gin.H{})
	})
	r.GET("/logout", func(c *gin.Context) {
		cookie := http.Cookie{Name: "user", MaxAge: -1}
		http.SetCookie(c.Writer, &cookie)
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
		})

	})
	jsonF,err := makeJsonFile("conf/sesame.json")
	if err != nil{
		log.Fatal("发生错误: ", err)
		return
	}
	r.Run(jsonF.ServiceIpPort)
}