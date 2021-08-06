document.addEventListener("DOMContentLoaded", () => {
  let socket = new WebSocket("ws://localhost:8081/wsbirge");

  socket.onopen = function () {
    console.log("Соединение установлено.");
    socket.send(JSON.stringify({ type: "getWsBinanceData" }));
  };

  let tickerask = 0;
  let tickerbid = 0
  let polyPointsAsks = [];
  let polyPointsBids = [];

  function renderBids(point) {

    var polyline = document.getElementById(`linemaxbid`);
    if (polyline == null) {
      polyline = document.createElementNS(
        "http://www.w3.org/2000/svg",
        "polyline"
      );
      polyline.setAttributeNS(null, "id", `linemaxbid`);
      svg = document.getElementById(`chartbid`);
      svg.appendChild(polyline);
    }
    
    let circles = [];
    tickerbid += 100;
      polyPointsBids.push(tickerbid + ", " + Math.round(point));
      let  circle = document.createElementNS(
          "http://www.w3.org/2000/svg",
          "circle"
        );
        circle.setAttributeNS(
          null,
          "r",
          75
        ); 
        circle.setAttributeNS(null, "class", "point");
        circles.push(circle);
        svg.appendChild(circle);
      circle.setAttributeNS(null, "cx", tickerbid);
      circle.setAttributeNS(null, "cy", point);
    polyline.setAttributeNS(null, "points", polyPointsBids.join(" "));
  }


  function renderAsks(point) {

    var polyline = document.getElementById(`lineminask`);
    if (polyline == null) {
      polyline = document.createElementNS(
        "http://www.w3.org/2000/svg",
        "polyline"
      );
      polyline.setAttributeNS(null, "id", `lineminask`);
      svg = document.getElementById(`chartask`);
      svg.appendChild(polyline);
    }
    
    let circles = [];
    tickerask += 100;
      polyPointsAsks.push(tickerask + ", " + Math.round(point));
      let  circle = document.createElementNS(
          "http://www.w3.org/2000/svg",
          "circle"
        );
        circle.setAttributeNS(
          null,
          "r",
          75
        ); 
        circle.setAttributeNS(null, "class", "point");
        circles.push(circle);
        svg.appendChild(circle);
      // }
      circle.setAttributeNS(null, "cx", tickerask);
      circle.setAttributeNS(null, "cy", point);
 
    polyline.setAttributeNS(null, "points", polyPointsAsks.join(" "));
  }



  socket.onmessage = function (e) {
    let message = JSON.parse(e.data);
    let spread = document.getElementById("spread");
    let diff = document.getElementById("diff")
    let maxBid = document.getElementById("maxBid");
    let minAsk = document.getElementById("minAsk");
    
    switch (message.type) {
      case "newdata":
        spread.textContent = `Spread data: ${message.spread} $`;
        maxBid.textContent = `Max bid: ${message.maxbid} $`;
        minAsk.textContent = `Min ask: ${message.minask} $`;
        diff.textContent = ` Max difference  ask bids: ${message.maxdiff} $`;
       
        //ask, bid
        renderBids(message.maxbid );
        renderAsks(message.minask);
        // setTimeout(function () {
        socket.send(JSON.stringify({ type: "getWsBinanceData" }));
      // }, 500);
      default:
        console.log("default case");
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
