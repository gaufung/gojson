package gojson

import "testing"

func TestJsonParser(t *testing.T) {
	var json = `{}`
	p := NewJsonParserFromString(json)
	if obj, err := p.Parse(); err != nil {
		t.Error("Parse() error")
	} else {
		maps := obj.(map[string]interface{})
		if len(maps) != 0 {
			t.Error("Parse result failed")
		}
	}

	json = `{"name" : "golang", "version":1.9, "companies":["google", "腾讯"],
		     "verified": true, "user":null}`
	p = NewJsonParserFromString(json)
	if obj, err:= p.Parse(); err!=nil{
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
