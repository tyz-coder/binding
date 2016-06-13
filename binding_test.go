package binding

import (
	"fmt"
	"testing"
	"time"
	"strconv"
)

type Human struct {
	Name     string    `form:"name"`
	Age      int       `form:"age"`
	Birthday time.Time `form:"birthday"`
}

func (this *Human) CleanedBirthday(n string) (time.Time, error) {
	return time.Parse("2006-01-02", n)
}

type Student struct {
	Human
	Number int    `form:"number"`
	Class  string `form:"class"`
}

//func (this *Student) DefaultClass() string {
//	return "Class one"
//}


func (this *Human) DefaultAge() int {
	return 100
}

func (this *Student) DefaultAge() int {
	return 200
}

func (this *Human) CleanedNumber(n string) (int, error) {
	var num, e = strconv.Atoi(n)
	return num+100, e
}

//func (this *Student) CleanedNumber(n string) (int, error) {
//	var num, e = strconv.Atoi(n)
//	return num+200, e
//}

var formData = map[string]interface{}{"name": "Yangfeng", "number": "9", "birthday": "2016-06-12"}

func TestSample(t *testing.T) {
	// 绑定
	var s *Student
	var err = BindWithTag(formData, &s, "form")
	if err != nil {
		fmt.Println("绑定失败")
	}
	fmt.Println(s)
}

//
//type MyString string
//
//////////////////////////////////////////////////////////////////////////////////////////////////////
//type Human struct {
//	CleanedData map[string]interface{}
//	Name        MyString  `binding:"name"`
//	Age         int       `binding:"age"`
//	Birthday    time.Time `binding:"birthday"`
//	List        []int     `binding:"list"`
//}
//
//func (this *Human) CleanedName(n string) (MyString, error) {
//	if len(n) > 0 {
//		return MyString(fmt.Sprintf("My name is %s", n)), nil
//	}
//	return "", errors.New("随便给点吧")
//}
//
//func (this *Human) CleanedBirthday(n string) (time.Time, error) {
//	return time.Parse("2006-01-02", n)
//}
//
//////////////////////////////////////////////////////////////////////////////////////////////////////
//type Class struct {
//	ClassName string `binding:"class_name"`
//}
//
//func (this *Class) DefaultClassName() string {
//	return "class 3"
//}
//
//////////////////////////////////////////////////////////////////////////////////////////////////////
//type Student struct {
//	Human
//	Number int `binding:"number"`
//	Class  Class
//}
//
//var source = map[string]interface{}{"list": []string{"123", "456"}, "name": "SmartWalle", "age": 123.5, "birthday": "2016-06-12", "number": 1234, "class_name_1": "class 1"}
//
//func TestBindPoint(t *testing.T) {
//	var s *Student
//	fmt.Println(Bind(source, &s))
//	if s != nil {
//		fmt.Println(s.CleanedData)
//		fmt.Println(s.Name, s.Age, s.Birthday, s.Number, s.Class.ClassName, s.List)
//	}
//}
