document.addEventListener("DOMContentLoaded", () => {
  let socket = new WebSocket("ws://localhost:8081/wsbirge");

  socket.onopen = function () {
    console.log("Соединение установлено.");
    socket.send(JSON.stringify({ type: "getWsBinanceData" }));
  };

  // let circles = [];
  // append current svg, new line, circle
  // time, last 4 elems, y - dollar

  function render(points) {
    // console.log(points.length);
    let circles = [];
    let polyline = document.getElementById("line");
    svg = document.getElementById("chartSvg");
    svg.innerHtml = "";
    if (polyline == null) {
      polyline = document.createElementNS(
        "http://www.w3.org/2000/svg",
        "polyline"
      );
      polyline.setAttributeNS(null, "id", "line");
      svg.appendChild(polyline);
    }

    let polylinePoints = [];
    for (let i = 0; i < points.length; i++) {
      polylinePoints.push(points[i].x + ", " + points[i].y);
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
          40
        ); /* This is the radius of the circle */
        circle.setAttributeNS(
          null,
          "class",
          "point"
        ); /* You can style individual shapes using CSS */
        circles.push(circle);
        svg.appendChild(circle);
      }
      circle.setAttributeNS(null, "cx", points[i].x);
      circle.setAttributeNS(null, "cy", points[i].y);
    }
    /* In case we modify the number of points */
    if (points.length < circles.length) {
      for (; i < circles.length; i++) {
        circles[i].remove();
      }
      circles.splice(points.length, circles.length);
    }
    polyline.setAttributeNS(null, "points", polylinePoints.join(" "));
  }

  socket.onmessage = function (e) {
    let app = document.getElementById("app");
    let chart = document.getElementById("chartSvg");
    let message = JSON.parse(e.data);
    //array data -? paint graphic ||
    //prev innerHtml = "", create new dynamic js
    switch (message[0].type) {
      case "newdata":
        let points = [];
        for (let [idx, v] of Object.entries(message)) {
          // app.textContent = ` ${v.bids}  : ${v.asks}, ${v.spread}`;
          // points.push({ x: v.time, y: v.maxask });
          // let t = v.time.toString().slice(5)
          console.log(10 + idx);
          points = [...points, { x: 100 + idx, y: v.maxask }];
          render(points);
        }
        //1, {110, 41.000}
        //2, {111, 41.200}
        socket.send(JSON.stringify({ type: "getWsBinanceData" }));
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
