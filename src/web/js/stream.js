console.log("STREAM");

var source = new EventSource('http://localhost:5000/stream/accelerometer');

source.addEventListener('message', function(e) {
    console.log(e.data);
}, false);

source.addEventListener('open', function(e) {
    console.log("OPEN!");
}, false);

source.addEventListener('error', function(e) {
    console.log("ERROR! " + e.readyState);
    if (e.readyState == EventSource.CLOSED) {
        // Connection was closed.
    }
}, false);
