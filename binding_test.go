package binding

import (
	"fmt"
	"testing"
	"github.com/smartwalle/going/time"
)

type MyString string

////////////////////////////////////////////////////////////////////////////////////////////////////
type Human struct {
	Name MyString `model:"name"`
	Age  int      `model:"age"`
	BD   time.Timestamp `model:bd`
}

//func (this *Human) CleanedName(n int) (MyString, error) {
//	if n != 0 {
//		return MyString(fmt.Sprintf("%d", n)), nil
//	}
//	return "", errors.New("随便给点吧")
//}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Class struct {
	ClassName string `model:"class_name"`
}

func (this *Class) DefaultClassName() string {
	return "haha"
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Student struct {
	Human
	Number int `model:"number"`
	Class  Class
}

var source = map[string]interface{}{"name": "name field", "age": "123", "number": 1234, "class_name1": "adfsf", "bd": 1444444444}

func TestBindPoint(t *testing.T) {
	var s *Student
	fmt.Println(Bind(source, &s))
	if s != nil {
		fmt.Println(s.Name, s.Age, s.Number, s.Class.ClassName, s.BD.Time())
	}
}
