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
      <div class="_reserver__placeholder"><p>|</p></div>
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
         height: 38px;
         margin: 0;
         padding: 0;
         font-size: 18px;
      }

      ._reserver__span {
         color: white;
         font-size: 18px;
         vertical-align: middle;
      }

      ._reserver__form {
         display: inline-block;
         margin: 0;
      }

      ._reserver__input {
         min-width: 512px;
         font-size: 18px;
         vertical-align: middle;
      }

      ._reserver__button {
         height: 24px;
         width: 32px;
         vertical-align: middle;
         position: static;
      }
   </style>
   `;

var StartPage =
   `
   <html>
      <head>
         <title>Reserver</title>

         <style>
            body {
               background-color: #121211;
            }

            ._reserver__container {
               position: absolute;
               top: 25%;
               left: 25%;
               background-color: #2b2b2b;
               border-radius: 10px;
               width: 700px;
               height: 300px
               padding: 10px;
            }

            ._reserver__text {
               font-weight: bold;
               color: white;
               font-size: 24px;
               text-align: center;
               font-family: sans-serif;
            }

            ._reserver__form {
            }

            ._reserver__input {
               display: block;
               font-size: 24px;
               width: 650px;
               margin-left: 25px;
               text-align: center;
            }

            /* Button gradients taken from imgur */
            ._reserver__button {
               display: block;

               background: #2b2b2b;
               background: -moz-linear-gradient(top,#2b2b2b 0,#121211 100%);
               background: -webkit-gradient(linear,left top,left bottom,color-stop(0%,#2b2b2b),color-stop(100%,#121211));
               background: -webkit-linear-gradient(top,#2b2b2b 0,#121211 100%);
               background: -o-linear-gradient(top,#2b2b2b 0,#121211 100%);
               background: -ms-linear-gradient(top,#2b2b2b 0,#121211 100%);
               background: linear-gradient(to bottom,#2b2b2b 0,#121211 100%);
               filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#2b2b2b', endColorstr='#121211', GradientType=0);

               border: 2px solid #444442;
               border-radius: 8px;

               width: 600px;
               height: 50px;
               margin-left: 50px;
               margin-top: 20px;
               font-size: 24px;
               font-family: monospace;
               color: white;
               font-weight: bold;
               cursor: pointer;
            }

            ._reserver__button:hover {
               background:#2b2b2b;
               background:-moz-linear-gradient(top,#2b2b2b 0,#444442 0,#121211 100%);
               background:-webkit-gradient(linear,left top,left bottom,color-stop(0%,#2b2b2b),color-stop(0%,#444442),color-stop(100%,#121211));
               background:-webkit-linear-gradient(top,#2b2b2b 0,#444442 0,#121211 100%);
               background:-o-linear-gradient(top,#2b2b2b 0,#444442 0,#121211 100%);
               background:-ms-linear-gradient(top,#2b2b2b 0,#444442 0,#121211 100%);
               background:linear-gradient(to bottom,#2b2b2b 0,#444442 0,#121211 100%);
               filter:progid:DXImageTransform.Microsoft.gradient(startColorstr='#444442', endColorstr='#121211', GradientType=0);
            }

            ._reserver__button:active {
               background: #121211;
               filter:progid:DXImageTransform.Microsoft.gradient(startColorstr='#121211', endColorstr='#121211', GradientType=0);
            }
         </style>
      </head>
      <body>
         <div class='_reserver__container'>
            <p class='_reserver__text'>Put a full URL in the box and GO!</p>
            <form class='_reserver__form' action="/">
               <input class="_reserver__input"
                        name="target"
                        type="text"
                        placeholder="http://www.myCoolSite.com" />
               <input type="hidden" name="type" value=0 />
               <input type="submit" class='_reserver__button' value="GO!">
            </form>
         </div>
      </body>
   </html>
   `;
