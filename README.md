webworker-upload
================

File upload using web workers to parallelize the upload

|- freader.html (main html file)
|- client
|-- A test client written in go (can use in lieu of browser)
|- srv
|-- The go web server that handles the upload
|- js
|-- the web worker files that perform the browser work
|- css
|-- any css assets
