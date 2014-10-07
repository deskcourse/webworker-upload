self.addEventListener('message', function(e) {
  var file = e.data;
  var buffers = [];

  var nr = 8;
  var rlen = Math.ceil(file.size / nr );

  for(var i=0; i<nr; i++) {
    var w = new Worker('worker2.js');
    w.onmessage = function(e) {
	    postMessage(e.data);
    }
    w.postMessage({ 'i' : i, 'f' : files[0], 's' : rlen});
    //readData(i, files[0], rlen);
  }

  // Read each file synchronously as an ArrayBuffer and
  // stash it in a global array to return to the main app.
  // 
  /*
  [].forEach.call(files, function(file) {
    var reader = new FileReaderSync();
    buffers.push(reader.readAsArrayBuffer(file));
  });

  postMessage(buffers); */
}, false);

self.readData = function(index, file, rbytes) {
  var rpos = index*rbytes;
  if (rpos > file.size) return;
  var bl = file.slice(rpos, rpos + rbytes);
  var rdr = new FileReaderSync();
  var txt = '' + rdr.readAsText(bl);
  var ijs = { 'rpos': rpos, 'rbytes' : rpos + rbytes, 'txt' : txt }; 
  postMessage(ijs);
  //postMessage('' + rpos + ':' + file.size + ':' + txt + ijs);

}

