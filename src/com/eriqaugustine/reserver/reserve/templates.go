package reserve;

import (
   "html/template"
);

var BasePageTemplate =
   template.Must(template.New("basePage").Parse(BasePageTemplateStr));

const BasePageTemplateStr = `
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
`;
