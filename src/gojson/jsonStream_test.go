package gojson

import (
	"testing"
)

var json = `{	"debugger" : false,
	      	"company" : ["google", "apple"],
	      	"tools" : null,
	       	"version" : 1.9,
	       	"scores" : [1.5,1.7,9],
	        "persons" :  {
	       		"name" : "golang",
	       		"address" : "american",
	       		"isverified": true
	       		}
	       }`

func TestJsonStream(t *testing.T) {
	stream := NewJsonStreamFromString(json)
	if obj, ok := stream.Parse(); !ok{
		t.Error("Parse() failed")
	}else{
		maps := obj.(map[string]interface{})
		if maps["debugger"].(bool) != false{
			t.Error("debugger failed")
		}
		if maps["company"].([]interface{})[0].(string) != "google" {
			t.Error("company failed")
		}
		if maps["company"].([]interface{})[1].(string) != "apple" {
			t.Error("company failed")
		}
		if maps["version"].(float64) != 1.9 {
			t.Error("version failed")
		}
		scores := maps["scores"].([]interface{})
		if scores[0]!=1.5 || scores[1] != 1.7 || scores[2].(float64) != 9{

			t.Error("scores failed")
		}
		persons := maps["persons"].(map[string]interface{})
		if persons["name"] != "golang" || persons["address"] != "american" ||
				persons["isverified"] != true{
			t.Error("persons failed")
		}
	}

}