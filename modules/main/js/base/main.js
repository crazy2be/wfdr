// main.js: JavaScript functions required for the home page of the website.
//

$(document).ready(function() {
	console.log("First document ready");
	$("#google-search input").focus();
	/*$("h2").wrapInner("<span>")
	$("h2 br")
		.before("<span class='spacer'>&nbsp;</span>")
		.after("<span class='spacer'>&nbsp;</span>");*/
	$('.anythingSlider').anythingSlider({
        easing: "swing",                // Anything other than "linear" or "swing" requires the easing plugin
        autoPlay: true,                 // This turns off the entire FUNCTIONALY, not just if it starts running or not
        startStopped: false,            // If autoPlay is on, this can force it to start stopped
        delay: 6000,                    // How long between slide transitions in AutoPlay mode
        animationTime: 600,             // How long the slide transition takes
        hashTags: true,                 // Should links change the hashtag in the URL?
        buildNavigation: false,          // If true, builds and list of anchor links to link to each slide
        pauseOnHover: true,             // If true, and autoPlay is enabled, the show will pause on hover
        startText: "Start",             // Start text
        stopText: "Stop",               // Stop text
        navigationFormatter: null       // Details at the top of the file on this use (advanced use)
	});
});

searchScript = false;

function handleSearch(results) {
	console.log(results)
	res = results.responseData;
	rescont = $('#google-search-results');
	num = res.cursor.estimatedResultCount;
	num = (num != undefined)?num:0;
	rescont.html('<p>Press enter to see all results on google. (About '+num+' results)</p>');
	for (i = 0; i < res.results.length; i++) {
		result = res.results[i];
		console.log(result);
		rescont.html(rescont.html()+getResultHTML(result.url, result.title, result.content));
	}
}

function getResultHTML(url, title, content) {
	return '<p><a href="'+url+'"><h3>'+title+'</h3></a>'+content+'<br /></p>';
}

function search(query) {
	if (searchScript !== false) {
		document.body.removeChild(searchScript);
	}
	searchScript = document.createElement("script");
	searchScript.type = "text/javascript";
	searchScript.src = 'http://ajax.googleapis.com/ajax/services/search/web?v=1.0&q='+encodeURIComponent(query)+'&callback=handleSearch&nocache='+Math.random();
	document.body.appendChild(searchScript);
}

$(document).ready(function() {
	console.log("Document ready");
	previousValue = ''
	animating = false;
	stuffHidden = false;
	$('#google-search-bar').keyup(function(e) {
		currentValue = this.value;
		function animations() {
			if (currentValue == '' && stuffHidden == true) {
				if (animating) return;
				animating = true
				$('#bottom-buttons,.anythingSlider').slideDown(400, function() {
					animating = false;
					stuffHidden = false
				});
				$('#google-search-results').slideUp(400);
			}
			if ((previousValue == '' || stuffHidden == false) && currentValue != '') {
				if (animating) return;
				animating = true
				dummyresults = '';
				for (i = 0; i < 4; i++) {
					dummyresults += getResultHTML("#", "Loading results...", "Loading the results. Please wait. While you are waiting, you might want to 1) check your internet is working, 2) flick your fingers together, and 3) ooh and aww about how pretty these animations are.");
				}
				$('#google-search-results').html(dummyresults);
				$('#bottom-buttons,.anythingSlider').slideUp(400, function() {
					animating = false;
					stuffHidden = true;
				});
				$('#google-search-results').slideDown(400);
			}
		};
		animations();
		search(this.value);
		previousValue = this.value;
	});
});

