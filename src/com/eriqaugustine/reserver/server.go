package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
    "fmt"
    "io/ioutil"
    //"strings"
)


var templ = template.Must(template.New("qr").Parse(templateStr))
var temp2 = template.Must(template.New("qr").Parse(templateStr2))

func main() {
    var port *string = flag.String("port", "3030", "service port");
    flag.Parse()

    http.Handle("/", http.HandlerFunc(QR))
    err := http.ListenAndServe(":" + *port, nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func QR(w http.ResponseWriter, req *http.Request) {
    // templ.Execute(w, req.FormValue("s"))
    response, err := http.Get(req.FormValue("s"));

   //TEST
   println("^^^^");
   println(req.URL.String());
   println(req.FormValue("s"));
   println("vvvv");

    if err != nil {
        fmt.Printf("%s\n", err)
        //os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s\n", err)
            //os.Exit(1)
        }
        //fmt.Printf("%s\n", string(contents))
        //templ.Execute(w, string(contents));
        //w.Write([]byte(string(contents)));
        //println(string(contents));
        //templ.Execute(w, template.HTML(string(contents)));
        //temp2.Execute(w, template.HTML(strings.Replace(string(contents), "'", "\\'", -1)));
        //temp2.Execute(w, strings.Replace(string(contents), "'", "\\'", -1));
        temp2.Execute(w, template.HTML(string(contents)));
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
<iframe id=targetContent seamless style="height: 100%; width: 100%;"></iframe>
<script type="text/javascript">
var doc = document.getElementById('targetContent').contentWindow.document;
doc.open();
doc.write('{{.}}');
doc.close();
</script>
{{else}}
<p>No Target</p>
{{end}}
</body>
</html>
`

/*
const templateStr2 = `
<html>
<head>
<title>Reserve</title>
</head>
<body>
{{if .}}
<iframe seamless style="height: 100%; width: 100%;">
{{.}}
</iframe>
{{end}}
</body>
</html>
`
*/

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
