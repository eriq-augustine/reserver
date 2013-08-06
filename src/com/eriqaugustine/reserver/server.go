package main

import (
   "com/eriqaugustine/reserver/reserve"
   "code.google.com/p/go-html-transform/h5"
   "code.google.com/p/go.net/html"
   "bytes"
   "flag"
   "fmt"
   "html/template"
   "io"
   "io/ioutil"
   "log"
   "net/http"
   "net/url"
   "regexp"
   "strconv"
   "strings"
   "time"
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
   var userAgent string = request.UserAgent();
   var target string = request.FormValue("target");

   if (target == "") {
      http.NotFound(response, request);
      return;
   }

   unescapeTarget, err := url.QueryUnescape(target);

   if (err != nil) {
      fmt.Printf("Unescape Error: %s\n", err)
   }

   // TODO(eriq): Does not handle user info.
   targetUrl, err := url.Parse(unescapeTarget);

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
         var contents *string = getModifiedMain(targetUrl, userAgent);
         if (contents != nil) {
            reserve.BasePageTemplate.Execute(response, template.HTML(*contents));
            return;
         }
      case REQUEST_TYPE_IMAGE:
         if (getAndServeResource(response, request, targetUrl, userAgent, nil)) {
            return;
         }
      case REQUEST_TYPE_JS:
         if (getAndServeResource(response, request, targetUrl, userAgent, nil)) {
            return;
         }
      case REQUEST_TYPE_CSS:
         if (getAndServeResource(response, request, targetUrl, userAgent, fixCSS)) {
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

func getAndServeResource(response http.ResponseWriter, request *http.Request,
                         targetUrl *url.URL, userAgent string,
                         contentFix func(string, *url.URL) string) bool {
   var contents *[]byte = getResource(targetUrl, userAgent);
   if (contents != nil) {
      var fixedContents string = string(*contents);
      if (contentFix != nil) {
        fixedContents = contentFix(fixedContents, targetUrl);
      }

      var modTime time.Time;
      var contentReader *bytes.Reader = bytes.NewReader([]byte(fixedContents));
      http.ServeContent(response, request, targetUrl.Path, modTime,  contentReader);
      return true;
   }

   return false;
}

func getResource(targetUrl *url.URL, userAgent string) *[]byte {
   client := &http.Client{};

   request, err := http.NewRequest("GET", targetUrl.String(), nil);
   if err != nil {
      fmt.Printf("Request Error: %s\n", err);
      return nil;
   }

   request.Header.Set("User-Agent", userAgent);

   response, err := client.Do(request);
   if err != nil {
      fmt.Printf("Fetch Error: %s\n", err)
      return nil;
   }

   defer response.Body.Close();
   contents, err := ioutil.ReadAll(response.Body);

   if (err != nil) {
      fmt.Printf("Resource Read Error: %s\n", err)
      return nil;
   }

   return &contents;
}

func getModifiedMain(targetUrl *url.URL, userAgent string) *string {
   client := &http.Client{};

   request, err := http.NewRequest("GET", targetUrl.String(), nil);
   if err != nil {
      fmt.Printf("Request Error: %s\n", err);
      return nil;
   }

   request.Header.Set("User-Agent", userAgent);

   response, err := client.Do(request);
   if err != nil {
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
               var link *string = getAttr(&node.Attr, "href");
               if (link != nil) {
                  var newLink string = fixLink(*link, REQUEST_TYPE_MAIN, targetUrl);
                  replaceAttr(&node.Attr, "href", newLink);
               }
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
            case "input":
               var link *string = getAttr(&node.Attr, "src");
               if (link != nil) {
                  var newLink = identifyAndFixLink(*link, targetUrl);
                  replaceAttr(&node.Attr, "src", newLink);
               }
         }
      }
   });

   var rtn = tree.String();
   return &rtn;
}

func identifyLink(link string) int {
   urlObj, err := url.Parse(link);

   if (err != nil) {
      return REQUEST_TYPE_UNKNOWN;
   }

   if (strings.HasSuffix(urlObj.Path, ".png") ||
       strings.HasSuffix(urlObj.Path, ".jpg") ||
       strings.HasSuffix(urlObj.Path, ".jpeg") ||
       strings.HasSuffix(urlObj.Path, ".ico") ||
       strings.HasSuffix(urlObj.Path, ".gif")) {
      return REQUEST_TYPE_IMAGE;
   } else if (strings.HasSuffix(urlObj.Path, ".css")) {
      return REQUEST_TYPE_CSS;
   } else if (strings.HasSuffix(urlObj.Path, ".js")) {
      return REQUEST_TYPE_JS;
   } else {
      return REQUEST_TYPE_UNKNOWN;
   }
}

func identifyAndFixLink(link string, targetUrl *url.URL) string {
   return fixLink(link, identifyLink(link), targetUrl);
}

func fixCSS(css string, targetUrl *url.URL) string {
   re := regexp.MustCompile(`url\s*\(['"]?(.*)['"]?\)`);
   return re.ReplaceAllStringFunc(css, func(urlRule string) string {
      var link = re.ReplaceAllString(urlRule, "$1");

      // TODO(eriq): The trailing "'" should not be caught by the group.
      //  Fix the regex and remove this.
      link = strings.TrimRight(link, "'\"");

      // TODO(eriq): Potential problem is url is unescaped (because of quotes).
      return fmt.Sprintf("url('%s')", identifyAndFixLink(link, targetUrl));
   });
}

func fixLink(link string, linkType int, targetUrl *url.URL) string {
   urlObj, err := targetUrl.Parse(link);

   if (err != nil) {
      return link;
   }

   var urlString string = url.QueryEscape(urlObj.String());
   return fmt.Sprintf("/?type=%d&target=%s", linkType, urlString);
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
