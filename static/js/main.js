$(() => {
  const ws = new WebSocket("ws://" + document.location.host + "/ws");

  let player;

  const $circle = $("<span class='nought'></span>");
  const $cross = $("<span class='cross'><span></span><span></span></span>");
  const $loader = $("span#loader");
  const $shitTalk = $("h3#shitTalk");
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
    $info.text("Connected");

    $("div>div").click(sendMark).hover(function(){
      $(this).append(players[player.number].chip.clone().addClass("hover"));
    }, function() {
      $(this).empty();
    });

    $shitTalk.text(["Your turn", "Opponents turn"][player.number-1])
    $("div.off").toggleClass("off");
  }
  
  function mark(msg, bool) {
    console.log(msg);
    const $boardPos = $(`div>div[data-pos="${msg.Position}"]`);
    $boardPos.off().append(players[msg.PlayerNumber].chip.clone());
    if (!bool) $shitTalk.text(["Opponents turn", "Your turn"][Math.abs(msg.PlayerNumber - player.number)])
  }
  
  function win(msg) {
    mark(msg, true)
    $("div>div").off();
    const { color, target, css } = players[msg.PlayerNumber];
    const positions = `div>div[data-pos="${msg.WinPos[0]}"]>${target}, div>div[data-pos="${msg.WinPos[1]}"]>${target}, div>div[data-pos="${msg.WinPos[2]}"]>${target}`
    $(positions).css(css, color);
    // I think that the below code makes it so that the UI updates before the alert, if my JS knowledge is good that is.
    // Im sure I was getting the alert popping up before the js could update the view, now it doesnt
    $shitTalk.text(msg.PlayerNumber === player.number ? "You won!": "Ur shit");
  }
  
  function draw(msg) {
    mark(msg, true);
    $("div>div").off();
    $shitTalk.text("Draw! Nobody wins! !!!11!one1!");
  }

  function winbydc() {
    $("div>div").off();
    $shitTalk.text("The other player disconnected.");
  }

  function error(errMessage) {
    alert(errMessage);
  }
  
  const actions = {
    welcome,
    start,
    mark,
    win,
    draw,
    winbydc,
    error,
  };

  function sendMark() {
    $boardPos = $(this);
    ws.send(parseInt($boardPos.attr("data-pos")));
  }

  ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    actions[data.Type](data.Data);
  }

  ws.onclose = function(event) {
    $info.text("Refresh the page to start a new game.");
  }
})