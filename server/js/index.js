function logout() {
    $.ajax({
        url: "/logout",
        method: "POST"
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