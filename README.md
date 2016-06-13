## Binding
Binding 是一个利用 Golang 的反射机制，将 map 对象的数据映射到 struct 的工具包，可用于将 HTTP 请求参数映射到指定的 struct。

业界关于将 HTTP 参数绑定到 struct 的工具库虽然已经有很多，但是大多都只是对数据进行简单的映射，缺少灵活的控制，本工具来自于实践，或者也适用于你。

#### 例子
```
import (
	"fmt"
	"testing"
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

func (this *Student) DefaultClass() string {
	return "Class one"
}

var formData = map[string]interface{}{"name": "Yangfeng", "age": "12", "number": "9", "birthday": "2016-06-12"}

func TestSample(t *testing.T) {
	// 绑定
	var s *Student
	var err = BindWithTag(formData, &s, "form")
	if err != nil {
		fmt.Println("绑定失败")
	}
	fmt.Println(s)
}
```

#### 自定义 Tag 名称
上面代码是一个将 map 对象和 struct 进行绑定的简单例子， formData 是我们要绑定的数据源，我们要它的数据绑定到 Student 结构体上，除了 BindWithTag 函数有第三个参数外，和其它用于数据绑定的组件没有什么区别，这也是本组件的灵活点之一，可以任意设定结构体属性和数据源数据对应的关系。

#### 默认值 (Default 方法)
如果您仔细看，会发现结构体 Student 有一个 DefaultClass 的方法，该方法用于返回 Student 结构体的 Class 属性的默认值。对于一些场景，比如 HTTP 请求，可能请求的数据源中并不存在某些参数，但是我们又需要这些参数，所以，你可以通过这种方式来指定某一个属性的默认值。具体的方法为给该结构体添加一个以 Default+属性名 的方法，该方法需要返回一个具体的值。当在被绑定的数据源中找不到对应的关系，就会调用该方法获取一个默认的值。

#### 清理数据 (Cleaned 方法)
对于一个未知的数据源，里面的数据可能是乱七八糟的，对于我们来讲有太多未知的因素，如果在绑定的时候就能够清理一次或者转换为我们需要的数据类型，着实可以减少不少 if 语句。

上面的代码中，结构体 Human 有一个名叫 CleanedBirthday 的方法，就是用于“清理” Birthday 属性的值。数据源中的 birthday 是一个字符串，我们需要将其转换为 time.Time 对象，所以特意添加了一个方法来对数据进行加工。注意方法名称为 Cleaned+属性名，方法有一个参数，为原始数据，其类型需要和原始数据类型保持一致。Cleaned 相关的方法会返回两个数据，一个是“加工”之后的数据，另一个是满足 error 接口的对象。如果返回的 error 不为空，则会在 Bind 方法返回该错误。

#### 特别说明

* 通过 Default 方法获取到的属性值，不会再经过 Cleaned 方法进行加工；
* Cleaned 方法返回的数据类型必须和对应属性的类型一致；
* 如果数据源中的数据类型和结构体的属性数据类型不一致，将把数据源的数据类型转换为结构体对应属性的类型；
* 支持 Slice，详见 test。

## HTTP
关于和 HTTP 参数进行绑定，可以参考 [FORM](https://github.com/smartwalle/form) 组件，该组件是本组件的一个延伸。


