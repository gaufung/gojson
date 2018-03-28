package gojson

import (
	"testing"
)

func TestJsonStream0(t *testing.T) {
	var json = `{}`
	reader := NewTokenReaderFromString(json)
	if obj, err := Parse(reader); err != nil {
		t.Error("Parse() error")
	} else {
		maps := obj.(map[string]interface{})
		if len(maps) != 0 {
			t.Error("Parse result failed")
		}
	}
}

func TestJsonStream1(t *testing.T) {
	var json = `{"name" : "golang", "version":1.9, "companies":["google", "腾讯"],
		     "verified": true, "user":null}`
	reader := NewTokenReaderFromString(json)
	if obj, err:= Parse(reader); err!=nil{
		t.Error(err.Error())
	}else{
		maps := obj.(map[string]interface{})
		if name, ok := maps["name"]; !ok{
			t.Error("not found name key")
		}else{
			if name.(string) != "golang" {
				t.Error("name value failed")
			}
		}
		if version, ok := maps["version"]; !ok{
			t.Error("not found version key")
		}else{
			if version.(float64) != 1.9 {
				t.Error("version value failed")
			}
		}
		if companies, ok := maps["companies"]; !ok{
			t.Error("not found companies key")
		}else{
			companiesNames := companies.([]interface{})
			if companiesNames[0].(string) != "google" {
				t.Error("companies failed")
			}
			if companiesNames[1].(string) != "腾讯" {
				t.Error("companies failed")
			}
		}
		if verified, ok := maps["verified"];!ok{
			t.Error("Not found verified")
		}else{
			if verified!=true{
				t.Error("verified failed")
			}

		}
		if user,ok := maps["user"]; !ok{
			t.Error("Not found user")
		}else{
			if user!=nil{
				t.Error("user failed")
			}
		}

	}
}


func TestJsonStream2(t *testing.T) {
	json := `{
		"student":
				{
					"name" : "Peter",
					"male" : true ,
					"scores" : [10.0, -4]
				}
		}`
	reader := NewTokenReaderFromString(json)
	if obj, err := Parse(reader); err != nil {
		t.Error("Parse() failed")
	} else {
		maps := obj.(map[string]interface{})
		if student, ok := maps["student"]; !ok{
			t.Error("not found student key")
		}else{
			studentsMap := student.(map[string]interface{})
			if name, ok := studentsMap["name"]; !ok {
				t.Error("not found name key")
			}else{
				if name != "Peter"{
					t.Error("Peter doesn't equal")
				}
			}
			if male, ok := studentsMap["male"]; !ok {
				t.Error("not found male key")
			}else{
				if male != true{
					t.Error("male doesn't equal")
				}
			}
			if scores, ok := studentsMap["scores"]; !ok {
				t.Error("not found socres key")
			}else{
				scoresArray := scores.([]interface{})
				if scoresArray[0] != 10.0 {
					t.Errorf("Actual value is %f", scoresArray[0])
				}
				if scoresArray[1] != -4.0{
					t.Errorf("Actual value is %f", scoresArray[1])
				}
			}

		}
	}
}

func TestJsonStream3(t *testing.T){
	var json = `
		{ "students" : [
					{
					"name" : "Peter"
				},
					{
					"name" : "Mary"
				}
			]
		}
	`
	reader := NewTokenReaderFromString(json)
	if obj, err := Parse(reader); err != nil {
		t.Error("Parse() students failed")
	}else{
		students := obj.(map[string]interface{})["students"].([]interface{})
		student1 := students[0].(map[string]interface{})
		if name, ok := student1["name"]; !ok {
			t.Error("not found name 1 key")
		}else{
			if name != "Peter" {
				t.Error("Peter doesn't equal")
			}
		}
	}
}

func TestJsonStream4(t *testing.T) {
	var json = `
	{
		"languages" : "golang",
		"version" : 1.9,
		"opensource" : true,
		"owener" : null,
		"companies" : [
				 {
				 	"name" : "google",
				 	"creator" : true,
				 	"versions" : [1.5, 1.9, 2.0]
				 	},
				 {
				 	"name" : "腾讯",
				 	"creator" : false,
				 	"versions" : [0.1e-1, -1.3]
				 }
			       ]
	}
	`
	reader := NewTokenReaderFromString(json)
	if obj, err := Parse(reader); err != nil {
		t.Error("Parse() stream4 failed: " + err.Error())
	}else{
		maps := obj.(map[string]interface{})
		if language, ok := maps["languages"]; !ok{
			t.Error("languages failed")
		}else{
			if language!="golang"{
				t.Error("golang doesn't equal")
			}
		}
		if version, ok := maps["version"]; !ok {
			t.Error("version key failed")
		}else{
			if version != 1.9 {
				t.Error("1.9 desn't equal")
			}
		}
		if opensource, ok := maps["opensource"]; !ok{
			t.Error("opensource key failed")
		}else{
			if opensource != true{
				t.Error("true doesn't equal")
			}
		}
		if owner, ok := maps["owener"]; !ok {
			t.Error("owener key failed")
		}else{
			if owner != nil {
				t.Error("nil doesn't equal")
			}
		}
		if companies, ok := maps["companies"]; !ok {
			t.Error("companies key failed")
		}else{
			companiesArray := companies.([]interface{})
			if len(companiesArray) != 2 {
				t.Error("companies count failed")
			}
 		}
	}
}