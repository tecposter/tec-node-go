package web

import (
	"path/filepath"
	"html/template"
	"log"
	"net/http"
	//"fmt"
)

var lp = filepath.Join("site/public", "index.html")
var homeTemplate, err = template.ParseFiles(lp)
var staticHost = "tecposter.cn:8197"

func Home(w http.ResponseWriter, r *http.Request) {
	if err != nil {
		log.Println(err.Error() + " - " + lp)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	homeTemplate.Execute(w, "http://" + staticHost)
}



/*
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="zh-cn" prefix="og: http://ogp.me/ns#">
<head>
<title>技术客 - TecPoster</title>

<meta property="og:type" content="website">
<meta property="og:site_name" content="技术客 - TecPoster">
<meta property="og:locale" content="zh-cn">
<meta property="og:titile" content="技术客 - TecPoster">
<meta property="ob:url" content="https://www.tecposter.cn/">

<meta charset="utf-8">
<meta content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" name="viewport">

<link rel="shortcut icon" type="image/vnd.microsoft.icon" href="/favicon.ico">
<link rel="icon" type="image/png" href="/img/favicon-16x16.png" size="16x16">
<link rel="icon" type="image/png" href="/img/favicon-32x32.png" size="32x32">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta http-equiv="Content-Type" content="text/html;charset=utf-8">

<link rel="stylesheet" href="{{.}}/dev/css/main.css" />
<link rel="stylesheet" href="//at.alicdn.com/t/font_227141_3nadvz6oqsm.css">
<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/animate.css/3.7.0/animate.min.css">
</head>

<body>

<div class="page"><div class="main"></div></div>

<script type="text/javascript" src="{{.}}/dev/js/main.js"></script>

</body>
</html>
`))
*/

/*
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
*/
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
