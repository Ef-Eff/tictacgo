$(() => {
  let playerNumber;

  function color() {
    return ["red", "blue"][playerNumber - 1]
  }

  function welcome(number) {
    playerNumber = number
    console.log(`Welcome! You are now player ${playerNumber}`)
  }

  // no fucking idea
  function start() {
  
  }
  
  function mark(position) {
    console.log(position);
  }
  
  function win() {
  
  }
  
  function lose() {
  
  }

  function error(data) {
    alert(data.Data)
  }
  
  const actions = {
    welcome,
    start,
    mark,
    win,
    lose,
    error,
  }

  const ws = new WebSocket("ws://" + document.location.host + "/ws");

  ws.onmessage = function(event) {
    const data = JSON.parse(event.data)
    actions[data.Type](data.Data)
  }

  ws.onclose = function(event) {
    console.log(event)
    alert("Server closed. I dunno doood, i dunno.")
  }
})