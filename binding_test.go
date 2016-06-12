package binding

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type MyString string

////////////////////////////////////////////////////////////////////////////////////////////////////
type Human struct {
	CleanedData map[string]interface{}
	Name        MyString  `model:"name"`
	Age         int       `model:"age"`
	Birthday    time.Time `model:"birthday"`
}

func (this *Human) CleanedName(n string) (MyString, error) {
	if len(n) > 0 {
		return MyString(fmt.Sprintf("My name is %s", n)), nil
	}
	return "", errors.New("随便给点吧")
}

func (this *Human) CleanedBirthday(n string) (time.Time, error) {
	return time.Parse("2006-01-02", n)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Class struct {
	ClassName string `model:"class_name"`
}

func (this *Class) DefaultClassName() string {
	return "class 3"
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Student struct {
	Human
	Number int `model:"number"`
	Class  Class
}

var source = map[string]interface{}{"name": "SmartWalle", "age": 123.5, "birthday": "2016-06-12", "number": 1234, "class_name_1": "class 1"}

func TestBindPoint(t *testing.T) {
	var s *Student
	fmt.Println(Bind(source, &s))
	if s != nil {
		fmt.Println(s.CleanedData)
		fmt.Println(s.Name, s.Age, s.Birthday, s.Number, s.Class.ClassName)
	}
}
