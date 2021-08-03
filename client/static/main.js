document.addEventListener("DOMContentLoaded", () => {
  let socket = new WebSocket("ws://localhost:8081/wsbirge");

  socket.onopen = function () {
    console.log("Соединение установлено.");
    socket.send(JSON.stringify({ type: "getWsBinanceData" }));
  };

  let ticker = 0;
  let polylinePoints = [];

  function render(points) {
    var polyline = document.getElementById("line");
    if (polyline == null) {
      polyline = document.createElementNS(
        "http://www.w3.org/2000/svg",
        "polyline"
      );
      polyline.setAttributeNS(null, "id", "line");
      svg = document.getElementById("chartSvg");
      svg.appendChild(polyline);
    }
    points[0].maxask = points[0].maxask;
    let circles = [];
    ticker += 100;
    for (let i = 0; i < points.length; i++) {
      polylinePoints.push(ticker + ", " + Math.round(points[i].maxask));
      let circle;
      // console.log(circles.length);
      if (i < circles.length) {
        circle = circles[i];
      } else {
        circle = document.createElementNS(
          "http://www.w3.org/2000/svg",
          "circle"
        );
        circle.setAttributeNS(
          null,
          "r",
          75
        ); /* This is the radius of the circle */
        circle.setAttributeNS(null, "class", "point");
        circle.textContent = points[i].minbid;
        /* You can style individual shapes using CSS */
        circles.push(circle);
        svg.appendChild(circle);
      }
      circle.setAttributeNS(null, "cx", ticker);
      circle.setAttributeNS(null, "cy", points[i].maxask);
    }
    /* In case we modify the number of points */
    if (points.length < circles.length) {
      for (; i < circles.length; i++) {
        circles[i].remove();
      }
      circles.splice(points.length, circles.length);
    }
    polyline.setAttributeNS(null, "points", polylinePoints.join(" "));
    // console.log(polyline.getAttribute("points"), 100);
  }
  socket.onmessage = function (e) {
    let message = JSON.parse(e.data);
    let spread = document.getElementById("spread");

    let maxBid = document.getElementById("maxBid");
    switch (message.type) {
      case "newdata":
        spread.textContent = `Spread data: ${message.spread} $`;
        maxBid.textContent = ` Max bid: ${message.maxbid} $`;
        //minAsk todo

        render([message]);
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
