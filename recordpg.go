package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type UatVersion struct {
	// 1. Create a struct for storing CSV lines and annotate it with JSON struct field tags
	Identy string `json:"identyid"`
	PyName string `json:"pyName"`
	UnitName string `json:"unitname"`
	AppType string `json:"apptype"`
	Ip  string `json:"ip"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
	Env      string `json:"env"`
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func findversion(fline string) string {
	re := regexp.MustCompile(`.*_[A-Z0-9]{5}_(AP|WB|DB)-[A-Z]-.*-.*-.*\.tar\.gz$`)
	fileName:= re.FindString(fline)

	sampleRegexp := regexp.MustCompile(`\d+\.\d+\.\d+`)

	version := sampleRegexp.FindAllString(fileName,1)
	//fmt.Printf("version: %s\n", version)
	if len(version) >0{
		return version[0]
	} else {
		return "0.0.0"
	}

}
func Split(r rune) bool {
	return r == '-' || r == '_'
}

func createshoplist(filename string) []UatVersion {
	//f, err := os.Open("hxuat_pkgversion_2022012017.csv")
	f,err:=os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()
	// 2. Read CSV file using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// 获取环境
	var envName = strings.Split(filepath.Base(filename), "_")[0]
	file, err := os.Create("uat_version.csv")
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	headerstr := []string{"ip", "filename", "version"}
	writer.Write(headerstr)
	defer writer.Flush()
	var sList []UatVersion
	for i, line := range data {
		if len(line) != 6 {
			continue
		}
		values := []string{}
		if i > 0 { // omit header line
			var rec UatVersion
			for j, field := range line {
				rec.Env=envName
				switch j {
				case 1:
					rec.Ip=field
				case 2:
					rec.Filename=field
					if len(strings.Split(field,"_")) == 3{
						rec.Identy=strings.Split(field,"-")[0]+"_"+rec.Ip
						rec.PyName=strings.Split(field,"_")[0]
						rec.AppType=strings.FieldsFunc(field,Split)[2]
						rec.UnitName=strings.Split(field,"_")[1]+"_"+rec.AppType
					} else {
						rec.PyName="null"
						rec.UnitName="null"
						rec.AppType="null"
					}

					rec.Version=findversion(field)
				}
			}
			values=append(values,rec.Identy,rec.PyName,rec.UnitName,rec.AppType,rec.Ip,rec.Filename,rec.Version,rec.Env)
			err = writer.Write(values)
			checkError("Cannot write to file", err)
			sList=append(sList,rec)
		}
	}
	return sList
}

// db connect

var db *sql.DB
var Confpath string = "./pg.conf"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func sqlInsert(filename string) {
	var err error
	//filename:="hxuat_pkgversion_2022012017.csv"
	//还可以是这种方式打开
	connetStr, err := ioutil.ReadFile("./pg.conf")
	//fmt.Println(connetStr)
	connetline:=strings.Replace(string(connetStr),"\n","",-1)
	if err != nil {
		fmt.Println("read fail", err)
	}
	fmt.Printf("%s?%s\n",connetline,"sslmode=disable")
	db, err := sql.Open("postgres", fmt.Sprintf("%s?%s",connetline,"sslmode=disable"))
	checkErr(err)
	//插入数据
	//tmstp:=time.Now().Unix()
	stmt, err := db.Prepare("INSERT INTO versions_history(sysid,unitname,ip,pkgname,vsion,environment,timestamp ) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id")
	checkErr(err)

	// 获取时间戳
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	tmstp := tm.Format("2006-01-02 03:04:05 PM")

	total := int64(0)
	shoppingList:= createshoplist(filename)
	for _, item := range shoppingList {
		if item.Identy != "" {
			log.Println("insert",item.PyName,item.UnitName,item.Ip,item.Filename,item.Version,item.Env,tmstp)
			res, err := stmt.Exec(item.PyName,item.UnitName,item.Ip, item.Filename,item.Version,item.Env,tmstp)
			//这里的三个参数就是对应上面的$1,$2,$3了
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			total+=affect

		}
	}
	fmt.Println("rows affect:", total)
}

func recordtodb(filename string) (int,error)  {
	sqlInsert(filename)
	return 0,nil
}

//需要赋值的变量
var version = ""
var fileName = ""
//通过flag包设置-version参数
var printVersion bool

func init() {
	flag.BoolVar(&printVersion, "version", false, "print program build version")
	flag.StringVar(&fileName,"file1", "", "版本文件")
	flag.Parse()
}

func main() {
	if printVersion {
		println(version)
		os.Exit(0)
	}
	recordtodb(fileName)
}
