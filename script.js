function login(){
    let email, password;
    email = $("#inputEmail").val();
    password = $("#inputPassword").val();
    if (!email || !password){
        alert("Email or Password Missing");
    }else{
        waitingDialog.show("Loging In...");
        var settings = {
            "async": true,
            "crossDomain": true,
            "url": "http://localhost:5000/api/login",
            "method": "POST",
            "headers": {
              "Content-Type": "application/json",
              "Authorization": "Basic "+ btoa(email + ":" + password),
            }
          }
          $.ajax(settings).done(function (response) {
            var resp = response;
            if(!resp){
                alert("Invalid email or password");
            }else{
                sessionStorage["email"] = resp['email'];
                sessionStorage["id"] = resp['id'];
                sessionStorage["name"] = resp['name'];
                sessionStorage["token"] = resp['token'];
                $("#myCarousel").carousel('next');
                generateAddress();
            }
            waitingDialog.hide();
          });
    }
}

function generateAddress(){
    let coin, username, userid;
    coin = "KMD";
    username = sessionStorage["username"];
    userid = sessionStorage["userid"];
    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "http://localhost:5000/api/getnewaddress",
        "method": "POST",
        "headers": {
          "X-Access-Token": sessionStorage.token,
          "Content-Type": "application/json"
        },
        "processData": false,
        "data": JSON.stringify({"coin" : "KMD" , "username" : sessionStorage.name, "userid" : sessionStorage.id})
      }
      
      $.ajax(settings).done(function (response) {
        $("#gaddress").val(response.address);
      });
}

function generateTx(){
    let address, amount;
    address = $("#raddress").val();
    amount = $("#ramount").val() * 100000000;
    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "http://localhost:5000/api/createrawtx",
        "method": "POST",
        "headers": {
          "X-Access-Token": sessionStorage.token,
          "Content-Type": "application/json"
        },
        "processData": false,
        "data": JSON.stringify({"coin" : "KMD" , "address" : address, "amount" : amount})
      }
      $.ajax(settings).done(function (response) {
        console.log(response);
        if(response['success'] == false){
          alert(response['message']);
        }else if(response['success'] == true){
          broadcastTx(response["raw_hash"]);
        }
      });
}

function broadcastTx(hash){
  var settings = {
    "async": true,
    "crossDomain": true,
    "url": "http://localhost:5000/api/broadcasttx",
    "method": "POST",
    "headers": {
      "X-Access-Token": sessionStorage.token,
      "Content-Type": "application/json"
    },
    "processData": false,
    "data": JSON.stringify({"hash" : hash})
  }
  $.ajax(settings).done(function (response) {
    console.log(response);
    if(response.success == true){
      $(".success_message").html(response["message"]);
      $(".success_tx a").html(response["txid"]);
      let link = "https://kmdexplorer.io/tx/"+response["txid"];
      $(".success_tx a").attr("href", link);
      $("#myModal").modal();
    }else if (response.sucesss == false){
      alert(response);
    }
  }).fail(function(error){
    console.log(error);
  });
}