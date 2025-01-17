const menu = document.getElementById('menu');
const roomNameInput = document.getElementById('chatNameInput');
const chatSection = document.getElementById('chatSection');
const roomTitle = document.getElementById('chatTitle');
const leaveButton = document.getElementById('leaveChatBtn');
const chatWindow = document.getElementById('chatWindow');

const chatForm = document.getElementById('chatForm');
const sendMsgBtn = document.getElementById('sendMessageBtn');
const msgInput = document.getElementById('messageInput');

const userList = document.getElementById('userList');

let userId = undefined;
initUserId();

document.getElementById('joinChatBtn').addEventListener('click', (e) => {
    enableChat(roomNameInput.value);
    joinRoom();
});

document.getElementById('logoutBtn').addEventListener('click', (e) => {
    document.cookie = 'session=; Max-Age=0';
    window.location.reload();
})

const joinRoom = () => {
    const roomName = roomNameInput.value.trim();
    
    let socket = new WebSocket(`ws://${window.location.host}/ws/joinRoom/${roomName}`);
    
    chatForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        if (msgInput.value.length == 0) {
            return
        }
    
        socket.send(msgInput.value);
    });

    socket.onopen = function(e) {
        leaveButton.addEventListener('click', () => {
            socket.close();
            disableChat();
        });

        roomTitle.value = 'Room: ' + roomNameInput.value;
    };
      
    socket.onmessage = function(e) {
        const message = JSON.parse(e.data);
        takeMessage(message);
    };
};

const errorMessage = (message) => {
    const errMsg = document.querySelector("#err-msg");
    errMsg.classList.remove("d-none");
    errMsg.textContent = message;
};

const enableChat = (roomName) => {
    roomTitle.textContent = `Room: ${roomName}`;
    menu.classList.add('d-none');
    chatSection.classList.remove('d-none');
    leaveButton.classList.remove('d-none');
};

const disableChat = () => {
    roomTitle.textContent = '';
    menu.classList.remove('d-none');
    chatSection.classList.add('d-none');
    leaveButton.classList.add('d-none');
};

const takeMessage = (message) => {
    let msg = '';

    if (message.type === "event") {
        updateUserList(roomNameInput.value);
    }

    if (message.userId === userId) {
        msg = `
        <div class="d-flex mb-2">
            <div class="bg-${message.color} text-white p-2 rounded ms-auto">${message.content}</div>
        </div>`;
    }else {
        msg = `
        <div class="d-flex mb-2">
            <div class="bg-${message.color} text-white p-2 rounded me-auto">${message.content}</div>
        </div>`;
    }

    chatWindow.innerHTML = chatWindow.innerHTML + msg;
};

function updateUserList(roomName) {
    fetch(`/ws/getClients/${roomName}`)
    .then((resp) => {
        resp.json()
        .then((data) => {
            userList.innerHTML = '';

            for (let user of data) {
                const newUser = `<li class="bg-${user.color} text-white list-group-item">${user.username}</li>`
                userList.innerHTML += newUser;
            }
        })
    })
    .catch(err => console.error(err));
}

function initUserId() {
    fetch('/get-session')
    .then(res => res.json()
    .then(data => userId = data.userId));
}

function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
}
