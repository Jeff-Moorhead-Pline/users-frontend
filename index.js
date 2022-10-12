"use strict";

// TODO: can we move these functions to another file?

// TODO: add error handling to fetch call
async function getUsers() {

    let url = "http://localhost:8080/api/users";

    await fetch(url)
    .then((response) => response.json())
    .then( createUsersTable );
}

function createUsersTable(data) {

    data.map( appendUserToTable );
}

// TODO: look up promises
async function createUser(event) {

    event.preventDefault();
    let url = "http://localhost:8080/api/users/new";
    let form = document.querySelector("#create-user-form");
    let fd = new FormData(form);
    let body = {};

    for (const [key, val] of fd) {
        body[key] = val;
    }

    await fetch(url, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then( (response) => response.json() )
    .then( (data) => { appendUserToTable(data); } );
}

function appendUserToTable(user) {

    let usersList = document.querySelector("#users-info");
    let newUserElem = document.createElement("tr");
    let id = document.createElement("td");
    let userSince = document.createElement("td");
    let userType = document.createElement("td");
    let tdClass = "user-data"

    id.innerText = user.id;
    id.setAttribute("class", tdClass);
    newUserElem.appendChild(id);

    userSince.innerText = user.userSince;
    userSince.setAttribute("class", tdClass);
    newUserElem.appendChild(userSince);

    userType.innerText = user.type;
    userType.setAttribute("class", tdClass);
    newUserElem.appendChild(userType);

    usersList.appendChild(newUserElem);
}

async function deleteUser() {

    // DELETE to /api/users/:id and refresh list
}

// Add event handler and load users on window load
document.getElementById("create-user-form").addEventListener("submit", createUser);
window.onload = getUsers;