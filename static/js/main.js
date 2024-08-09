/**
 * @typedef {0|1|2|3|4|5|6|7|8} BoardNumbers
 *
 * @typedef {{
 *    1: {
 *      number: 1;
 *      chip: any;
 *      color: "red";
 *      css: "border-color";
 *      target: "span.nought"
 *    };
 *    2: {
 *      number: 2;
 *      chip: any;
 *      color: "blue";
 *      css: "background-color";
 *      target: "span.cross>span";
 *    };
 *  }} Players
 *
 * @typedef {Players[1]|Players[2]} Player
 *
 * @typedef {"welcome"|"start"|"mark"|"win"|"draw"|"winbydc"|"error"} MessageType
 *
 * @typedef {{
 *  PlayerNumber: 1|2;
 *  Position: BoardNumbers;
 * }} ActionMessage
 *
 * @typedef {{
 *  WinPos: [BoardNumbers, BoardNumbers, BoardNumbers];
 * } & ActionMessage} WinMessage
 *
 * @typedef {{
 *    Type: MessageType;
 *    Data: any;
 * }} Message
 */

$(() => {
  const ws = document.location.host.includes('localhost:')
    ? new WebSocket(`ws://${document.location.host}/ws`)
    : new WebSocket(`wss://${document.location.host}/ws`);

  /** @type {Player} */
  let player;

  const $circle = $("<span class='nought'></span>");
  const $cross = $("<span class='cross'><span></span><span></span></span>");
  const $loader = $('span#loader');
  const $topText = $('h3#top-text');
  const $infoText = $('h3#info-text');
  const turns = ['Opponents turn', 'Your turn'];

  /** @type {Players} */
  const players = {
    1: {
      number: 1,
      chip: $circle,
      color: 'red',
      css: 'border-color',
      target: 'span.nought',
    },
    2: {
      number: 2,
      chip: $cross,
      color: 'blue',
      css: 'background-color',
      target: 'span.cross>span',
    },
  };

  const actions = {
    welcome,
    start,
    mark,
    win,
    draw,
    winbydc,
    error,
  };

  /** @param {1|2} number */
  function welcome(number) {
    player = players[number];
    console.log(`Welcome! You are now player ${number}`);
    $('#player').text(`Player ${number}`).css('color', player.color);
  }

  function start() {
    $loader.remove();
    $infoText.text('Connected');

    $('div>div')
      .click(sendMark)
      .hover(
        function () {
          $(this).append(players[player.number].chip.clone().addClass('hover'));
        },
        function () {
          $(this).empty();
        }
      );

    $topText.text(turns[player.number % 2]);
    $('div.off').toggleClass('off');
  }

  /**
   * @param {ActionMessage} msg
   * @param {Boolean} bool
   */
  function mark(msg, bool) {
    const $boardPos = $(`div>div[data-pos="${msg.Position}"]`);
    $boardPos.empty().off().append(players[msg.PlayerNumber].chip.clone());
    if (!bool) {
      $topText.text(turns[Math.abs(msg.PlayerNumber - player.number)]);
    }
  }

  /** @param {WinMessage} msg */
  function win(msg) {
    mark(msg, true);
    $('div>div').off();
    const { color, target, css } = players[msg.PlayerNumber];
    const positions = `
      div>div[data-pos="${msg.WinPos[0]}"]>${target},
      div>div[data-pos="${msg.WinPos[1]}"]>${target},
      div>div[data-pos="${msg.WinPos[2]}"]>${target}`;
    $(positions).css(css, color);
    $topText.text(
      msg.PlayerNumber === player.number ? 'You won!' : 'You lost.'
    );
  }

  /** @param {ActionMessage} msg */
  function draw(msg) {
    mark(msg, true);
    $('div>div').off();
    $topText.text('Draw!');
  }

  function winbydc() {
    $('div>div').off();
    $topText.text('The other player disconnected.');
  }

  /** @param {string} errMessage */
  function error(errMessage) {
    alert(errMessage);
  }

  function sendMark() {
    ws.send(parseInt($(this).attr('data-pos')));
  }

  ws.onmessage = function (event) {
    /** @type {Message} */
    const data = JSON.parse(event.data);
    actions[data.Type](data.Data);
  };

  ws.onclose = function () {
    $('div>div').off();
    $infoText.text('Refresh the page to start a new game.');
  };
});
