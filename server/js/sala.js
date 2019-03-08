var urlCreator = window.URL || window.webkitURL;
var imageUrl;

function frame(anfitrion) {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/stream/" + anfitrion + "/frame");
    xhr.responseType = "blob";
    xhr.onload = function(e) {
        if(imageUrl != null) {
            URL.revokeObjectURL(imageUrl);
            delete imageUrl
        }
        imageUrl = urlCreator.createObjectURL(this.response);
        var view = document.getElementById("imgView");
        view.src = imageUrl;
        view.onload = function() {
            $(".view").css("background-image","url('"+imageUrl+"')");
        }
    }
    xhr.send();
}

var ws;

function iniciarWS() {
    var socket = new WebSocket("ws://" + window.location.host + "/stream/" + anfitrion + "/chat/" + usuario);
    socket.onopen = (function() {console.log("Socket de chat abierto")});
    socket.onmessage = (function(e) {
        m = JSON.parse(e.data);
        var parrafo = $("<p><b>" +m.Autor + "</b>: " + (m.Autor == "Bot" ? "<i>" : "") + m.Texto + (m.Autor == "Bot" ? "</i>" : "") + "</p>");
        $("#panel").append(parrafo);
        $("#panel").animate({ scrollTop: $("#panel").prop("scrollHeight")}, 300);
    });
    socket.onclose = (function(e) {
        //Deshabilitar botones de chat
        $("#panel").append("<p><b>Bot</b>: <i>El stream ha finalizado</i></p>");
        $("#panel").animate({ scrollTop: $("#panel").prop("scrollHeight")}, 300);
        console.log("Socket de chat cerrado");
        console.log(e);
    });
    return socket;
}

$(document).ready(function () {
    if(anfitrion != '') {
        setInterval(function() {
            frame(anfitrion)
        }, 200);
    }

   if(window.WebSocket === undefined) {
       alert("Tu navegador no soporta websockets. El chat está desactivado");
   }
   else {
       if(anfitrion != "") {
           ws = iniciarWS();
       }
   }

   $("#btnEnviar").click(function() {
       enviarMensaje();
   });
   
   $("#txtMensaje").keydown(function(e) {
       if(e.keyCode == 13) {
           if(!e.shiftKey) {
               e.preventDefault();
               enviarMensaje();
           }
       }
   });

   $("#fullscreen").click(function () {
       fullscreen();
   });

});

function enviarMensaje() {
    var texto = $("#txtMensaje").val();
    if(texto != null && texto != '') {
        texto = texto.replace("\n", "<br>");
        var mensaje = {
            Anfitrion: anfitrion,
            Autor: usuario,
            Texto: texto 
        };
        if(ws.readyState === ws.OPEN) {
            ws.send(JSON.stringify(mensaje));
        }
        $("#txtMensaje").val('');
    }
}

function closeStream() {
    $.ajax({
        url: window.location.protocol + "//" + window.location.host + "/stream/" + anfitrion + "/exit",
        method: "POST"
    })
    .done(function() {
        console.log("CloseStream: OK");
    })
    .fail(function() {
        console.log("CloseStream: KO");
    });
}

function fullscreen() {
    var elem = document.getElementById("view");
    if (elem.requestFullscreen) {
        elem.requestFullscreen();
      } else if (elem.mozRequestFullScreen) { /* Firefox */
        elem.mozRequestFullScreen();
      } else if (elem.webkitRequestFullscreen) { /* Chrome, Safari and Opera */
        elem.webkitRequestFullscreen();
      } else if (elem.msRequestFullscreen) { /* IE/Edge */
        elem.msRequestFullscreen();
      }
}

