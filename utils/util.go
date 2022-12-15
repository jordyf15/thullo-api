package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"reflect"
	"time"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func DataResponse(data interface{}, metadata interface{}) map[string]interface{} {
	return map[string]interface{}{"data": data, "meta": metadata}
}

func ToSHA256(toHash string) string {
	hash := sha256.New()
	hash.Write([]byte(toHash))
	return hex.EncodeToString(hash.Sum(nil))
}

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func RandFileName(prefix, suffix string) string {
	return prefix + RandString(8) + suffix
}

func ToBSON(object interface{}) bson.M {
	var tagValue string

	data := bson.M{}
	var element reflect.Value
	if reflect.ValueOf(object).Kind() == reflect.Ptr {
		element = reflect.ValueOf(object).Elem()
	} else {
		element = reflect.ValueOf(object)
	}

	for i := 0; i < element.NumField(); i++ {
		typeField := element.Type().Field(i)
		tag := typeField.Tag

		tagValue = tag.Get("bson")

		if tagValue == "-" {
			continue
		}

		switch element.Field(i).Kind() {
		case reflect.Array, reflect.Slice:
			if objectID, ok := element.Field(i).Interface().(primitive.ObjectID); ok {
				data[tagValue] = objectID
				continue
			} else if byteArr, ok := element.Field(i).Interface().([]byte); ok {
				data[tagValue] = byteArr
				continue
			}

			value := bson.A{}
			arr := element.Field(i)
			for j := 0; j < arr.Len(); j++ {
				item := arr.Index(j)

				switch item.Kind() {
				case reflect.Struct:
					switch item.Interface().(type) {
					case time.Time:
						value = append(value, item.Interface())
					default:
						obj := item.Interface()
						value = append(value, ToBSON(&obj))
					}
				default:
					obj := objectFromKind(item)
					if obj == nil {
						continue
					}

					value = append(value, obj)
				}
			}

			data[tagValue] = value
		case reflect.Struct:
			value := element.Field(i).Interface()
			switch value.(type) {
			case time.Time:
				data[tagValue] = value
			default:
				data[tagValue] = ToBSON(value)
			}
		default:
			value := objectFromKind(element.Field(i))
			if value == nil {
				continue
			}

			data[tagValue] = value
		}
	}

	return data
}

func objectFromKind(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return value.Bool()
	case reflect.Int:
		return value.Int()
	case reflect.Int8:
		return int8(value.Int())
	case reflect.Int16:
		return int16(value.Int())
	case reflect.Uint8:
		return uint8(value.Uint())
	default:
		return value.Interface()
	}
}
