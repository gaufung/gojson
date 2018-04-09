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

type student struct {
	name   string
	age    int
	isMale bool
	scores []float64
}

type class struct {
	name     string
	students []student
}

func TestJsonParser1(t *testing.T) {
	if bytes,err := Encode(nil); err!=nil{
		t.Error(err.Error())
	}else{
		if string(bytes) != "null"{
			t.Error("nil failed")
		}
	}
	if bytes, err:= Encode(10); err!=nil{
		t.Error(err.Error())
	}else{
		if string(bytes) != "10"{
			t.Error("int failed")
		}
	}
	if bytes, err:= Encode(14.2); err!=nil{
		t.Error(err.Error())
	}else{
		if string(bytes) != "14.2"{
			t.Error("float failed")
		}
	}
	if bytes, err:=Encode("123");err!=nil{
		t.Error(err.Error())
	}else{
		if string(bytes) != `"123"`{
			t.Error("string failed")
		}
	}
	if bytes, err:=Encode([]int{1,2,3});err!=nil{
		t.Error(err.Error())
	}else{
		if string(bytes) != `[1,2,3]`{
			t.Error("array failed")
		}
	}
	if bytes, err:=Encode(true);err!=nil{
		t.Error(err.Error())
	}else{
		if string(bytes)!="true"{
			t.Error("bool failed")
		}
	}
	if bytes, err:=Encode(false);err!=nil{
		t.Error(err.Error())
	}else{
		if string(bytes)!="false"{
			t.Error("bool failed")
		}
	}
	maps := make(map[int]string)
	if _, err:=Encode(maps); err==nil{
		t.Error("map success")
	}

	student1 := student{name:"tom",age:10,isMale:true, scores:[]float64{1.2,3.2}}
	if bytes, err:=Encode(student1);err!=nil{
		t.Error(err.Error())
	}else{
		expect := `{"name":"tom","age":10,"isMale":true,"scores":[1.2,3.2]}`
		if string(bytes)!=expect{
			t.Error("struct failed")
		}
	}
	class1 := class{"Grade1",[]student{student1}}
	if bytes, err:=Encode(class1);err!=nil{
		t.Error(err.Error())
	}else{
		expect:=`{"name":"Grade1","students":[{"name":"tom","age":10,"isMale":true,"scores":[1.2,3.2]}]}`
		if string(bytes)!=expect{
			t.Error("struct failed")
		}
	}

}
