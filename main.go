package main

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/xingliuhua/gostrings/model"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

//go:embed utils.go.txt
var src string

func main() {
	fmt.Println("will generate r directory \nkey will be translate to camel style...")

	err := generateStringResource()
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
	fmt.Println("generate r success!")
}

func generateStringResource() error {
	err := createRFile()
	if err != nil {
		return err
	}
	err = createUtilFile()
	if err != nil {
		return err
	}
	stringReses, err := parseAllXML()
	if err != nil {
		return err
	}
	allKey := make(map[string]string)
	allStringData := make(map[string]map[string]string)
	allStringArrayData := make(map[string]map[string][]string)
	for local, stringRes := range stringReses {
		localKeys := make(map[string]string)
		for _, v := range stringRes.StringArrays {
			if _, exist := localKeys[strcase.ToCamel(v.Name)]; exist {
				return errors.New("Duplicate name:" + v.Name)
			}
			localKeys[strcase.ToCamel(v.Name)] = strcase.ToCamel(v.Name)
			allKey[strcase.ToCamel(v.Name)] = strcase.ToCamel(v.Name)
			strings := make([]string, 0)
			for _, string := range v.Items {
				strings = append(strings, string)
			}
			if allStringArrayData[local] == nil {
				allStringArrayData[local] = make(map[string][]string)
			}
			allStringArrayData[local][strcase.ToCamel(v.Name)] = strings
		}
		for _, v := range stringRes.Strings {
			if _, exist := localKeys[strcase.ToCamel(v.Name)]; exist {
				return errors.New("Duplicate name:" + v.Name)
			}
			localKeys[strcase.ToCamel(v.Name)] = strcase.ToCamel(v.Name)
			allKey[strcase.ToCamel(v.Name)] = strcase.ToCamel(v.Name)
			if allStringData[local] == nil {
				allStringData[local] = make(map[string]string)
			}
			allStringData[local][strcase.ToCamel(v.Name)] = v.Text
		}
	}
	err = writeKeyData(allKey)
	if err != nil {
		return err
	}
	err = writeInitData(allStringData, allStringArrayData)
	if err != nil {
		return err
	}

	exec.Command("bash", "-c", "go fmt ./r/r.go").Run()
	return nil
}
func writeInitData(stringsData map[string]map[string]string, stringArrayData map[string]map[string][]string) error {
	r, err := os.OpenFile("./r/r.go", os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer r.Close()
	bufferString := bytes.NewBufferString("func init() {\n")
	_, err = bufferString.WriteString("allString=map[string]map[string]string{\n")
	if err != nil {
		return err
	}

	for local, v := range stringsData {
		_, err = bufferString.WriteString("\"" + local + "\": {\n")
		if err != nil {
			return err
		}
		for k, v2 := range v {
			_, err = bufferString.WriteString("\"" + k + "\":\"" + v2 + "\",\n")
			if err != nil {
				return err
			}
		}

		_, err = bufferString.WriteString("},\n")
		if err != nil {
			return err
		}
	}
	_, err = bufferString.WriteString("}\n")
	if err != nil {
		return err
	}

	_, err = bufferString.WriteString("allStringArray = map[string]map[string][]string{\n")
	if err != nil {
		return err
	}
	for local, v := range stringArrayData {
		_, err = bufferString.WriteString("\"" + local + "\": {\n")
		if err != nil {
			return err
		}

		for k, v2 := range v {
			_, err = bufferString.WriteString("\"" + k + "\": []string{\n")
			if err != nil {
				return err
			}
			for _, v3 := range v2 {
				_, err = bufferString.WriteString("\"" + v3 + "\",\n")
				if err != nil {
					return err
				}
			}
			_, err = bufferString.WriteString("},\n")
			if err != nil {
				return err
			}
		}
		_, err = bufferString.WriteString("},\n")
		if err != nil {
			return err
		}
	}
	_, err = bufferString.WriteString("}\n}\n")
	if err != nil {
		return err
	}

	_, err = r.WriteString(bufferString.String())
	if err != nil {
		return err
	}
	return nil
}
func writeKeyData(allKeys map[string]string) error {
	r, err := os.OpenFile("./r/r.go", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer r.Close()
	bufferString := bytes.NewBufferString("")
	_, err = bufferString.WriteString("package r\n")
	if err != nil {
		return err
	}
	_, err = bufferString.WriteString("const (\n")
	if err != nil {
		return err
	}
	for k, v := range allKeys {
		_, err = bufferString.WriteString(k + "= \"" + v + "\"\n")
		if err != nil {
			return err
		}
	}
	_, err = bufferString.WriteString(")\n")
	if err != nil {
		return err
	}
	_, err = r.WriteString(bufferString.String())
	if err != nil {
		return err
	}
	return nil
}
func createRFile() error {
	file, err := os.Open("./r")
	defer file.Close()
	if exist := os.IsNotExist(err); exist {
		err := os.Mkdir("r", 0766)
		if err != nil {
			return err
		}
	}
	_, err = os.Open("./r/r.go")
	if err == nil {
		err := os.Remove("./r/r.go")
		if err != nil {
			return err
		}
	}
	return nil
}
func createUtilFile() error {
	file, err := os.Open("./r")
	defer file.Close()
	if exist := os.IsNotExist(err); exist {
		err := os.Mkdir("r", 0766)
		if err != nil {
			return err
		}
	}
	_, err = os.Open("./r/util.go")
	if err == nil {
		err := os.Remove("./r/util.go")
		if err != nil {
			return err
		}
	}
	utilFile, err := os.OpenFile("./r/util.go", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer utilFile.Close()
	bufferString := bytes.NewBufferString(src)
	_, err = utilFile.WriteString(bufferString.String())
	if err != nil {
		return err
	}
	return nil
}
func parseAllXML() (map[string]model.StringRes, error) {

	r, err := os.Open("./strings")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	readdir, err := r.Readdir(-1)
	if err != nil {
		return nil, err
	}
	m := make(map[string]model.StringRes)
	filterName, err := regexp.Compile("string.*\\.xml")
	if err != nil {
		return nil, err
	}
	subName, err := regexp.Compile("-.*\\.xml")
	if err != nil {
		return nil, err
	}
	for _, v := range readdir {
		if v.IsDir() {
			continue
		}
		b := filterName.MatchString(v.Name())
		if !b {
			continue
		}
		matchString := subName.FindString(v.Name())
		res, err := parseXML(v.Name())
		if err != nil {
			return nil, err
		}
		if matchString == "" {
			m[""] = res
		} else {
			matchString = matchString[1 : len(matchString)-4]
			m[matchString] = res
		}
	}
	return m, nil
}
func parseXML(fileName string) (model.StringRes, error) {
	res := model.StringRes{}

	file, err := os.Open("./strings/" + fileName)
	if exist := os.IsNotExist(err); exist {
		return res, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return res, nil
	}
	err = xml.Unmarshal(data, &res)
	if err != nil {
		return res, nil
	}
	return res, nil
}
