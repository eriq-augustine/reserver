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
   "regexp"
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

   // TODO(eriq): Does not handle user info.
   targetUrl, err := url.Parse(target);

   if (err != nil) {
      fmt.Printf("Url parse: %s\n", err)
   }

   var urlBase = fmt.Sprintf("%s://%s", targetUrl.Scheme, targetUrl.Host);

   switch request.FormValue("type") {
      case "main":
         var contents *string = getModifiedMain(target, urlBase);
         if (contents != nil) {
            reserve.BasePageTemplate.Execute(response, template.HTML(*contents));
            return;
         }
      case "image":
         var contents *[]byte = getResource(target);
         if (contents != nil) {
            response.Write(*contents);
            return;
         }
      case "js":
         println("TODO: js");
      case "css":
         println("TODO: css");
         var contents *[]byte = getResource(target);
         if (contents != nil) {
            response.Write(*contents);
            return;
         }
      default:
         println("TODO: Default");
         println(request.URL.String());
   }

   // Fall through to 404.
   http.NotFound(response, request);
}

func getResource(target string) *[]byte {
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
               // TODO(eriq).
            case "img":
               var link *string = getAttr(&node.Attr, "src");
               var newLink string = fixLink(*link, "image", urlBase);
               replaceAttr(&node.Attr, "src", newLink);
            case "style":
               // inline CSS
               fmt.Println("CSS 1: ", node.FirstChild.Data);
               node.FirstChild.Data = fixCSS(node.FirstChild.Data, urlBase);
               fmt.Println("CSS 2: ", node.FirstChild.Data);
            case "link":
               // CSS, favicon?
               var link *string = getAttr(&node.Attr, "href");
               var newLink = identifyAndFixLink(*link, urlBase);
               replaceAttr(&node.Attr, "href", newLink);
         }
      }
   });

   var rtn = tree.String();
   return &rtn;
}

func identifyAndFixLink(link string, urlBase string) string {
   if (strings.HasSuffix(link, ".png") ||
       strings.HasSuffix(link, ".jpg") ||
       strings.HasSuffix(link, ".jpeg") ||
       strings.HasSuffix(link, ".ico") ||
       strings.HasSuffix(link, ".gif")) {
      return fixLink(link, "image", urlBase);
   } else if (strings.HasSuffix(link, ".css")) {
      return fixLink(link, "css", urlBase);
   } else if (strings.HasSuffix(link, ".js")) {
      return fixLink(link, "js", urlBase);
   } else {
      return fixLink(link, "main", urlBase);
   }
}

func fixCSS(css string, urlBase string) string {
   re := regexp.MustCompile(`url\s*\(['|"]?(.*)['|"]?\)`);
   return re.ReplaceAllStringFunc(css, func(urlRule string) string {
      var link = re.ReplaceAllString(urlRule, "$1");

      // TODO(eriq): Potential problem is url is unescaped (because of quotes).
      return fmt.Sprintf("url: ('%s')", identifyAndFixLink(link, urlBase));
   });
}

func fixLink(link string, linkType string, urlBase string) string {
   //TODO(eriq): Relative links.
   if (strings.HasPrefix(link, "//")) {
      // Absolute.
      return fmt.Sprintf("/?type=%s&target=http:%s", linkType, link);
   } else if (strings.HasPrefix(link, "/")) {
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
