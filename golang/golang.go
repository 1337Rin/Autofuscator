package golang

import (
	"fmt"
	"strings"
	"io/ioutil"
	"math/rand"
	"time"
)


func Slicecontains(slice []string, keyword string) bool {
    for _, test := range slice {
        if strings.Contains(test, keyword) {
            return true
            break
        }
    }
    return false
}

func Cleanlist(list []string) []string {
	var newlist = []string{}
	for i:=0;i<len(list);i++ {
		if strings.Contains(list[i], "(") || strings.Contains(list[i], ")") {
			continue
		} else if Slicecontains(newlist, list[i]) == false {
			newlist = append(newlist, list[i])
		}
	}

	return newlist
}

func Findvariables(gofile []string) []string {
	var variables = []string{}
	for line:=0;line<len(gofile);line++ {
		if strings.Contains(gofile[line], " :=") {
			brokenLine := strings.Split(gofile[line], " ")

			for i:=1;i<len(brokenLine);i++ {
				if brokenLine[i] == ":="  && brokenLine[i-1] != "_" && brokenLine[i-1] != "err" {
					variables = append(variables, brokenLine[i-1])
				} else if brokenLine[i-1] == "_" || brokenLine[i-1] == "err" {
					continue
				}
			}

		} else if strings.Contains(gofile[line], "var") {
			brokenLine := strings.Split(gofile[line], " ")

			for i:=0;i<len(brokenLine);i++ {
				if brokenLine[i] == "var" {
					variables = append(variables, brokenLine[i+1])
				}
			}
		}
	}
	return Cleanlist(variables)
}

func Findimports(gofile []string) []string {
	var imports = []string{}
	var open bool
	for line:=0;line<len(gofile);line++ {

		if strings.Contains(gofile[line], "import (") == true {
			open = true
			continue
		}
		if strings.Contains(gofile[line], ")") == true {
			open = false
			break
		}
		if open == true {
		imports = append(imports, strings.Replace(gofile[line], `"`, ``, -1))
		}
	}
	return Cleanlist(imports)
}

func FindMethods(gofile []string, imports []string) []string {
	var methods []string
	var method string
	for line:=0;line<len(gofile);line++ {
		for imported:=0;imported<len(imports);imported++ {
			if strings.Contains(gofile[line], fmt.Sprintf(imports[imported]+".")) {
				
				index := strings.Index(gofile[line], fmt.Sprintf(imports[imported]+"."))

				for i:=index;i<len(gofile[line]);i++ {
					if i > 1 {
						if fmt.Sprintf(gofile[line][i-1:i]) == "(" {
							break
						}
					}
					method = fmt.Sprintf(gofile[line][index:i])
				}
				methods = append(methods, method)
				
			}
		}
	}
	return Cleanlist(methods)
}

func Findstrings(gofile []string) []string {
	var fileStrings = []string{}
	var tmpString string
	var open bool = false
	for line:=0;line<len(gofile);line++ {
			if strings.Contains(gofile[line], `"`) {
				brokenLine := strings.Split(gofile[line], "")
				for i:=0;i<len(brokenLine);i++ {
					if brokenLine[i] == `"` && open == false {
						open = true
						//tmpString = fmt.Sprintf(tmpString + brokenLine[i])
					} else if brokenLine[i] != `"` && open == true {
						tmpString = fmt.Sprintf(tmpString + brokenLine[i])
					} else if brokenLine[i] == `"` && open == true {
						open = false
						//tmpString = fmt.Sprintf(tmpString + brokenLine[i])
						fileStrings = append(fileStrings, tmpString)
						tmpString = ""
					}
				}
			}
		}
	return Cleanlist(fileStrings)
}

func Preobfuscation(gofile string) []string {
	content, _ := ioutil.ReadFile(gofile)
	contentProcessing := fmt.Sprintf("%s", content)
	lines := strings.Split(contentProcessing, "\n")
	for i:=0;i<len(lines);i++ {
		lines[i] = strings.TrimSpace(lines[i])
	}
	for j:=0;j<len(lines);j++ {
		if lines[j] == "" {
			lines = append(lines[:j], lines[j+1:]...)
		}
	}
	return lines
}

func RandomName() string {
	var characters = []string{"Q", "O", "I", "l"}
	var randstr string
	for i:=0;i<20;i++ {
		rand.Seed(time.Now().UnixNano())
		v := rand.Intn(4-0) + 0
		randstr += characters[v]
	}
	return randstr
}

func ReplaceVariables(gofile []string, variables []string ) []string {
	var gofileComplete string
	var brokenLine []string  
	for line:=0; line<len(gofile);line++ {
		gofileComplete = fmt.Sprintf(gofileComplete + gofile[line] + "\n")
	}
	for variable:=0; variable<len(variables);variable++ {
		gofileComplete = strings.Replace(gofileComplete, variables[variable], RandomName(), -1)
	}
	brokenLine = strings.Split(gofileComplete, "\n")
	return brokenLine
}

func ReplaceStrings(gofile []string, stringslst []string, importlst []string) []string {
	var gofileComplete string
	var brokenLine []string
	for line:=0; line<len(gofile);line++ {
		gofileComplete = fmt.Sprintf(gofileComplete + gofile[line] + "\n")
	}
	for i:=0; i<len(stringslst);i++ {
		if stringslst[i] == " " {
			continue
		}
		if Slicecontains(importlst, stringslst[i]) {
			continue
		}
		if stringslst[i] == "\\r" || stringslst[i] == "\\n" || stringslst[i] == "\\t" {
			continue
		}
		gofileComplete = strings.Replace(gofileComplete, fmt.Sprintf(`"`+stringslst[i]+`"`), ObfuscateString(stringslst[i]), -1)
	}
	brokenLine = strings.Split(gofileComplete, "\n")
	return brokenLine
}

