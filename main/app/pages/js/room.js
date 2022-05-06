url = "localhost:8000/room/";

url += window.location.href.replace("http://localhost:8000/", "");

var exampleSocket = new WebSocket("ws://" + url);
const Action = {
  Login: 1,
  Logout: 2,
  GetUsers: 3,
  SetAdmin: 4,
  GetField: 5,
  Up: 6,
  Down: 7,
  Right: 8,
  Left: 9,
};
const ActionBack = {
  NewUser: 1,
  UserLeft: 2,
  AllUsers: 3,
  LogedIn: 4,
  Field: 5,
  Lost: 6,
};

document.getElementById("up").onclick = function () {
  var msg = {
    action: Action.Up,
  };

  exampleSocket.send(JSON.stringify(msg));
};

document.getElementById("down").onclick = function () {
  var msg = {
    action: Action.Down,
  };

  exampleSocket.send(JSON.stringify(msg));
};

document.getElementById("right").onclick = function () {
  var msg = {
    action: Action.Right,
  };

  exampleSocket.send(JSON.stringify(msg));
};

document.getElementById("left").onclick = function () {
  var msg = {
    action: Action.Left,
  };

  exampleSocket.send(JSON.stringify(msg));
};

document.getElementById("login").onclick = function () {
  data = document.getElementById("name").value;
  if (data === "") {
    alert("name can`t be empty");
    return;
  }

  var msg = {
    action: Action.Login,
    data: data,
  };

  exampleSocket.send(JSON.stringify(msg));
};

document.getElementById("logout").onclick = function () {
  document.getElementById("body").style.display = "none";

  document.getElementById("login-div").style.display = "block";
  document.getElementById("hello").innerText = "Hello,";

  var msg = {
    action: Action.Logout,
  };

  exampleSocket.send(JSON.stringify(msg));
};

exampleSocket.onmessage = function (event) {
  msg = JSON.parse(event.data);

  if (msg.error) {
    alert(msg.error);
    return;
  }

  switch (msg.action) {
    case ActionBack.NewUser:
      var element = document.createElement("li");
      element.innerText = msg.data;
      document.getElementById("users").appendChild(element);
      alert("new user added " + msg.data);
      break;

    case ActionBack.UserLeft:
      var msgNew = {
        action: Action.GetUsers,
      };

      exampleSocket.send(JSON.stringify(msgNew));
      alert("user left " + msg.data);
      break;

    case ActionBack.AllUsers:
      users = JSON.parse(msg.data);
      document.getElementById("users").innerHTML = "";
      users.forEach(function (user) {
        var element = document.createElement("li");
        element.innerText = user;
        document.getElementById("users").appendChild(element);
      });
      break;

    case ActionBack.LogedIn:
      var msg = {
        action: Action.GetUsers,
      };

      exampleSocket.send(JSON.stringify(msg));

      var msg = {
        action: Action.GetField,
      };

      exampleSocket.send(JSON.stringify(msg));
      document.getElementById("hello").innerText += " ";
      document.getElementById("hello").innerText += data;

      document.getElementById("body").style.display = "block";

      document.getElementById("login-div").style.display = "none";
      break;
    case ActionBack.Field:
      field = JSON.parse(msg.data);
      document.getElementById("field").innerHTML = " ";

      field.forEach(function (row) {
        var tr = document.createElement("tr");
        row.forEach(function (num) {
          var th = document.createElement("th");
          th.innerText = num;
          tr.appendChild(th);
        });
        document.getElementById("field").appendChild(tr);
      });
      break;
    case ActionBack.Lost:
      alert("You lost(((");
      break;
    default:
      alert("unknown action");
  }
};
