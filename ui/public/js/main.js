$(document).ready(function (){
    const api = 'http://127.0.0.1:8080/api/v1';
    const pageIds = {
        loading: 'loading',
        games: 'games',
        game: 'game',
    }
    const pageElems = {
        loading: $('#'+pageIds.loading),
        games: $('#'+pageIds.games),
        game: $('#'+pageIds.game),
    }
    switchPages(pageIds.loading);

    const blankBoard = '---------';

    const gameStatusRunning = "RUNNING";
    const gameStatusCrossWon = "X_WON";
    const gameStatusNoughtWon = "O_WON";
    const gameStatusDraw = "DRAW";

    getGames(function (games) {
        games.forEach(function (game) {
            $('#games-list').append(
                drawGame(game),
            );
        });
        switchPages(pageIds.games);
    });

    $('#games-new').append(drawGameBoard('', blankBoard, gameStatusRunning));

    $('#computer').click(function () {
        createGame(blankBoard, function (id) {
            getGame(id, function (game) {
                $('#games-list').append(
                    drawGame(game),
                );
            });
        });
    });

    function getGames(callback) {
        $.ajax({
            url: api + '/games',
            type: 'GET',
            success: function (games) {
                callback(games);
            },
            error: onError,
        });
    }

    function getGame(id, callback) {
        $.ajax({
            url: api + '/games/'+id,
            type: 'GET',
            success: callback,
            error: onError,
        });
    }

    function createGame(board, callback) {
        $.ajax({
            url: api + '/games',
            type: 'POST',
            data: JSON.stringify({
                board: board,
            }),
            dataType: "json",
            success: function (location) {
                let id = location.location.substring(location.location.lastIndexOf('/') + 1);
                callback(id);
            },
            error: onError,
        });
    }

    function updateGame(id, board, callback) {
        $.ajax({
            url: api + '/games/'+id,
            type: 'PUT',
            data: JSON.stringify({
                board: board,
            }),
            dataType: "json",
            success: callback,
            error: onError,
        });
    }

    function switchPages(show) {
        $('#error').text('')
        for (page in pageElems) {
            if (show === page) {
                pageElems[page].show();
            } else {
                pageElems[page].hide();
            }
        }
    }

    function onError(error) {
        console.log(error);
        $('#error').html(error.status+' '+error.statusText+'<br>'+error.responseText);
    }

    function drawGame(game) {
        let html = $('<div></div>');
        html.append('<div>ID: '+game.id+'</div>');
        html.append(drawGameBoard(game.id, game.board, game.status));
        switch (game.status) {
            case gameStatusDraw:
                html.append('Game is draw!');
                break;
            case gameStatusCrossWon:
                html.append('X won!');
                break;
            case gameStatusNoughtWon:
                html.append('0 won!');
                break;
            default:
                html.append('Make a move!');
        }
        html.append('<hr>');
        return html;
    }

    function drawGameBoard(id, board, status) {
        let html = '<table data-id="'+id+'" data-board="'+board+'" data-char="X">';
        let n = 0;
        for (let i = 0; i < 3; i++) {
            html += '<tr>';
            for (let j = 0; j < 3; j++) {
                html += '<td data-cell="'+n+'">'+board[n]+'</td>'
                n++
            }
            html += '</tr>';
        }
        html += '</table>';
        let elem = $(html);
        if (status === gameStatusRunning) {
            elem.find('td').click(onCellClick);
        }
        return elem
    }

    function onCellClick() {
        let n = $(this).attr('data-cell');
        let t = $(this).closest('table');
        let c = t.attr('data-char');
        let id = t.attr('data-id');
        let board = t.attr('data-board');

        if (id) {
            let chars = board.split('');
            chars[n] = c;
            board = chars.join('');
            updateGame(id, board, function (game) {
                t.closest('div').replaceWith(drawGame(game));
            })
        } else {
            let chars = board.split('');
            chars[n] = $('input[name="char"]:checked').val();
            board = chars.join('');
            createGame(board, function (id) {
                getGame(id, function (game) {
                    $('#games-list').append(
                        drawGame(game),
                    );
                    t.replaceWith(drawGameBoard('', blankBoard, gameStatusRunning));
                });
            });
        }
    }
})