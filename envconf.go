package envconf

import (
	"log"
	"os"
	"reflect"
	"strconv"
)

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

		// get the custom tags, env = environmental variable to parse, default = default value if not found
		envVar := field.Tag.Get("env")
		defaultVal := field.Tag.Get("default")

		if len(envVar) > 0 || len(defaultVal) > 0 {
			val := defaultEnv(envVar, defaultVal)

			switch field.Type.Kind() {
			case reflect.String:
				v.Elem().Field(i).SetString(val)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intVal, err := strconv.Atoi(val)
				if err != nil {
					log.Fatalln(err)
				}

				v.Elem().Field(i).SetInt(int64(intVal))
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(val)
				if err != nil {
					log.Fatalln(err)
				}

				v.Elem().Field(i).SetBool(boolVal)
			}

		}
	}
}
