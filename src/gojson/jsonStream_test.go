package gojson

//import (
//	"testing"
//)
//
//func TestJsonStream0(t *testing.T) {
//	var json = `{}`
//	stream := NewJsonStreamFromString(json)
//	if obj, ok := stream.Parse(); !ok {
//		t.Error("Parse() error")
//	} else {
//		maps := obj.(map[string]interface{})
//		if len(maps) != 0 {
//			t.Error("Parse result failed")
//		}
//	}
//}
//
//func TestJsonStream1(t *testing.T) {
//	//var json =`{
//	//		"companies" :[{"name":"google"}, {"name":"facebook"}]
//	//	   }`
//	var json = `{
//		"companies":[{"name":"google"}]
//	}`
//	stream := NewJsonStreamFromString(json)
//	if obj, ok := stream.Parse(); !ok {
//		t.Error("Parser() error")
//	} else {
//		maps := obj.(map[string]interface{})
//		companies := maps["companies"].([]interface{})
//		if len(companies) != 1 {
//			t.Error("companies key error")
//		}
//	}
//}
//
//func TestJsonStream2(t *testing.T) {
//	var json = `{	"debugger" : false,
//	      	"company" : ["google", "apple"],
//	      	"tools" : null,
//	       	"version" : 1.9,
//	       	"scores" : [1.5,1.7,9],
//	        "persons" :  {
//	       		"name" : "golang",
//	       		"address" : "american",
//	       		"isverified": true
//	       		}
//	       }`
//	stream := NewJsonStreamFromString(json)
//	if obj, ok := stream.Parse(); !ok {
//		t.Error("Parse() failed")
//	} else {
//		maps := obj.(map[string]interface{})
//		if maps["debugger"].(bool) != false {
//			t.Error("debugger failed")
//		}
//		if maps["company"].([]interface{})[0].(string) != "google" {
//			t.Error("company failed")
//		}
//		if maps["company"].([]interface{})[1].(string) != "apple" {
//			t.Error("company failed")
//		}
//		if maps["version"].(float64) != 1.9 {
//			t.Error("version failed")
//		}
//		scores := maps["scores"].([]interface{})
//		if scores[0] != 1.5 || scores[1] != 1.7 || scores[2].(float64) != 9 {
//
//			t.Error("scores failed")
//		}
//		persons := maps["persons"].(map[string]interface{})
//		if persons["name"] != "golang" || persons["address"] != "american" ||
//			persons["isverified"] != true {
//			t.Error("persons failed")
//		}
//	}
//
//}
