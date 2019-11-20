package mux

import (
	"testing"
	"strings"
	"fmt"
)

func TestParseMatch(t *testing.T) {
	pattern:="/db/:key/meng/:value/huang"
	i:=strings.Index(pattern,":")
	prefix:=pattern[:i]
	match := strings.Split(pattern[i:], "/")
	params:=make(map[string]string)
	key:=""
	for i:=0;i<len(match);i++{
		if strings.Contains(match[i],":"){
			match[i]=strings.Trim(match[i],":")
			params[match[i]]=""
			if i>0{
				key+="/"
			}
		}else {
			key+="/"+match[i]
			match[i]=""
		}
	}
	path:="/db/123/meng/456/huang"
	strs := strings.Split(strings.Trim(path,prefix), "/")
	if len(strs)==len(match){
		for i:=0;i<len(strs);i++{
			if match[i]!=""{
				if _,ok:=params[ match[i]];ok{
					params[ match[i]]=strs[i]
				}
			}
		}
	}
	fmt.Println(params)
}