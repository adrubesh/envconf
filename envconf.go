package envconf

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

var (
	errHandler = DefaultErrHandler()
)

type ErrHandler func(err error)

func DefaultErrHandler() ErrHandler {
	return func(err error) {
		log.Fatalf("envconf: %s\n", err)
	}
}

/* Sets the error handler */
func SetErrHandler(h ErrHandler) {
	errHandler = h
}

/*
	envconf makes it simple to parse environmental variables to a struct

	it utilizes struct tags to initialize default values, overriding them with corresponding
	environmental variables if they are found
*/

// defaultEnv looks up the environmental variable specified by key
// returns defaultValue if it doesn't exist or the key is empty
func defaultEnv(key, defaultValue string) string {
	v, exists := os.LookupEnv(key)
	if !exists || len(v) == 0 {
		return defaultValue
	}
	return v
}

func LoadConfig(c interface{}) {
	v := reflect.ValueOf(c)
	t := v.Elem().Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		var err error
		// get the custom tags, env = environmental variable to parse, default = default value if not found
		envVar := field.Tag.Get("env")
		defaultVal := field.Tag.Get("default")

		required := false
		if required, err = strconv.ParseBool(field.Tag.Get("required")); err != nil {
			required = false
		}

		if len(envVar) > 0 || len(defaultVal) > 0 {
			val := defaultEnv(envVar, defaultVal)

			if required && len(defaultVal) == 0 {
				errHandler(fmt.Errorf("required field %s [%s] has no default value", field.Name, envVar))
			}

			switch field.Type.Kind() {
			case reflect.String:
				v.Elem().Field(i).SetString(val)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intVal, err := strconv.Atoi(val)
				if err != nil {
					errHandler(fmt.Errorf("failed to parse %s [%s] as int: %s", field.Name, envVar, err))
				}

				v.Elem().Field(i).SetInt(int64(intVal))
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(val)
				if err != nil {
					errHandler(fmt.Errorf("failed to parse %s [%s] as bool: %s", field.Name, envVar, err))
				}

				v.Elem().Field(i).SetBool(boolVal)
			}

		}
	}
}
