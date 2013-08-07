package reserve;

import (
   "html/template"
);

var BasePageTemplate =
   template.Must(template.New("basePage").Parse(BasePageTemplateStr));

const BasePageTemplateStr = `{{.}}`;

var TopBar string =
   `
   <div class='_reserver__container'>
      <div class="_reserver__placeholder"></div>
      <div class='_reserver__bar'>
         <span class='_reserver__span'>http://</span>
         <form class='_reserver__form' action="/">
            <input class="_reserver__input"
                     name="target"
                     type="text"
                     placeholder="www.myCoolSite.com" />
            <input type="hidden" name="type" value=0 />
            <input type="submit" class='_reserver__button' value="Go!">
         </form>
      </div>
   </div>
   `;

var InjectedStyle =
   `
   <style>
      ._reserver__container {
      }

      ._reserver__bar {
         background-color: #2d2d2d;
         border-bottom: 2px solid white;
         padding: 5px;
         width: 100%;
         position: fixed;
         top: 0;
         left: 0;
         z-index: 999999;
      }

      ._reserver__placeholder {
         width: 100%;
         height: 32px;
         margin: 0;
         padding: 0;
      }

      ._reserver__span {
         color: white;
      }

      ._reserver__form {
         display: inline-block;
         margin: 0;
      }

      ._reserver__input {
         min-width: 512px;
      }

      ._reserver__button {
      }
   </style>
   `;
