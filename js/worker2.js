var i = -1;
var fdet = null;
self.addEventListener('message', function(e) {

  //if (e == null || e.data == null) self.close();
  if (fdet == null) {
    fdet = {};
    fdet.i = e.data.i;
    fdet.f = e.data.f;
    fdet.fsz = e.data.f.size;
    fdet.t = e.data.t;
    fdet.rsz = e.data.rsz;
    fdet.rpos = fdet.i * fdet.rsz;
    fdet.jump = fdet.rsz * fdet.t;
    fdet.active = false;
    i = e.data.i;
    self.postMessage({ i: i, val : 'ready' });
  }
  else {
    if (e.data === 'terminate') {
      wrapUp();
    } else {
      if(fdet.rpos < fdet.fsz) {
        self.postMessage({ i: i, val : 'sending'});
        readData();
      } else {
        self.postMessage({ i: i, val : 'done'});
      }
    }
  }

}, false);

self.wrapUp = function() {
    if (fdet.active == true) {
      setTimeout(self.wrapUp, 1000);
    } else {
      self.close();
    }
}

self.updateProgress = function(e) {
    if (e.lengthComputable) {
      var value = (e.loaded / e.total) * 100;
      self.postMessage({ i: i, val : value });
      console.log(value);
    }
}

self.upload = function(index, name, blobOrFile, offset, ijs) {

  var xhr = new XMLHttpRequest();
  xhr.open('POST', '/wrdr/uploadChunk', true);
  xhr.setRequestHeader('X-File-Offset', '' + offset);
  xhr.setRequestHeader('X-File-Name', '' + name);
  xhr.onload = function(e) { 
    fdet.active = false;
    self.postMessage(ijs);
  }

  // Listen to the upload progress.
  //xhr.addEventListener('progress', updateProgress, false);

  xhr.send(blobOrFile);
}
self.uploadSync = function(index, name, blobOrFile, offset, ijs) {

  var xhr = new XMLHttpRequest();
  xhr.open('POST', '/wrdr/uploadChunk', false);
  xhr.setRequestHeader('X-File-Offset', '' + offset);
  xhr.setRequestHeader('X-File-Name', '' + name);
  xhr.send(blobOrFile);
  if (xhr.status == 200 || xhr.status == 302) {
    fdet.active = false;
    self.postMessage(ijs);
    return 0;
  }
  return 1;

}

//self.readData = function(index, file, pt, rbytes) {
self.readData = function() {
  while(fdet.rpos < fdet.fsz) {
  //if(fdet.rpos < fdet.fsz) {
    var epos = fdet.rpos + fdet.rsz;
    if (epos > fdet.fsz) epos = fdet.fsz;
    var bl = fdet.f.slice(fdet.rpos, epos);
    var bys = epos - fdet.rpos;
    var ijs = { 'i' : fdet.i, 'rpos': fdet.rpos, 'val': 'progress', 'bytesSent': bys }; 
    fdet.active = true;
    //upload(fdet.i, fdet.f.name, bl, fdet.rpos, ijs);
    uploadSync(fdet.i, fdet.f.name, bl, fdet.rpos, ijs);
    fdet.rpos = fdet.rpos + fdet.jump;
  }
  //comment below if using an if
  self.postMessage({i: fdet.i, val: "done"});

}

self.readDataDebug = function(index, file, rbytes) {
  var rpos = index*rbytes;
  if (rpos > file.size) return;
  var bl = file.slice(rpos, rpos + rbytes);
  var ijs2 = { 'i' : index, 'rpos': rpos, 'bl' : bl }; 
  //postMessage(ijs2);
  var rdr = new FileReaderSync();
  var txt = '' + rdr.readAsArrayBuffer(bl);
  txt = null;
  var ijs = { 'i' : index, 'rpos': rpos, 'rbytes' : rpos + rbytes, 'bl' :bl }; 
  postMessage('' + rpos + ':' + file.size + ':' + txt + ijs);

}

