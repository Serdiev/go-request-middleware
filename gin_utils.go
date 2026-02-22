package gin_utils

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidatePath[T any, T1 any](next func(*gin.Context, T, *T1), svc *T1) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query T

		_ = c.ShouldBindJSON(&query)

		v := reflect.ValueOf(&query).Elem()
		t := v.Type()

		for _, paramName := range c.Params {
			field := findFieldByName(t, paramName.Key)
			if field != nil && v.FieldByIndex(field.Index).CanSet() {
				val := c.Param(paramName.Key)
				fv := v.FieldByIndex(field.Index)
				if fv.Kind() == reflect.String {
					fv.SetString(val)
				}
			}
		}

		next(c, query, svc)
	}
}

func findFieldByName(t reflect.Type, name string) *reflect.StructField {
	nameLower := strings.ToLower(name)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if strings.ToLower(f.Name) == nameLower {
			return &f
		}
	}
	return nil
}

func BindRequest[T any, T1 any](next func(*gin.Context, T, *T1), svc *T1) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request T

		v := reflect.ValueOf(&request).Elem()
		t := v.Type()

		var dataField reflect.Value
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == "Data" && field.Type.Kind() == reflect.Map {
				dataField = v.Field(i)
				break
			}
		}

		if dataField.IsValid() {
			var data map[string]any
			if err := c.ShouldBindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			dataField.Set(reflect.ValueOf(data))
		} else {
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		if field := v.FieldByName("Id"); field.IsValid() && field.CanSet() {
			if id := c.Param("id"); id != "" {
				if field.Kind() == reflect.String {
					field.SetString(id)
				}
			}
		}

		next(c, request, svc)
	}
}
