package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/xingliuhua/gostrings/internal/model"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	logo := "                                   ___                                                          \n                                 ,--.'|_             ,--,                                       \n              ,---.              |  | :,'   __  ,-.,--.'|         ,---,                         \n  ,----._,.  '   ,'\\   .--.--.   :  : ' : ,' ,'/ /||  |,      ,-+-. /  |  ,----._,.  .--.--.    \n /   /  ' / /   /   | /  /    '.;__,'  /  '  | |' |`--'_     ,--.'|'   | /   /  ' / /  /    '   \n|   :     |.   ; ,. :|  :  /`./|  |   |   |  |   ,',' ,'|   |   |  ,\"' ||   :     ||  :  /`./   \n|   | .\\  .'   | |: :|  :  ;_  :__,'| :   '  :  /  '  | |   |   | /  | ||   | .\\  .|  :  ;_     \n.   ; ';  |'   | .; : \\  \\    `. '  : |__ |  | '   |  | :   |   | |  | |.   ; ';  | \\  \\    `.  \n'   .   . ||   :    |  `----.   \\|  | '.'|;  : |   '  : |__ |   | |  |/ '   .   . |  `----.   \\ \n `---`-'| | \\   \\  /  /  /`--'  /;  :    ;|  , ;   |  | '.'||   | |--'   `---`-'| | /  /`--'  / \n .'__/\\_: |  `----'  '--'.     / |  ,   /  ---'    ;  :    ;|   |/       .'__/\\_: |'--'.     /  \n |   :    :            `--'---'   ---`-'           |  ,   / '---'        |   :    :  `--'---'   \n  \\   \\  /                                          ---`-'                \\   \\  /              \n   `--`-'                                                                  `--`-'               \n"
	fmt.Println(logo)
	fmt.Println("generate r directory...")

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
	//err = createUtilFile()
	//if err != nil {
	//	return err
	//}
	stringReses, err := parseAllXML()
	if err != nil {
		return err
	}
	allKey := make(map[string]string)
	allStringData := make(map[string]map[string]string)
	allStringArrayData := make(map[string]map[string][]string)
	for local, stringRes := range stringReses {
		if allStringData[local] == nil {
			allStringData[local] = make(map[string]string)
		}
		for _, v := range stringRes.Strings {
			if _, exist := allStringData[local][v.Name]; exist {
				return errors.New("Duplicate name:" + v.Name)
			}
			allKey[wrapStringKey(v.Name)] = v.Name

			allStringData[local][v.Name] = v.Text
		}

		if allStringArrayData[local] == nil {
			allStringArrayData[local] = make(map[string][]string)
		}
		for _, v := range stringRes.StringArrays {

			if _, exist := allStringArrayData[local][v.Name]; exist {
				return errors.New("Duplicate name:" + v.Name)
			}

			allKey[wrapArrayKey(v.Name)] = v.Name
			strings := make([]string, 0)
			for _, string := range v.Items {
				strings = append(strings, string)
			}
			allStringArrayData[local][v.Name] = strings
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
func wrapStringKey(key string) string {
	return "String_" + key
}
func wrapArrayKey(key string) string {
	return "Array_" + key
}
func writeInitData(stringsData map[string]map[string]string, stringArrayData map[string]map[string][]string) error {
	r, err := os.OpenFile("./r/r.go", os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer r.Close()
	bufferString := bytes.NewBufferString("func init() {\n")
	_, err = bufferString.WriteString("allString := map[string]map[string]string{\n")
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

	_, err = bufferString.WriteString("allStringArray := map[string]map[string][]string{\n")
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
	_, err = bufferString.WriteString("}\n")
	if err != nil {
		return err
	}
	_, err = bufferString.WriteString("strutil.SetData(allString,allStringArray)")
	if err != nil {
		return err
	}
	_, err = bufferString.WriteString("}\n")
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
	_, err = bufferString.WriteString("package r\nimport \"github.com/xingliuhua/gostrings/pkg/strutil\"\n")
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
