document.getElementById("connect").onclick = function () {
  data = document.getElementById("room-name").value;
  if (data === "") {
    alert("room can`t be empty");
    return;
  }
  document.location.replace("/" + data);
};
