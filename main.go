package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	//"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"strings"
	//"github.com/xuri/excelize/v2"
)

type VersionRecord struct {
	// 1. Create a struct for storing CSV lines and annotate it with JSON struct field tags
	Sysname  string `json:"sysname"`
	Unitname string `json:"unitname"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	Version  string `json:"version"`
	Env      string `json:"env"`
}

func createShoppingList(data [][]string) []VersionRecord {
	// convert csv lines to array of structs
	var shoppingList []VersionRecord
	for i, line := range data {
		//if len(line) != 7 {
		//	continue
		//}
		if i > 0 { // omit header line
			var rec VersionRecord
			for j, field := range line {
				rec.Env = "production"
				switch j {
				case 0:
					rec.Sysname = strings.ToLower(field)
				case 1:
					rec.Unitname = strings.ToLower(field)
				case 2:
					rec.Hostname = field
				case 3:
					rec.Ip = field
				case 7:
					rec.Version = field
				}
			}
			shoppingList = append(shoppingList, rec)
		}
	}

	return shoppingList
}

func Parsecsv(csvfile string) map[string]map[string]string {
	//subMap := make(map[string]string)
	mainMap := make(map[string]map[string]string)
	// open file
	f, err := os.Open(csvfile)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// 2. Read CSV file using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// 3. Assign successive lines of raw CSV data to fields of the created structs
	shoppingList := createShoppingList(data)

	// 4. Convert an array of structs to JSON using marshaling functions from the encoding/json package
	//jsonData, err := json.MarshalIndent(shoppingList, "", "  ")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(string(jsonData))
	//fmt.Printf("type:%T \n", shoppingList)
	//fmt.Printf("type:%T \n", string(jsonData))
	for _, item := range shoppingList {
		if item.Unitname != "" && item.Sysname != "" {
			// 初始化多维map
			if _, exist := mainMap[item.Sysname]; exist {
				// 如果map值不存在，初始化下
				if mainMap[item.Sysname][item.Unitname] == "" {
					mainMap[item.Sysname][item.Unitname] = item.Version
				} else {
					if mainMap[item.Sysname][item.Unitname] < item.Version {
						mainMap[item.Sysname][item.Unitname] = item.Version
					}
				}
			} else {
				//map初始化
				c := make(map[string]string)
				c[item.Unitname] = item.Version
				mainMap[item.Sysname] = c
			}
		}
	}

	return mainMap
}

//func checkError(message string, err error) {
//	if err != nil {
//		log.Fatal(message, err)
//	}
//}

//func findversion(fline string) string {
//	re := regexp.MustCompile(`.*_[A-Z0-9]{5}_AP-[A-Z]-.*-.*-.*\.tar\.gz$`)
//	fileName:= re.FindString("WPA010_COGED_AP-A-0.9.2-210816-1439.HF-1030-A.tar.gz")
//
//	sampleRegexp := regexp.MustCompile(`\d+\.\d+\.\d+`)
//
//	version := sampleRegexp.FindAllString(fileName,1)
//	//fmt.Printf("version: %s\n", version)
//	return version[0]
//}
// sit的特例版本
//func getsitversion() map[string]string{
//	f, err := excelize.OpenFile("1024versionchange.xlsx")
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rv:= make(map[string]string) // Get all the rows in the Sheet1 section.
//	rows, err := f.GetRows("Sheet1")
//	for _, row := range rows {
//		if len(row) >=4{
//			//fmt.Println(row[0],row[2])
//			rv[strings.ToLower(row[0])]=rv[strings.ToLower(row[2])]
//		}
//	}
//	return rv
//}

func main() {
	file1 := flag.String("file1", "", "版本文件一")
	file2 := flag.String("file2", "", "版本文件二")
	flag.Parse()
	// open file
	//sitvs:=make(map[string]string)
	fmt.Printf("var1:%s\n", *file1)
	fmt.Printf("var2:%s\n", *file2)
	firstResult := Parsecsv(*file1)
	secontResult := Parsecsv(*file2)
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	headerstr := []string{"sysName", "unitName", "firstfile", "secondfile"}
	writer.Write(headerstr)
	defer writer.Flush()
	//sitvs=getsitversion()
	for sysName, m := range firstResult {
		for unitName, _ := range m {
			values := []string{}

			//// sit 环境由于版本变更0.9.* 其实是1.0.*,如果非sit环境，可以去掉
			//if _, ok := sitvs[unitName]; ok {
			//	if _, exist := secontResult[sysName][unitName]; exist {
			//		if secontResult[sysName][unitName] < "1.0.0" {
			//			secontResult[sysName][unitName] = "1.0.0"
			//		}
			//	}
			//}
			//// end

			values = append(values, sysName, unitName, firstResult[sysName][unitName], secontResult[sysName][unitName])
			err := writer.Write(values)
			checkError("Cannot write to file", err)
		}
	}

	//fmt.Println(findversion("/home/ap/appoper/WPA010_COGED_AP-A-0.9.2-210816-1439.HF-1030-A.tar.gz"))
}
