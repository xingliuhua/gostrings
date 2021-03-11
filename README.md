[中文版](https://github.com/xingliuhua/gostrings/blob/master/README.cn.md)
# gramod

This is an International Library of go strings. It can automatically generate go files according to XML resource files, which is convenient for unified string management.
## Background
* Hard coding of string in code is not convenient to manage.
* There is a need for multilingualism.

## Feature
* String resources are put in the strings folder for easy management.
* XML configuration string, easy to reuse, management, multi language comparison.
* Before compiling, the user calls the command to automatically generate go file, which is more efficient at runtime.

## Install
go get github.com/xingliuhua/gostrings

## Usage
1. Create the strings folder at the root of the project.
2. Create an. XML file in the strings folder (it must be in sting [language].XML format, and the abbreviation of the language is unlimited); for example string.xml 、string- en.xml 、string-en- us.xml 、string- zh.xml 、string- unknown.xml and so on.
3. the normal string and string array can be placed in the XML file. For specific format, please refer to:
string.xml
```xml
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="confirm">confirm</string>
    <string name="cancel">cancel</string>
    <string_array name="city">
        <item>beijing</item>
        <item>new York</item>
        <item>Hong Kong</item>
    </string_array>
</resources>
```
string-zh.xml
```xml
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="confirm">确定</string>
    <string name="cancel">取消</string>
    <string_array name="city">
        <item>北京</item>
        <item>纽约</item>
        <item>香港</item>
    </string_array>
</resources>
```
4. Switch to the root directory of the project and execute the "gostrings" command (remember to add gopath / bin to the path environment variable), which will automatically generate the R folder and the go files in it. Do not move the files in it.
5. in code:
```go
	str, err := r.GetString("", r.Cancel) // from string.xml
	str, err := r.GetString("unknown", r.Cancel) // from string-unknown.xml
	str, err := r.GetString("zh", r.Cancel) // from string-zh.xml
	strArray, err := r.GetStringArray("en-us", r.City) // from string-en-us.xml
```


## Maintainers

[@xingliuhua](https://github.com/xingliuhua).

## Contributing

Feel free to dive in! [Open an issue] or submit PRs.