func ReplaceMethods(gofile []string, methodlst []string) []string {
	var randomname string
	var currentmethod string
	for line:=1;line<len(gofile);line++ {
		if strings.Contains(gofile[line-1], ")") {
			for i:=0;i<len(methodlst);i++ {
				
				randomname = RandomName()
				currentmethod = methodlst[i]
				gofile = append(gofile[:line+1], gofile[line:]...)
				gofile[line] = fmt.Sprintf("var " + randomname + " = " + methodlst[i])

				for j:=line+1;j<len(gofile);j++ {
					if strings.Contains(gofile[j], currentmethod) {
						gofile[j] = strings.Replace(gofile[j], methodlst[i], randomname, -1)
					}
				}
			}
			break
		}
	}
	return gofile
}

//string shit
func init() {
	rand.Seed(time.Now().UnixNano())
}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ObfuscateString(plainStr string) string {
	var obfStr []string
	var obfuscatedString string
	obfStr = append(obfStr,plainStr)


	for _,val := range obfStr{
		b := Str2byte(val)
		idx := GenerateRangeNum(0,6)

		switch idx{
		case 0:
			obfuscatedString = b.tostring1()
		case 1:
			obfuscatedString = b.tostring2(false)
		case 2:
			obfuscatedString = b.tostring3(false)
		case 3:
			obfuscatedString = b.tostring4(false)
		case 4:
			obfuscatedString = b.tostring3(true)
		case 5:
			obfuscatedString = b.tostring4(true)
		case 6:
			obfuscatedString = b.tostring2(true)
		}
	}
	return obfuscatedString

}

type strByte struct{
	rawstr string
	str []string
	bytes []byte
}

func Str2byte(str0 string)*strByte{
	var s0 []string
	b0 := []byte(str0)
	for _,v := range b0{
		s0 = append(s0,string(v))
	}
	return &strByte{rawstr:str0, str:s0, bytes: b0}
}

func (sb *strByte)tostring1() string {
	buildStr := fmt.Sprintf("string([]byte{")
	for i := 0;i < len(sb.str)-1;i++{
		buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("'%s',",sb.str[i]))
	}
	buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("'%s'",sb.str[len(sb.str)-1]) + "})")
	return buildStr
}

func (sb *strByte)tostring2(offset bool) string {
	var buildStr string
	if offset == false{

		buildStr = fmt.Sprintf("string([]byte{") 
		for i := 0;i < len(sb.bytes)-1;i++{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x, ",sb.bytes[i]))
		}
		buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x",sb.bytes[len(sb.bytes)-1]) + "})")
	}else{
		buildStr = fmt.Sprintf("string([]byte{")
		for i := 0;i < len(sb.bytes)-1;i++{
			rand0 := GenerateRangeNum(-16,16)
			if rand0 >0{
				buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x-%d,",sb.bytes[i]+byte(rand0),rand0))
			}else if rand0 <0{
				buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x+%d,",sb.bytes[i]+byte(rand0),-rand0))
			}else{
				buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x,",sb.bytes[i],rand0))
			}

		}
		rand0 := GenerateRangeNum(-16,16)
		if rand0 >0{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x-%d",sb.bytes[len(sb.bytes)-1]+byte(rand0),rand0))
		}else if rand0 <0{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x+%d",sb.bytes[len(sb.bytes)-1]+byte(rand0),-rand0))
		}else{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf("0x%x",sb.bytes[len(sb.bytes)-1],rand0))
		}
		buildStr = fmt.Sprintf(buildStr + "})")
	}
	return buildStr
}

func (sb *strByte)tostring3(offset bool) string {
	var buildStr string
	/* broken
	if offset == false {
		buildStr = fmt.Sprintf("4string(append([]byte{}")
		for i := 0; i < len(sb.bytes); i++ {
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", byte(0x%x)", sb.bytes[i]) + "))")
		}

	}else {
*/
	buildStr = fmt.Sprintf("string(append([]byte{}")
	for i := 0; i < len(sb.bytes); i++ {
		rand0 := GenerateRangeNum(-16,16)
		if rand0 == 0{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", byte(0x%x)", sb.bytes[i]))
		}else if rand0 < 0{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", byte(0x%x+%d)", sb.bytes[i]+byte(rand0),-rand0))
		}else{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", byte(0x%x-%d)", sb.bytes[i]+byte(rand0),rand0))
		}
	}
	buildStr = fmt.Sprintf(buildStr + "))")
	//}
	return buildStr
}


func (sb *strByte)tostring4(offset bool) string {
	var buildStr string
	if offset == false{
		buildStr = fmt.Sprintf("string(append([]byte{}")
		for i := 0;i < len(sb.bytes);i++{
			buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", []byte{0x%x}[0]",sb.bytes[i]))
		}
		buildStr = fmt.Sprintf(buildStr + "))")
	} else{
		buildStr = fmt.Sprintf("string(append([]byte{}")
		for i := 0;i < len(sb.bytes);i++{
			rand0 := GenerateRangeNum(-16,16)
			if rand0 == 0{
				buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", []byte{0x%x}[0]",sb.bytes[i]))
			}else if rand0 < 0{
				buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", []byte{0x%x+%d}[0]",sb.bytes[i]+byte(rand0),-rand0))
			}else{
				buildStr = fmt.Sprintf(buildStr + fmt.Sprintf(", []byte{0x%x-%d}[0]",sb.bytes[i]+byte(rand0),rand0))
			}
		}
		buildStr = fmt.Sprintf(buildStr + "))")
	}
	return buildStr
}

func GenerateRangeNum(min, max int) int {
	randNum := rand.Intn(max - min) + min
	if randNum == 0{
		randNum++
	}
	return randNum
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
