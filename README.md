**gojson**

![](./gojson.png)

*gojson* is json parser in go language.

# 1 Install

To start using `gojson`, please install Go and run `go get`

```sh
go get -u github/gaufung/gojson
```

# 2 Usage

After installing  `gojson`, you can use it as follow.

```go
import (
    "io"
    "github/gaufung/gojson"
    "fmt"
)
var json = `{	"debugger" : false,
	      	    "companies" : ["google", "腾讯"],
	      	    "tools" : null,
	       	    "version" : 1.9,
	       	    "scores" : [[1.5,1.7],[[1],[10]]],
	            "persons" : [  {
	       		                    "name" : "C#",
	       		                    "address" : "micorsoft"
	       		               },
	       		               {
	       		                    "name" : "golang",
	       		                    "address" : "american",
	       		                    "scores" : [
                                                    {
	       		                                        "math" : 98.3,
	       		                                        "computer" : 100
                                                    }
                                                ]
                                }
                            ]
	        }`
func main(){
    tokenReader := gojson.NewTokenReaderFromString(json)
    if obj, err := gojson.Parse(tokenReader); err!=nil {
		maps := obj.(map[string]interface{})
        fmt.Println(maps["debugger"]) // false
        fmt.Println(maps["companies"]) // [google 腾讯]
        fmt.Println(maps["tools"]) // null
        fmt.Println(maps["version"]) // 1.9
        fmt.Println(maps["scores"]) // [[1.5,1.7],[[1],[10]]]
        persons := maps["persons"].([]interface{})
        //...
	}else{
		fmt.Println("error")
	}
}
```


# 3 RoadMap

- [x] Fix nesting json parsing.
- [ ] Add serialize features.
- [ ] Appeal to stack rather than recursion to achive high performance.
- [ ] Benchmark testing
