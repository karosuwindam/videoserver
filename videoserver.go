package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	//	"os"
)

type Filelist struct {
	Name []string
	//Folder []bool
}

type FileTpath struct {
	//cgipath  string `json:"cgipath"`
	Filepath string `json:"filepath"`
	Port     string `json:"port"`
	//playpath string `json:"playpath"`
}

const CGIPATH = "cgi.html"
const PLAYPATH = "play.html"

func Readjison() FileTpath {
	var tmp []FileTpath
	byte, err := ioutil.ReadFile("videoserver.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(byte, &tmp); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(tmp)
	return tmp[0]
}
func (t *Filelist) Read() {
	//filedatapath
	//files, _ := ioutil.ReadDir(FILEPATH)
	var tmp []string
	files, _ := ioutil.ReadDir(filedatapath.Filepath)
	for _, f := range files {
		//t.Name = append(t.Name, f.Name())
		tmp = append(tmp, f.Name())
		//fmt.Println(f.Name())
	}
	t.Name = tmp
	//fmt.Println(t.Name)
}
func cgiEditOutput() string {
	var output string
	for i := 0; i < len(filelist.Name); i++ {
		tmpdata := filelist.Name[i]
		if strings.Index(tmpdata, ".mp4") >= 0 {
			ai := strconv.Itoa(i)
			output += "<div><a href=\"play?id=" + ai + "\">" + tmpdata + "</a></div>\n"
		} else {
			output += "<div>" + tmpdata + "<br><li>" + "</li></div>\n"
		}

	}
	return output
}

func cgiRun(w http.ResponseWriter, r *http.Request) {
	var fp *os.File
	var err error
	var tmp string
	fmt.Println(r)
	filelist.Read()
	fp, err = os.Open(CGIPATH)
	if err != nil {

	}
	reader := bufio.NewReaderSize(fp, 4096)
	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		if strings.Index(line, "<%output%>") >= 0 {
			//line = cgiEditOutput()
			line = strings.Replace(line, "<%output%>", cgiEditOutput(), 1)

		}
		tmp += line
		//fmt.Print(line)
	}
	if err != io.EOF {
		panic(err)
	}
	//tmp := r.URL.RawQuery
	fmt.Fprintf(w, "%s", tmp)
}
func cgiPlay(w http.ResponseWriter, r *http.Request) {
	var fp *os.File
	var err error
	var tmp string
	var n int

	idtmp := strings.Split(r.URL.RawQuery, "=")
	id, _ := strconv.Atoi(idtmp[1])
	fp, err = os.Open(PLAYPATH)
	if err != nil {

	} /*
		reader := bufio.NewReaderSize(fp, 4096)
		for line := ""; err == nil; line, err = reader.ReadString('\n') {
			if strings.Index(line, "<%output%>") >= 0 {
				//line = strings.Repeat(line, "<%output%>", filelist.Name[id], 1)
				//line = cgiEditOutput()
			}
			tmp += line
			//fmt.Print(line)
		}
		if err != io.EOF {
			panic(err)
		}*/
	buf := make([]byte, 1024)
	for {
		n, err = fp.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		tmp += string(buf[:n])
	}
	fmt.Println(filelist.Name[id])
	str1 := url.QueryEscape(filelist.Name[id])
	str1 = regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(str1, "$1%20")
	/*str1 := filelist.Name[id]
	if strings.Index(str1, "#") >= 0 {
		str1 = strings.Replace(str1, "#", "%2388", 1)
	}*/
	if strings.Index(tmp, "<%output%>") >= 0 {
		tmp = strings.Replace(tmp, "<%output%>", str1, 1)
		//line = cgiEditOutput()
	}
	//fmt.Println(tmp)
	//tmp := r.URL.RawQuery
	fmt.Fprintf(w, "%s", tmp)

}

var filelist Filelist
var filedatapath FileTpath

func main() {
	filedatapath = Readjison()
	if filedatapath.Filepath == "" {
		return
	}
	if filedatapath.Port == "" {
		return
	}
	//fmt.Print(filedatapath)
	//tmp := new(Filelist)
	filelist.Read()
	//fmt.Println(filelist.Name)
	http.HandleFunc("/play", cgiPlay)
	http.HandleFunc("/cgi", cgiRun)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(filedatapath.Filepath))))
	http.ListenAndServe(":"+filedatapath.Port, nil)

}
