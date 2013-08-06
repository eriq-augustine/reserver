package main

import (
   "flag"
   "html/template"
   "log"
   "net/http"
   "net/url"
   "fmt"
   "io/ioutil"
   "io"
   "strings"
   "com/eriqaugustine/reserver/reserve"
   "code.google.com/p/go-html-transform/h5"
   "code.google.com/p/go.net/html"
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
   var target string = request.FormValue("target");
   //TEST
   println("^^^^");
   println(request.URL.String());
   println(target);
   println(request.FormValue("type"));
   println("vvvv");

   // TODO(eriq): Does not handle user info.
   targetUrl, err := url.Parse(target);

   if (err != nil) {
      fmt.Printf("Url parse: %s\n", err)
   }

   var urlBase = fmt.Sprintf("%s://%s", targetUrl.Scheme, targetUrl.Host);

   //TEST
   println("Base: " + urlBase);

   switch request.FormValue("type") {
      case "main":
         println("main");
         var contents *string = getModifiedMain(target, urlBase);
         if (contents != nil) {
            reserve.BasePageTemplate.Execute(response, template.HTML(*contents));
            return;
         }
      case "image":
         println("image");
         var contents *[]byte = getImage(target);
         if (contents != nil) {
            response.Write(*contents);
            return;
         }
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

// TODO(eriq): Make general getResource().
func getImage(target string) *[]byte {
   response, err := http.Get(target);

   if (err != nil) {
      // TODO(eriq): Better error logging.
      fmt.Printf("Image Fetch Error: %s\n", err)
      return nil;
   }

   defer response.Body.Close();
   contents, err := ioutil.ReadAll(response.Body);

   if (err != nil) {
      fmt.Printf("Image Read Error: %s\n", err)
      return nil;
   }

   return &contents;
}

func getModifiedMain(target string, urlBase string) *string {
   response, err := http.Get(target);

   if (err != nil) {
      fmt.Printf("Fetch Error: %s\n", err)
      return nil;
   }

   defer response.Body.Close();

   var rtn *string = replaceLinks(response.Body, urlBase);
   if (rtn != nil) {
      return rtn;
   }
   return nil;

   /*TEST
   contents, err := ioutil.ReadAll(response.Body);

   if (err != nil) {
      fmt.Printf("Read Body Error: %s\n", err)
      return nil;
   }

   var rtn string = string(contents);
   return &rtn;
   */
}

func replaceLinks(responseBody io.Reader, urlBase string) *string {
   tree, err := h5.New(responseBody);

   if (err != nil) {
      return nil;
   }

   tree.Walk(func(node *html.Node) {
      if (node.Type == html.ElementNode) {
         switch node.Data {
            case "a":
            case "img":
               var link *string = getAttr(&node.Attr, "src");
               var newLink string = fixLink(*link, "image", urlBase);
               replaceAttr(&node.Attr, "src", newLink);
         }
      }
   });

   var rtn = tree.String();
   return &rtn;
}

func fixLink(link string, linkType string, urlBase string) string {
   //TODO(eriq): Relative links.
   if (strings.HasPrefix(link, "/")) {
      return fmt.Sprintf("/?type=%s&target=%s%s", linkType, urlBase, link);
   } else {
      return fmt.Sprintf("/?type=%s&target=%s", linkType, link);
   }

   return link;
}

func getAttr(attrs *[]html.Attribute, key string) *string {
   for _, attr := range *attrs {
      if (attr.Key == key) {
         return &attr.Val;
      }
   }

   return nil;
}

func replaceAttr(attrs *[]html.Attribute, key string, value string) {
   for i := 0; i < len(*attrs); i++ {
      if ((*attrs)[i].Key == key) {
         (*attrs)[i].Val = value;
      }
   }
}
