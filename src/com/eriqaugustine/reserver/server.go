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
   "strconv"
   "time"
   "bytes"
)

const (
   REQUEST_TYPE_MAIN int = iota
   REQUEST_TYPE_IMAGE
   REQUEST_TYPE_JS
   REQUEST_TYPE_CSS
   REQUEST_TYPE_UNKNOWN
   NUM_REQUEST_TYPES
);

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

   if (target == "") {
      http.NotFound(response, request);
      return;
   }

   // TODO(eriq): Does not handle user info.
   targetUrl, err := url.Parse(target);

   if (err != nil) {
      fmt.Printf("Url parse: %s\n", err)
   }

   intType, err := strconv.Atoi(request.FormValue("type"));

   if (err != nil) {
      fmt.Println("Type is not an int: ", request.FormValue("type"));
      intType = REQUEST_TYPE_UNKNOWN;
   }

   switch intType {
      case REQUEST_TYPE_MAIN:
         var contents *string = getModifiedMain(target, targetUrl);
         if (contents != nil) {
            reserve.BasePageTemplate.Execute(response, template.HTML(*contents));
            return;
         }
      case REQUEST_TYPE_IMAGE:
         var contents *[]byte = getResource(target);
         if (contents != nil) {
            var modTime time.Time;
            var contentReader *bytes.Reader = bytes.NewReader(*contents);
            http.ServeContent(response, request, targetUrl.Path, modTime,  contentReader);
            return;
         }
      case REQUEST_TYPE_JS:
         var contents *[]byte = getResource(target);
         if (contents != nil) {
            var modTime time.Time;
            var contentReader *bytes.Reader = bytes.NewReader(*contents);
            http.ServeContent(response, request, targetUrl.Path, modTime,  contentReader);
            return;
         }
      case REQUEST_TYPE_CSS:
         var contents *[]byte = getResource(target);
         if (contents != nil) {
            var modTime time.Time;
            var contentReader *bytes.Reader = bytes.NewReader(*contents);
            http.ServeContent(response, request, targetUrl.Path, modTime,  contentReader);
            return;
         }
      default:
         println("TODO: Default");
         fmt.Println("   Target: ", targetUrl);
         fmt.Println("   Url: ", request.URL);
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

func getModifiedMain(target string, targetUrl *url.URL) *string {
   response, err := http.Get(target);

   if (err != nil) {
      fmt.Printf("Fetch Error: %s\n", err)
      return nil;
   }

   defer response.Body.Close();

   var rtn *string = replaceLinks(response.Body, targetUrl);
   if (rtn != nil) {
      return rtn;
   }
   return nil;
}

func replaceLinks(responseBody io.Reader, targetUrl *url.URL) *string {
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
               if (link != nil) {
                  var newLink string = fixLink(*link, REQUEST_TYPE_IMAGE, targetUrl);
                  replaceAttr(&node.Attr, "src", newLink);
               }
            case "style":
               // inline CSS
               node.FirstChild.Data = fixCSS(node.FirstChild.Data, targetUrl);
            case "link":
               // CSS, favicon?
               var link *string = getAttr(&node.Attr, "href");
               if (link != nil) {
                  var newLink = identifyAndFixLink(*link, targetUrl);
                  replaceAttr(&node.Attr, "href", newLink);
               }
            case "script":
               var link *string = getAttr(&node.Attr, "src");
               if (link != nil) {
                  var newLink string = fixLink(*link, REQUEST_TYPE_JS, targetUrl);
                  replaceAttr(&node.Attr, "src", newLink);
               }
         }
      }
   });

   var rtn = tree.String();
   return &rtn;
}

func identifyLink(link string) int {
   if (strings.HasSuffix(link, ".png") ||
       strings.HasSuffix(link, ".jpg") ||
       strings.HasSuffix(link, ".jpeg") ||
       strings.HasSuffix(link, ".ico") ||
       strings.HasSuffix(link, ".gif")) {
      return REQUEST_TYPE_IMAGE;
   } else if (strings.HasSuffix(link, ".css")) {
      return REQUEST_TYPE_CSS;
   } else if (strings.HasSuffix(link, ".js")) {
      return REQUEST_TYPE_JS;
   } else {
      return REQUEST_TYPE_UNKNOWN;
   }
}

func identifyAndFixLink(link string, targetUrl *url.URL) string {
   return fixLink(link, identifyLink(link), targetUrl);
}

func fixCSS(css string, targetUrl *url.URL) string {
   re := regexp.MustCompile(`url\s*\(['|"]?(.*)['|"]?\)`);
   return re.ReplaceAllStringFunc(css, func(urlRule string) string {
      var link = re.ReplaceAllString(urlRule, "$1");

      // TODO(eriq): Potential problem is url is unescaped (because of quotes).
      return fmt.Sprintf("url('%s')", identifyAndFixLink(link, targetUrl));
   });
}

func fixLink(link string, linkType int, targetUrl *url.URL) string {
   url, err := targetUrl.Parse(link);

   if (err != nil) {
      return link;
   }

   return fmt.Sprintf("/?type=%d&target=%s", linkType, url.String());
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
