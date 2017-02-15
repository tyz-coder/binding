package main

import (
	"reflect"
	"fmt"
	"github.com/smartwalle/binding"
	"github.com/gin-gonic/gin"
)

func main() {
	var s = gin.Default()

	s.GET("/", BindingForm(IDForm{}), GinHandlerFuncWrapper(func(c *gin.Context, form *IDForm) {
		fmt.Println(form.Id)
	}))

	s.Run(":9009")
}

type Handler interface {}

type IDForm struct {
	Id string `binding:"id"`
}

func BindingForm(form interface{}) gin.HandlerFunc {
	var formType = reflect.TypeOf(form)
	if formType.Kind() == reflect.Ptr {
		formType = formType.Elem()
	}

	return func(c *gin.Context) {
		var newForm = reflect.New(formType)

		var err = binding.BindWithTag(map[string]interface{}{"id": "aaaa"}, newForm.Interface(), "binding")
		if err != nil {
			fmt.Println("err", err)
			c.AbortWithStatus(404)
			return
		}

		c.Set("binding_form", newForm.Interface())
		c.Next()
	}
}

func GinHandlerFuncWrapper(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var funValue = reflect.ValueOf(h)
		if funValue.IsValid() {
			funValue.Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(c.MustGet("binding_form"))})
		}
		c.Next()
	}
}
