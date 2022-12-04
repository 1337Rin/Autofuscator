package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"flag"
	"github.com/1337Rin/Autofuscator/golang"
	//"github.com/1337Rin/Autofuscator/powershell"
)

func exists(fileName string) bool {
    _, err := os.Stat(fileName)
    if errors.Is(err, os.ErrNotExist) {
        return false
    }
    return true
}

func is_readable(fileName string) bool {
    file, err := os.Open(fileName)
    if err != nil {
        return false
    }
    defer file.Close()
    return true
}

func Slicecontains(slice []string, keyword string) bool {
    for _, test := range slice {
        if strings.Contains(test, keyword) {
            return true
            break
        }
    }
    return false
}

func main() {
	var validLanguage = true
	var validInfile = false

	var outfile string
	var infile string

	var language string

	var stat bool

	var b1 bool
	var b2 bool
	var b3 bool

	flag.StringVar(&language, "l", "", "language")
	flag.BoolVar(&stat, "stat", false, "stat mode")
	flag.StringVar(&outfile, "o", "", "out file")
	flag.StringVar(&infile, "f", "", "in file")

	flag.BoolVar(&b1, "b1", false, "replace variable names")
	flag.BoolVar(&b2, "b2", false, "replace package and method names")
	flag.BoolVar(&b3, "b3", false, "encode strings")

	flag.Parse()


	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	if language == "go" || language == "golang" {
		validLanguage = true
	} else if language == "ps" || language == "powershell" {
		validLanguage = true
	} else {
		fmt.Println(language, "is not a suported language")
		validLanguage = false
		os.Exit(1)
	}

	if is_readable(infile) && exists(infile) {
		validInfile = true
	} else {
		if exists(infile) == false {
			fmt.Println("input file '", infile , "' does not exist")
		} else if is_readable(infile) == false {
			fmt.Println("input file '", infile, "' is not readable")
		}
	}

	if validLanguage && validInfile && stat == false {
		if language == "go" || language == "golang" {
			lines := golang.Preobfuscation(infile)
			variables := golang.Findvariables(lines)
			imports := golang.Findimports(lines)
			strings := golang.Findstrings(lines)
			methods := golang.FindMethods(lines, imports)

			if b1 {
				lines = golang.ReplaceVariables(lines, variables)
			}
			if b2 {
				lines = golang.ReplaceMethods(lines, methods)
			}
			if b3 {
				lines = golang.ReplaceStrings(lines, strings, imports)
			}

			if outfile == "" {
				for i:=0;i<len(lines);i++ {
					fmt.Println(lines[i])
				}
			}
		}

	} else if validLanguage && stat && validInfile {
		if language == "go" || language == "golang" {
			lines := golang.Preobfuscation(infile)
			variables := golang.Findvariables(lines)
			imports := golang.Findimports(lines)
			strings := golang.Findstrings(lines)
			methods := golang.FindMethods(lines, imports)

			fmt.Println("variables:")
			for i:=0;i<len(variables);i++ {
				fmt.Println(`	"` + variables[i] + `"`)
			}

			fmt.Println("\nimports:")
			for i:=0;i<len(imports);i++ {
				fmt.Println("	`" + imports[i] + "`")
			}

			fmt.Println("\nstrings:")
			for i:=0;i<len(strings);i++ {
				fmt.Println(`	"` + strings[i] + `"`)
			}
			fmt.Println("\nmethods:")
			for i:=0;i<len(methods);i++ {
				fmt.Println(`	"` + methods[i] + `"`)
			}

		} else if language == "ps" || language == "powershell" {
		}
	} 
}