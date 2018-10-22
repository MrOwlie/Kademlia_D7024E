
//Input
var hostURL = <HTMLInputElement>document.getElementById('hostURL')
var fileHash = <HTMLInputElement>document.getElementById('fileHash')
var file = <HTMLInputElement>document.getElementById('file')

//Buttons
var pinButton = <HTMLInputElement>document.getElementById('pin')
var unpinButton = <HTMLInputElement>document.getElementById('unpin')
var fetchButton = <HTMLInputElement>document.getElementById('fetch')
var storeButton = <HTMLInputElement>document.getElementById('store')

//Response field
var responseField = <HTMLInputElement>document.getElementById('response')

pinButton.addEventListener("click", (e:Event) => pin())
unpinButton.addEventListener("click", (e:Event) => unpin())
fetchButton.addEventListener("click", (e:Event) => fetch())
storeButton.addEventListener("click", (e:Event) => store())


function pin() {

  let requestURL =  "http://" + hostURL.value + '/pin?hash=' + fileHash.value
  console.log(requestURL)
  request(requestURL, 'patch', null, function cb(status, res){
    if(status != 200){
      responseField.value = "FAILED"
    } else {
      responseField.value = res.body
    }
  })
}

function unpin() {
  let requestURL =  "http://" + hostURL.value + '/unpin?hash=' + fileHash.value
  console.log(requestURL)
  request(requestURL, 'patch', null, function cb(status, res){
    if(status != 200){
      responseField.value = "FAILED"
    } else {
      responseField.value = res.body
    }
  })
}

function fetch() {
  var requestURL = "http://" + hostURL.value + '/store';
  console.log(requestURL)
  request(requestURL, 'get', null,  function cb(status, res){
    if(status != 200){
      responseField.value = "FAILED"
    } else if(res == null){
      console.log("No Response")
    } else {
      console.log(res)
      responseField.value = res.body
    }
  })
}

function store() {
  let requestURL =  "http://" + hostURL.value + '/pin?hash=' + fileHash.value
  console.log(requestURL)
  request(requestURL, 'post', file.value, function cb(status, res){
    if(status != 200){
      responseField.value = "FAILED"
    } else {
      responseField.value = res.body
    }
  })
}

function request(url, method, file, callback) {
  let xhr = new XMLHttpRequest();
  xhr.timeout = 2000;
  xhr.onreadystatechange = function(e) {
    if (xhr.readyState == 4) {
      if (xhr.status === 200) {
        callback(xhr.status, xhr.response)
      } else {
        callback(xhr.status, xhr.response)
      }
    }
  }
  xhr.ontimeout = function () {
    console.log("Timeout...")
  }
  xhr.open(method, url, true)
  if(file != null){
    xhr.setRequestHeader("Content-Type", "multipart/form-data")
    xhr.send(file)
  }else {
    xhr.send();
  }
}