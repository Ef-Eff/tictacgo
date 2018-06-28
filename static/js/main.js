$(() => {
  const ws = new WebSocket("ws://" + document.location.host + "/ws");

  let playerNumber;

  const $circle = $("<span class='nought'></span>");
  const $cross = $("<span class='cross'><span></span><span></span></span>");
  const fixYourShit = [$circle, $cross];

  function color() {
    return ["red", "blue"][playerNumber - 1];
  }

  function welcome(number) {
    playerNumber = number;
    console.log(`Welcome! You are now player ${playerNumber}`);
    $("#player").text(`Player ${playerNumber}`).css("color", color());
  }

  // no fucking idea
  function start() {
    console.log(`Two players have been found and now the match begins! You go ${["1st", "2nd"][playerNumber-1]}`);
    $("div>div").click(sendMark);
    $("h3#shitTalk").text(["Your turn", "Opponents turn"][playerNumber-1])
  }
  
  function mark(msg) {
    console.log(msg);
    const $boardPos = $(`div>div[data-pos="${msg.Position}"]`);
    $boardPos.off().append(fixYourShit[msg.PlayerNumber].clone());
  }
  
  function win(msg) {
    mark(msg)
    $("div>div").off();
    $(`div>div[data-keys*="${msg.Key}"]`).css("background-color", "green");
    setTimeout(() => {
      alert(msg.PlayerNumber === playerNumber ? "You won!": "Ur shit")
    }, 0)
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
    console.log("Server closed. I dunno doood, i dunno.");
  }
})