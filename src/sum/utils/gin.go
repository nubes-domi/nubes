package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func CtxGet[T any](c *gin.Context, key string) (T, bool) {
	iVal, ok := c.Get(key)
	if !ok {
		var blank T
		return blank, false
	}

	val, ok := iVal.(T)
	if !ok {
		log.Panicf("invalid conversion from context")
	}

	return val, true
}

func CtxMustGet[T any](c *gin.Context, key string) T {
	val, ok := CtxGet[T](c, key)
	if !ok {
		log.Panicf("could not Get key %s", key)
	}

	return val
}
