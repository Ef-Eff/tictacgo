$(() => {
  const ws = new WebSocket("ws://" + document.location.host + "/ws");

  let player;

  const $circle = $("<span class='nought'></span>");
  const $cross = $("<span class='cross'><span></span><span></span></span>");
  const $loader = $("span#loader");
  const $info = $("h3#info");

  const players = {
    1: {
      number: 1,
      chip: $circle,
      color: "red",
      css: "border-color",
      target: "span.nought"
    },
    2: {
      number: 2,
      chip: $cross,
      color: "blue",
      css: "background-color",
      target: "span.cross>span"
    }
  };

  function welcome(number) {
    player = players[number];
    console.log(`Welcome! You are now player ${number}`);
    $("#player").text(`Player ${number}`).css("color", player.color);
  }

  // no fucking idea
  function start() {
    $loader.remove();
    $info.text("Connected")
    $("div>div").click(sendMark).hover(function(){
      $(this).append(players[player.number].chip.clone().addClass("hover"));
    }, function() {
      $(this).empty();
    });
    $("h3#shitTalk").text(["Your turn", "Opponents turn"][player.number-1])
    $("div.off").toggleClass("off");
  }
  
  function mark(msg) {
    console.log(msg);
    const $boardPos = $(`div>div[data-pos="${msg.Position}"]`);
    $boardPos.off().append(players[msg.PlayerNumber].chip.clone());
    $("h3#shitTalk").text(["Opponents turn", "Your turn"][Math.abs(msg.PlayerNumber - player.number)])
  }
  
  function win(msg) {
    mark(msg)
    $("div>div").off();
    const { color, target, css } = players[msg.PlayerNumber];
    $(`div>div[data-keys*="${msg.Key}"]>${target}`).css(css, color);
    // I think that the below code makes it so that the UI updates before the alert, if my JS knowledge is good that is.
    // Im sure I was getting the alert popping up before the js could update the view, now it doesnt
    $("h3#shitTalk").text(msg.PlayerNumber === player.number ? "You won!": "Ur shit");
    $("h3#info").text("Refresh the page to start a new game.");
    setTimeout(() => {
      console.log(msg.PlayerNumber === player.number ? "You won!": "Ur shit");
    }, 0);
  }
  
  // :thinking: maybe later
  // function lose() {
  
  // }

  function error(errMessage) {
    alert(errMessage);
  }
  
  const actions = {
    welcome,
    start,
    mark,
    win,
    // lose,
    error,
  };

  function sendMark() {
    $boardPos = $(this);
    const message = {
      Position: parseInt($boardPos.attr("data-pos")), 
      Keys: $boardPos.attr("data-keys").split(" ") 
    };
    console.log("Sending play:", message);
    ws.send(JSON.stringify(message));
  }

  ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    actions[data.Type](data.Data);
  }

  ws.onclose = function(event) {
    // alert("Server closed. Refresh to reconnect");
  }
})