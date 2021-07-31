console.log("start")
    document.addEventListener("DOMContentLoaded", () => {
        //ws use ? 
        // var socket = new WebSocket("ws://localhost:8080/");
        const socket = new WebSocket.Server({ port: 8080 });
        // let res = fetch(`http://localhost:8080/birge`) 

        socket.onopen = function() {
            alert("Соединение установлено.");
          };
          
          socket.onclose = function(event) {
            if (event.wasClean) {
              alert('Соединение закрыто чисто');
            } else {
              alert('Обрыв соединения'); // например, "убит" процесс сервера
            }
            alert('Код: ' + event.code + ' причина: ' + event.reason);
          };
          
          socket.onmessage = function(event) {
            alert("Получены данные " + event.data);
          };
          
          socket.onerror = function(error) {
            alert("Ошибка " + error.message);
          };
        console.log(res, 1)
        //svg use  paint data

      });
       