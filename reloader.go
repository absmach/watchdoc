package main

const reloaderScript = `
<script>
(function() {
	const ws = new WebSocket('ws://' + window.location.host + '/ws');
	ws.onmessage = function(event) {
		if (event.data === 'reload') {
			console.log('Files changed, reloading...');
			window.location.reload();
		}
	};
	ws.onclose = function() {
		console.log('Dev server disconnected, retrying...');
		setTimeout(() => window.location.reload(), 1000);
	};
})();
</script>`
