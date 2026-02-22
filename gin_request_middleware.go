package grm

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidatePath[T any](next func(*gin.Context, T)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query T

		_ = c.ShouldBindJSON(&query)
		_ = c.ShouldBindQuery(&query)

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

		next(c, query)
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

func ValidateRequest[T any](next func(*gin.Context, T)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query T

		if err := c.ShouldBindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		v := reflect.ValueOf(&query).Elem()
		if field := v.FieldByName("Id"); field.IsValid() && field.CanSet() {
			if id := c.Param("id"); id != "" {
				if field.Kind() == reflect.String {
					field.SetString(id)
				}
			}
		}

		next(c, query)
	}
}

func ValidateQuery[T any](next func(*gin.Context, T)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query T

		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		v := reflect.ValueOf(&query).Elem()
		if field := v.FieldByName("Id"); field.IsValid() && field.CanSet() {
			if id := c.Param("id"); id != "" {
				if field.Kind() == reflect.String {
					field.SetString(id)
				}
			}
		}

		next(c, query)
	}
}
