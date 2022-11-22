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

		imports = append(imports, gofile[line])
		}
	}
	return Cleanlist(imports)
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
					} else if brokenLine[i] != `"` && open == true {
						tmpString = fmt.Sprintf(tmpString + brokenLine[i])
					} else if brokenLine[i] == `"` && open == true {
						open = false
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
	var characters = []string{"0", "O", "I", "l"}
	var randstr string
	for i:=0;i<20;i++ {
		rand.Seed(time.Now().UnixNano())
		v := rand.Intn(4-0) + 0
		randstr += characters[v]
	}
	return randstr

}
