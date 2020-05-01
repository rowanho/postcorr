var searchString = window.location.search
var params = new URLSearchParams(searchString);
var pathname = params.get("filename");
// Get the modal

$(document).ready(function(){
	var globaldict;
	var modal = document.getElementById("reuse-modal");
	if (pathname) {
		$.post( "/serve_reuse", {"filename": pathname}, function( data ) {
			console.log(data)
			$("#main").html(data['segments']);
			globaldict = data;
		});
	} else {
		console.log('links')
		$.post( "/serve_links",{}, function( data ) {
			$("#main").html(data['links']);
		});
	}

	$(".close").on("click", function() {
		modal.style.display = "none";
	});

	// When the user clicks anywhere outside of the modal, close it
	window.onclick = function(event) {
  	if (event.target == modal) {
    	modal.style.display = "none";
  	}
	}

	$("#main").on("click", ".reused", function() {
		console.log(globaldict['reuse_map']);
		content = globaldict['reuse_map'][$(this).attr("uid")];
		$(".modal-table").html(content)
		modal.style.display = "block";
	});
});
