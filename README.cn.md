[English](https://github.com/xingliuhua/gostrings/blob/master/README.md)
# gostrings

这是一个go字符串国际化库，根据xml资源文件自动生成go文件，便于字符串统一管理。

## 背景
* 代码中对字符串硬编码不方便管理。
* 有多语言的需求。
## 功能特点
* 字符串资源统一放在strings文件夹中，便于管理。
* xml配置字符串，便于复用、管理、多语言对比。
* 编译前用户调用命令，自动生成go文件，运行时效率更好


## 安装
执行go get github.com/xingliuhua/gostrings

## 使用
1. 在项目根目录创建strings文件夹。
2. 在strings文件夹中创建.xml文件(必须是sting[语言].xml格式，语言简称不限制);比如string.xml、string-en.xml、string-en-us.xml、string-zh.xml、string-unknown.xml等。
3. .xml文件中可以放普通字符串和字符串数组。具体格式可以参考：
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
4. 切换到项目根目录下执行"gostrings"命令(记得把gopath/bin加入到path环境变量中），会自动生成r文件夹和里面的go文件，不要动里面的文件。如果找不到该命令请检查gopath的bin路径是否加入到环境变量中。
```shell script
gostrings
```
5. 代码中查询
```go
	import "github.com/xingliuhua/gostrings/pkg/gostrs"
    
    str, err := gostrs.ShouldGetString("", r.Cancel) // from string.xml
    str := gostrs.GetString("unknown", r.Cancel) // from string-unknown.xml
    str := gostrs.GetStringWithDefault("zh", r.Cancel) // from string-zh.xml
    strArray, err := gostrs.ShouldGetStringArray("en-us", r.City) // from string-en-us.xml
```
在实际的开发中，可以根据http参数或Accept-Language头字段选择语言。
## 维护

[@xingliuhua](https://github.com/xingliuhua).

## 贡献

Feel free to dive in! [Open an issue] or submit PRs.
