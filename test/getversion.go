package main

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
)

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
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

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
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("uat_version.csv")
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	headerstr := []string{"ip", "filename", "version"}
	writer.Write(headerstr)
	defer writer.Flush()
	var sList []UatVersion
	for i, line := range data {
		values := []string{}
		if i > 0 { // omit header line
			var rec UatVersion
			for j, field := range line {
				rec.Env="pre"
				switch j {
				case 1:
					rec.Ip=field
				case 2:
					rec.Filename=field
					if len(strings.Split(field,"_")) == 3{
						rec.Identy=strings.Split(field,"-")[0]+"_"+rec.Ip
						rec.PyName=strings.Split(field,"_")[0]
						rec.UnitName=strings.Split(field,"_")[1]
						rec.AppType=strings.Split(strings.Split(field,"_")[2],"-")[0]
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

func getversionfile(filename string) map[string]string  {
	mainMap1 := make(map[string]string)
	shoppingList:= createshoplist(filename)
	for _, item := range shoppingList {
		if item.Identy != "" {
			// 初始化多维map
			if _, exist := mainMap1[item.Identy]; exist {
				// 如果map值不存在，初始化下
				if mainMap1[item.Identy] == "" {
					mainMap1[item.Identy] = item.Version
				} else {
					if mainMap1[item.Identy] < item.Version {
						mainMap1[item.Identy] = item.Version
					}
				}
			} else {
				//map初始化
				mainMap1[item.Identy] = item.Version
			}
		}
	}
	return mainMap1
}

func main() {
	firstResult:=getversionfile("hxpre_pkgversion_2022012710.csv")
	// open file
	file, err := os.Create("uat_result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	headerstr := []string{"sysName", "unitName", "firstfile", "secondfile"}
	writer.Write(headerstr)
	defer writer.Flush()
	//sitvs=getsitversion()
	for id, v := range firstResult {
		values := []string{}
		writer.Write(append(values, id,v))
	}

}
