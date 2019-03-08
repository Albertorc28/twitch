
$(document).ready(function () {

    $("#btnIniciar").click(function() {
        IniciarSesion();
    });

    $("#txtUser").keyup(function(e) {
        if(e != null) {    
            var code = e.which;
            if(code == 13) {
                IniciarSesion();
            }
        }
    });

    $("#txtPass").keyup(function(e) {
        if(e != null) {
            var code = e.which;
            if(code == 13) {
                IniciarSesion();
            }
        }
    });
});

function IniciarSesion() {
    $.ajax({
        url: "/login",
        method: "POST",
        contentType: "application/json",
        dataType: "json",
        data: JSON.stringify({ User: $("#txtUser").val(), Password: $("#txtPass").val()})
    })
    .done(function(data) {
        console.log("OK");
        console.log(data);
        window.location.href = "/";
    })
    .fail(function(data) {
        console.log("KO");
        console.log(data);
    });
}