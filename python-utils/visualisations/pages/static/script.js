var searchString = window.location.search
var params = new URLSearchParams(searchString);
var pathname = params.get("filename");
var globaldict;
// Get the modal

$(document).ready(function(){
	var modal = document.getElementById("reuse-modal");
	$.post( "/serve_reuse", {"filename": pathname}, function( data ) {
		console.log(data)
		$("#main").html(data['segments']);
		globaldict = data;
	});

	$(".close").on("click", function() {
		modal.style.display = "none";
	});

	$("#main").on("click", ".reused", function() {
		console.log('clicked on reused');
		modal.style.display = "block";
	});
});
