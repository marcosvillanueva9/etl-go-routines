package util

import (
	"fmt"
	"strings"
)

func Trim(line []string, dest map[string]string, params map[string]interface{}) {
	fmt.Println("trim function started")
	
	// Todo
}

func Parse(line []string, dest map[string]string, params map[string]interface{}) {
	fmt.Println("parse function started")
	
	// Todo
}

func Concat(line []string, dest map[string]string, params map[string]interface{}, columnMapper map[int]string) {
	fmt.Println("concat function started")


	paramColumns := strings.Split(params["columns"].(string), ",")
	
	for _, column := range paramColumns {
		for i, columnname := range columnMapper {
			if columnname == column {
				dest[params["destination"].(string)] += line[i]
			}
		}
	}
	
}