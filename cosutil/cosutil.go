package cosutil

import (
	"net/http"
	"io"
	"os"
	"io/ioutil"
	"html/template"
	"fmt"
)

const APPID = "1253927884"
const SecretId = "AKIDCkrAPEgkonSxhRHsLiWIxqnJmAHZtlFf"
const SecretKey = "2BEEFLRdfbPJyhkpl06DUjpmZ3mUH2eA"
const Bucket = "tejia-1253927884"
const Region = "ap-chengdu"
const Host = "https://tejia-1253927884.cos.ap-chengdu.myqcloud.com"
const TmplPath = "https://tejia-1253927884.cos.ap-chengdu.myqcloud.com/web/templates/"
const SCFTempPath = "/tmp/"


func GetXML(url string) ([]byte, error) {
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
func GetTemplate(tmpl string) (*template.Template, error) {
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
