document.addEventListener("DOMContentLoaded", () => {
  let socket = new WebSocket("ws://localhost:8081/wsbirge");

  socket.onopen = function () {
    console.log("Соединение установлено.");
    socket.send(JSON.stringify({ type: "getWsBinanceData" }));
  };

  socket.onmessage = function (e) {
    let app = document.getElementById("app");
    let chart = document.getElementById("chart");
    let message = JSON.parse(e.data);
    console.log("Получены данные " + Object.entries(message), e.data);

    switch (message[0].type) {
      case "newdata":
        for (let [k, v] of Object.entries(message)) {
          // app.textContent = ` ${v.bids}  : ${v.asks}, ${v.spread}`;
        }
        socket.send(JSON.stringify({ type: "getWsBinanceData" }));
      default:
        console.log("default");
    }
  };

  socket.onclose = function (event) {
    if (event.wasClean) {
      console.log("Соединение закрыто чисто");
    } else {
      console.log("Обрыв соединения"); // например, "убит" процесс сервера
    }
    console.log("Код: " + event.code + " причина: " + event.reason);
  };

  socket.onerror = function (error) {
    console.log("Ошибка " + error.message);
  };
  //svg use  paint data
});
