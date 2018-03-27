package gojson

import (
	"testing"
)

func TestJsonStream0(t *testing.T) {
	var json = `{}`
	reader := NewTokenReaderFromString(json)
	if obj, err := Parser(reader); err != nil {
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
	if obj, err:= Parser(reader); err!=nil{
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
	//var json = `{	"debugger" : false,
	//      	"company" : ["google", "apple"],
	//      	"tools" : null,
	//       	"version" : 1.9,
	//       	"scores" : [1.5,1.7,9],
	//        "persons" :  {
	//       		"name" : "golang",
	//       		"address" : "american",
	//       		"isverified": true
	//       		}
	//       }`
	json := `{
		"students": [
				{
					"name" : "Peter",
					"male" : true ,
					"scores" : [10.0, -4]
				},
				{
					"name" : "Mary",
					"male" : false,
					"scores" : [12, 3e4]
				}
			     ]
		}`
	reader := NewTokenReaderFromString(json)
	if obj, err := Parser(reader); err != nil {
		t.Error("Parse() failed")
	} else {
		t.Error(obj)
		//maps := obj.(map[string]interface{})
		//if _, ok := maps["students"]; !ok{
		//	t.Error("not found stduents key")
		//}else{
		//	//studentsArray := students.([]interface{})
		//	//student1 := studentsArray[0].(map[string]interface{})
		//	//if name1, ok1 := student1["name"]; !ok1{
		//	//	t.Error("not found name key")
		//	//}else{
		//	//	if name1 != "Peter"{
		//	//		t.Error("name1 failed")
		//	//	}
		//	//}
		//
		//
		//}
	}

}
