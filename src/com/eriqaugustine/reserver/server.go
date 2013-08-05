package main

import (
   "flag"
   "html/template"
   "log"
   "net/http"
   "fmt"
   "io/ioutil"
   //"strings"
   "com/eriqaugustine/reserver/reserve"
)

func main() {
   var port *string = flag.String("port", "3030", "service port");
   flag.Parse()

   http.Handle("/", http.HandlerFunc(Reserve))
   err := http.ListenAndServe(":" + *port, nil)
   if err != nil {
      log.Fatal("ListenAndServe:", err)
   }
}

func Reserve(response http.ResponseWriter, request *http.Request) {
   //TEST
   println("^^^^");
   println(request.URL.String());
   println(request.FormValue("target"));
   println(request.FormValue("type"));
   println("vvvv");

   switch request.FormValue("type") {
      case "main":
         println("main");
         var contents *string = getModifiedMain(request.FormValue("target"));
         if (contents != nil) {
            reserve.BasePageTemplate.Execute(response, template.HTML(*contents));
         }
         return;
      case "image":
         println("image");
      case "js":
         println("js");
      case "css":
         println("css");
      default:
         println("Default");
   }

   // Fall through to 404.
   http.NotFound(response, request);

   //fmt.Printf("%s\n", string(contents))
   //templ.Execute(w, string(contents));
   //w.Write([]byte(string(contents)));
   //println(string(contents));
   //templ.Execute(w, template.HTML(string(contents)));
   //temp2.Execute(w, template.HTML(strings.Replace(string(contents), "'", "\\'", -1)));
   //temp2.Execute(w, strings.Replace(string(contents), "'", "\\'", -1));
}

func getModifiedMain(target string) *string {
   response, err := http.Get(target);

   if (err != nil) {
      fmt.Printf("Fetch Error: %s\n", err)
      return nil;
   }

   defer response.Body.Close();
   contents, err := ioutil.ReadAll(response.Body);

   if (err != nil) {
      fmt.Printf("Read Body Error: %s\n", err)
      return nil;
   }

   var rtn string = string(contents);
   return &rtn;
}
