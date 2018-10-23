"use strict";
//Input
var hostURL = document.getElementById('hostURL');
var fileHash = document.getElementById('fileHash');
var form = document.getElementById('fileForm');
var file = document.getElementById('file');
//Buttons
var pinButton = document.getElementById('pin');
var unpinButton = document.getElementById('unpin');
var fetchButton = document.getElementById('fetch');
var storeButton = document.getElementById('store');
//Response field
var responseField = document.getElementById('response');
pinButton.addEventListener("click", function (e) { return pin(); });
unpinButton.addEventListener("click", function (e) { return unpin(); });
fetchButton.addEventListener("click", function (e) { return fetch(); });
storeButton.addEventListener("click", function (e) { return store(); });
function pin() {
    var requestURL = "http://" + hostURL.value + '/pin?hash=' + fileHash.value;
    console.log(requestURL);
    request(requestURL, 'patch', null, function cb(status, res) {
        if (status != 200) {
            responseField.value = "FAILED";
        }
        else {
            responseField.value = res;
        }
    });
}
function unpin() {
    var requestURL = "http://" + hostURL.value + '/unpin?hash=' + fileHash.value;
    console.log(requestURL);
    request(requestURL, 'patch', null, function cb(status, res) {
        if (status != 200) {
            responseField.value = "FAILED";
        }
        else {
            responseField.value = res;
        }
    });
}
function fetch() {
    var requestURL = "http://" + hostURL.value + '/fetch?hash=' + fileHash.value;
    window.open(requestURL);
   /* request(requestURL, 'get', null, function cb(status, res) {
        if (status != 200) {
            responseField.value = "FAILED";
        }
        else if (res == null) {
            console.log("No Response");
        }
        else {
            console.log(res);
            responseField.value = res;
        }
    });*/
}
function store() {
    var requestURL = "http://" + hostURL.value + '/store';
    console.log(requestURL);
    var form1 = document.getElementById('file');
    var formData = new FormData(form1);
    request(requestURL, 'post', formData, function cb(status, res) {
        if (status != 200) {
            responseField.value = "FAILED";
        }
        else {
            responseField.value = res;
        }
    });
}
function request(url, method, formData, callback) {
    var xhr = new XMLHttpRequest();
    xhr.timeout = 30000;
    xhr.onreadystatechange = function (e) {
        if (xhr.readyState == 4) {
            if (xhr.status === 200) {
                callback(xhr.status, xhr.response);
            }
            else {
                callback(xhr.status, xhr.response);
            }
        }
    };
    xhr.ontimeout = function () {
        console.log("Timeout...");
    };
    xhr.open(method, url, true);
    if (formData != null) {
        xhr.send(formData);
    }
    else {
        xhr.send();
    }
}
