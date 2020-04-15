var searchString = window.location.search
var params = new URLSearchParams(searchString);
var pathname = params.get("filename");
$.post( "/serve_reuse", {"filename": pathname}, function( data ) {
	console.log(data)
	$("#main").html(data['segments']);
});