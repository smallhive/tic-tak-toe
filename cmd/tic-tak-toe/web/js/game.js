const TypeInit = 1
const TypeConnect = 2
const TypeGameStarted = 3
const TypeYouTurn = 4
const TypeNotYouTurn = 5
const TypeStep = 6
const TypeFieldUpdate = 7
const TypeGameEnded = 8
const TypeGameFailed = 9
const TypeOpponentUnexpectedDisconnect = 11
const TypeSetNick = 12
const TypeSetOpponentNick = 13
const TypeAreYouReady = 14
const TypeIamReady = 15

let playerMark = undefined;
let playerNick = undefined;

let conn = undefined;

let turnNotificationCell = undefined;
let gameResult = undefined;

function subscribeCellHandlers() {
    document.querySelectorAll('.game-cell').forEach(item => {
        item.addEventListener('click', clickHandler)
    })
}

function unSubscribeCellHandlers() {
    document.querySelectorAll('.game-cell').forEach(item => {
        item.removeEventListener('click', clickHandler)
    })
}

function ready() {
    if (!window["WebSocket"]) {
        let item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        return;
    }

    while (true) {
        playerNick = prompt("You should provide you nick. Min 3 symbols");
        if (playerNick.length >= 3) {
            break;
        }
    }

    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onclose = function (evt) {
        console.log('Connection closed.');
    };

    conn.onmessage = function (evt) {
        let e = JSON.parse(evt.data);
        handleEvent(e);
    };

    conn.onopen = function (evt) {
        setPlayerNick(playerNick);

        turnNotificationCell = document.getElementById('turnNotification');
        gameResult = document.getElementById('gameResult');
        gameResult.innerHTML = 'We are waiting for opponent';
    }
}

function handleEvent(e) {
    switch (e.type) {
        case TypeInit:
            playerMark = e.data.Label;
            document.title = 'GameID ' + e.data.GameID;

            let playerMarkCell = document.getElementById('playerMark');
            playerMarkCell.innerHTML = 'You mark is ' + playerMark;
            break;
        case TypeGameStarted:
            gameResult.innerHTML = '';
            sendPlayerNickUpdate(playerNick);
            subscribeCellHandlers();
            break;
        case TypeYouTurn:
            turnNotificationCell.innerHTML = 'You turn';
            turnNotificationCell.classList.remove('opponent-step');
            turnNotificationCell.classList.add('you-step');
            break;
        case TypeNotYouTurn:
            turnNotificationCell.innerHTML = 'Oponent turn';
            turnNotificationCell.classList.remove('you-step');
            turnNotificationCell.classList.add('opponent-step');
            break;
        case TypeFieldUpdate:
            let field = e.data.Field;
            for (let i = 0; i < field.length; i++) {
                for (let j = 0; j < field[i].length; j++) {
                    let id = `cell_${i}_${j}`;
                    let cell = document.getElementById(id);
                    cell.innerHTML = field[i][j];
                }
            }
            break;
        case TypeGameEnded:
            if (e.data.IsWin) {
                document.title = 'Congrats! You are winner!';
            } else {
                document.title = 'My sorry! You are lose!';
            }

            e.data.Condition.forEach(point => {
                let id = `cell_${point[0]}_${point[1]}`;
                let cell = document.getElementById(id);
                cell.classList.add('win-condition');
            });

            gameResult.innerHTML = document.title;
            turnNotificationCell.innerHTML = '';
            unSubscribeCellHandlers();
            break
        case TypeGameFailed:
            document.title = 'No winners! Reconnect to attempt one more time ;)';
            gameResult.innerHTML = document.title;
            turnNotificationCell.innerHTML = '';
            unSubscribeCellHandlers();
            break
        case TypeOpponentUnexpectedDisconnect:
            gameResult.innerHTML = 'Opponent has left. You win!';
            unSubscribeCellHandlers();
            break;
        case TypeSetOpponentNick:
            setOpponentNick(e.data.Nick);
            break;
        case TypeAreYouReady:
            sendIamReady();
            break;
    }
}

const clickHandler = function (event) {
    let cell = event.currentTarget;

    let rowId = cell.getAttribute('data-row-id');
    let collId = cell.getAttribute('data-coll-id');

    let click = {
        "type": TypeStep,
        "data": {
            "Row": parseInt(rowId),
            "Coll": parseInt(collId),
        }
    }

    conn.send(JSON.stringify(click));
}

function setPlayerNick(nick) {
    let label = document.getElementById('playerNick');
    label.innerHTML = nick;
}

function sendPlayerNickUpdate(nick) {
    let cmd = {
        "type": TypeSetNick,
        "data": {
            "Nick": nick,
        }
    }
    conn.send(JSON.stringify(cmd));
}

function sendIamReady() {
    let cmd = {
        "type": TypeIamReady,
        "data": {}
    }
    conn.send(JSON.stringify(cmd));
}


function setOpponentNick(nick) {
    let label = document.getElementById('opponentNick');
    label.innerHTML = nick;
}
