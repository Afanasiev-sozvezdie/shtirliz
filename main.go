package main

import (
	"cblock"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	ACT_ENCODE = 1
	ACT_DECODE = 2
)

type Setting struct {
	Act                  byte
	InFile, OutFile, Key string
}

func main() {
	var sets Setting
	cmds := os.Args[1:]
	i := 0
	clen := len(cmds)

	//Parse comand line params
	for i < clen {
		switch cmds[i] {
		case "-p":
			sets.Key = cmds[i+1]
			i++
		case "-f":
			sets.InFile = cmds[i+1]
			i++
		case "-o":
			sets.OutFile = cmds[i+1]
			i++
		case "-e":
			sets.Act = ACT_ENCODE
		case "-d":
			sets.Act = ACT_DECODE
		}
		i++
	}

	//check settings
	if sets.InFile == "" {
		fmt.Println("Enter File:")
	}

	if sets.Key == "" {
		fmt.Println("Enter Passprase:")
	}

	if sets.OutFile == "" {
		sets.OutFile = sets.InFile + ".coded"
	}

	if sets.Act == 0 {
		fmt.Println("No action selected. File", sets.InFile, "will be encoded to file", sets.OutFile)
		sets.Act = ACT_ENCODE
	}

	Run(sets)
}

func Run(sets Setting) {
	switch sets.Act {
	case ACT_DECODE:
		DecodeFile(sets.InFile, sets.OutFile, sets.Key)
	case ACT_ENCODE:
		EncodeFile(sets.InFile, sets.OutFile, sets.Key)
	}
}

func EncodeFile(filename, codename, passphrase string) {
	buf, e := ioutil.ReadFile(filename)
	if e == nil {
		shifr := cblock.MakeShifr(passphrase)
		code := shifr.Encode(buf)

		f, fe := os.Create(codename)
		if fe == nil {
			fmt.Println("wtite to file", codename)
			f.Write(code)
			f.Close()
		}
	}
}

func DecodeFile(filename, codename, passphrase string) {
	buf, e := ioutil.ReadFile(filename)
	if e == nil {
		shifr := cblock.MakeShifr(passphrase)
		msg := shifr.Decode(buf)

		f, fe := os.Create(codename)
		if fe == nil {
			fmt.Println("wtite to file", codename)
			f.Write(msg)
			f.Close()
		}
	}
}
