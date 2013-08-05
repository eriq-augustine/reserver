package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
    "fmt"
    "io/ioutil"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))
var temp2 = template.Must(template.New("qr").Parse(templateStr2))

func main() {
    flag.Parse()
    http.Handle("/", http.HandlerFunc(QR))
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func QR(w http.ResponseWriter, req *http.Request) {
    // templ.Execute(w, req.FormValue("s"))
    response, err := http.Get(req.FormValue("s"));
    if err != nil {
        fmt.Printf("%s", err)
        //os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            //os.Exit(1)
        }
        //fmt.Printf("%s\n", string(contents))
        //templ.Execute(w, string(contents));
        w.Write([]byte(string(contents)));
        //println(string(contents));
        //temp2.Execute(w, string(contents));
    }
}

const templateStr = `
{{if .}}
{{.}}
{{end}}
`

const templateStr2 = `
<html>
<head>
<title>Reserve</title>
</head>
<body>
{{if .}}
<iframe style="height: 100%; width: 100%;">
{{.}}
</iframe>
{{end}}
</body>
</html>
`

/*
const templateStr = `
<html>
<head>
<title>Reserve</title>
</head>
<body>
{{if .}}
<iframe style="height: 100%; width: 100%;"src="{{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
</body>
</html>
`
*/
