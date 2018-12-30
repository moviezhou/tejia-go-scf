package main

import (
    "context"
	"fmt"
	"bytes"
	"encoding/xml"
	"html/template"
	"os"
	"io"
	"io/ioutil"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
    "github.com/tencentyun/scf-go-lib/cloudfunction"
)

const APPID = "1253927884"
const SecretId = "AKIDCkrAPEgkonSxhRHsLiWIxqnJmAHZtlFf"
const SecretKey = "2BEEFLRdfbPJyhkpl06DUjpmZ3mUH2eA"
const Bucket = "tejia-1253927884"
const Region = "ap-chengdu"
const Host = "https://tejia-1253927884.cos.ap-chengdu.myqcloud.com"
const TmplPath = "https://tejia-1253927884.cos.ap-chengdu.myqcloud.com/web/templates/"
const SCFTempPath = "/tmp/"


type DefineEvent struct {
    // test event define
    Key1 string `json:"key1"`
    Key2 string `json:"key2"`
}

type ListBucketResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Name string `xml:"Name"`
	Prefix string `xml:"Prefix"`
	Marker string `xml:"Marker"`
	MaxKeys uint32 `xml:"MaxKeys"`
	IsTruncated bool `xml:"IsTruncated"`
	Contents   []Contents   `xml:"Contents"`
}

type Contents struct {
	XMLName xml.Name `xml:"Contents"`
	Key    string   `xml:"Key"`
	LastModified    string   `xml:"LastModified"`
	ETag  string   `xml:"ETag"`
	Size uint32 `xml:"Size"`
	Owner Owner `xml:"Owner"`
	StorageClass string `xml:"StorageClass"`
}

type Response struct{ 
    IsBase64 bool `json:"IsBase64"` 
    StatusCode uint32 `json:"statusCode"` 
    Headers map[string]string `json:"headers"` 
	Body string `json:"body"` 
} 

type Owner struct {
	XMLName  xml.Name `xml:"Owner"`
	ID  string   `xml:"ID"`
	DisplayName  string   `xml:"DisplayName"`
}


func getXML(url string) ([]byte, error) {
	url = url + "?prefix=xitaihua_kjc"
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status Error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

//  Get go template by it's name, if it's not in local /tmp path then download it from Tecent cloud COS bucket.
func getTemplate(tmpl string) (*template.Template, error) {
	if _,err := os.Stat(SCFTempPath + tmpl); !os.IsNotExist(err) {
		t := template.New(tmpl)
		return t, nil
	} else {
		tmplUrl := TmplPath + tmpl
		fileName := SCFTempPath + tmpl
		err := downloadFile(fileName, tmplUrl)
		if err != nil{
			fmt.Println(err.Error())
			t := template.New("errorTmpl")
			return t, err
		} else {
			t := template.New(tmpl)
			if err != nil {
				fmt.Println(err.Error())
			}
			return t, nil
		}
	}
}

func downloadFile(filePath string, url string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func hello(ctx context.Context, event DefineEvent) (Response, error) {
	fmt.Println("key1:", event)
	fmt.Println(ctx)
	fmt.Println("key2:", event.Key2)


	db, err := sql.Open("mysql", "root:JsJs6773@tcp(cd-cdb-hhl9rz48.sql.tencentcdb.com:63394)/tejia?charset=utf8")
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("Test", "X部门", "2018-12-28")
	defer stmt.Close()
	checkErr(err)
	fmt.Println(res.RowsAffected())
	defer db.Close()
	
	byteValue, _ := getXML(Host)
	var objects ListBucketResult
	
	xml.Unmarshal(byteValue, &objects)
	
	tmpl, err := getTemplate("index.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		ret := Response{} 
		ret.IsBase64 = false 
		ret.StatusCode = 200 
		ret.Headers = map[string]string{} 
		ret.Headers["Content-Type"] = "text/html" 
		ret.Body = "Template error."
		return ret, nil
	}

	tmpl, err = tmpl.ParseFiles(SCFTempPath + "index.tmpl")
	if err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println("object key:", objects.Contents[1].Key)

    ret := Response{} 
    ret.IsBase64 = false 
    ret.StatusCode = 200 
    ret.Headers = map[string]string{} 
    ret.Headers["Content-Type"] = "text/html" 
    
	// Escape the first object which is the folder name
	objects.Contents = objects.Contents[1:len(objects.Contents)] 

	var tpl bytes.Buffer
	tmpl.Execute(&tpl, objects)
	ret.Body = tpl.String()
    return ret, nil 
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
    // Make the handler available for Remote Procedure Call by Cloud Function
    cloudfunction.Start(hello)
}