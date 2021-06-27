package commands

import (
	"os"
	"log"
	"io/ioutil"
	"strings"
)

func WriteKey(path string, content []byte) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	f.Write(content)
}

func SetFa12AndBifrostContractAddress(bipath, fa12path string) (string, string){
	biAddr, err := ioutil.ReadFile(bipath)
	if err != nil {
		log.Fatalln(err)
	}
	fa12Addr, err := ioutil.ReadFile(fa12path)
	if err != nil {
		log.Fatalln(err)
	}

	c_addr := strings.TrimSpace(string(biAddr))
	fa_addr := strings.TrimSpace(string(fa12Addr))

	return c_addr, fa_addr
}