package reserve;

import (
   "html/template"
);

var BasePageTemplate =
   template.Must(template.New("basePage").Parse(BasePageTemplateStr));

const BasePageTemplateStr = `{{.}}`;
