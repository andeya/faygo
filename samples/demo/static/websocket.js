var wsUri = "ws://localhost:8080/ws_server";
var ws = new WebSocket(wsUri);
ws.onopen = function() {
    console.log("connected to " + wsUri);
    var msg = "get server time";
    ws.send(JSON.stringify({
        "testsend": msg
    }));
    console.log("send:", msg);
};
ws.onclose = function(e) {
    console.log("connection closed (" + wsUri + " : " + e.code + "," + e.reason + ")");
}
ws.onerror = function(e) {
    for (var p in e) {
        console.log(p + "=" + e[p]);
    }
};
ws.onmessage = function(m) {
    console.log("receive:", m.data);
};
window.onbeforeunload = function() {
    ws.close();
    console.log("closed websocket");
    return
}
