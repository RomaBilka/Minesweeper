let newGameForm = document.getElementById("newGameForm");
let newGame = document.getElementById("newGame");

newGameForm.addEventListener("submit", (e) => {
    e.preventDefault();
    fetch("/newGame", {
        method: "POST",
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        },
        body: JSON.stringify({
            n: parseInt(document.getElementById('n').value),
            m: parseInt(document.getElementById('m').value),
            numberBlackHoles: parseInt(document.getElementById('numberBlackHoles').value),
        }),
    }).then((response) => response.json())
        .then((json) => createTable(json.cells));
});

newGame.addEventListener("click", () => {
    document.getElementById('gameOver').hidden = true;
    document.getElementById("form-container").hidden = false;
    document.getElementById("cells").hidden = true;
});

function createTable(cells) {
    const table = document.getElementById("cells");
    table.innerHTML = '';document.getElementById("form-container").hidden = true;
    for (let n in cells) {
        let tr = document.createElement('tr')
        for (let m in cells[n]) {
            let td = document.createElement('td')
            td.dataset.n = n
            td.dataset.m = m
            td.setAttribute('id', n.toString() + m.toString());
            tr.addEventListener("click", openCell);
            tr.addEventListener("contextmenu", DisabledEnabledCell);
            tr.appendChild(td)
        }
        table.appendChild(tr);
    }

    document.getElementById("form-container").hidden = true;
    table.hidden = false;
}

function openCell(e) {
    const path = "/openCell"
    const data = {
        n: parseInt(e.composedPath()[0].dataset.n),
        m: parseInt(e.composedPath()[0].dataset.m),
    }

    cellAction(data, path)
}

function DisabledEnabledCell(e) {
    e.preventDefault()
    const path = "/disabledEnabledCell"
    const data = {
        n: parseInt(e.composedPath()[0].dataset.n),
        m: parseInt(e.composedPath()[0].dataset.m),
    }

    cellAction(data, path)
}

function cellAction(data, path){
    fetch(path, {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }).then((response) => response.json())
        .then((json) => updateGame(json));
}
function updateGame(data){
    updateCells(data.cells)
    switch (parseInt(data.gameStatus)){
        case 1:
            gameOver("You won")
            break;
        case 2:
            gameOver("You los")
            break;
    }
}

function gameOver(text){
    document.getElementById('finalMessage').textContent = text
    document.getElementById('gameOver').hidden = false;
}

function updateCells(cells){
    for (let n in cells) {
        for (let m in cells[n]) {
            const id = n.toString() + m.toString()
            const td = document.getElementById(id)
            td.className = ''

            if(cells[n][m].isOpen){
                if(cells[n][m].isBlackHole){
                    td.classList.add("blackHall");
                    continue;
                }

                if(cells[n][m].numberNeighborhoodBlackHole>0){
                    td.textContent = cells[n][m].numberNeighborhoodBlackHole;
                    switch (parseInt(cells[n][m].numberNeighborhoodBlackHole)){
                        case 1:
                            td.classList.add("number1");
                            break;
                        case 2:
                            td.classList.add("number2");
                            break;
                        case 3:
                            td.classList.add("number3");
                            break;
                        default:
                            td.classList.add("number4");
                    }
                }
                td.classList.add("open");
            }
            if(cells[n][m].isDisabled){
                td.classList.add("disabled");
            }
        }
    }
}