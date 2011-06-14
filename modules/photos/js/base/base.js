$(document).ready(function() {
	$(".picture a img").hover(function() {
		console.log("hovered over", this);
		pic = $(this);
		hover = $('#hover-image');
		hoverin = $('#hover-image-inner');
		hoverimg = $('#hover-image-img');
		hoverimgsmall = $('#hover-image-img-small');
		lgimg = JSON.parse(pic.attr('data-lgimg'));
		
		// Set low and high res images individually, so that it is never blank. As the high res version loads, it is displayed overtop of the low res one.
		hover.css({'display':'none'});
		hoverimgsmall.attr('src', pic.attr('src'));
		hoverimg.attr('src', lgimg.src);
		
		hover.css({'left':pic.offset().left+2, 'top':pic.offset().top+2});
		hover.width(pic.width());
		hover.height(pic.height());
		
		// Reset to image coords then animate out.
		hoverin.css({'left':0,'right':0,'bottom':0,'top':0});
		relpos = getRelPos(pic, lgimg, $('#content-inner'));
		hoverin.animate(relpos, {'queue':false});
		
		//hover.fadeIn(100);
		hover.stop().css({'display':'inline', 'opacity': 1});
	}, function() {
		// Hover out
		$('#hover-image').fadeOut();
	}).click(function() {
		
	});
	
	// Gets the relative positioning (map with top, left, right, and bottom) that the large element should use within the small element, such that it becomes the proper size, within the bounds bound. small and bound should be elements, large should be a map with width and height keys.
	function getRelPos(small, large, bound) {
		// TODO: Calculate the outside of the parent element's box, and keep the pictures within that.
		var left = (large.width-small.width())/2;
		var right = left;
		var top = (large.height-small.height())/2;
		var bottom = top;
		
		do {
			bor = bound.offset().left+bound.width();
			lor = pic.offset().left+pic.width()+right-bor;
			console.log("right", bor, lor);
			left++;
			right--;
		} while (lor > -3)
		
		do {
			bol = bound.offset().left;
			// large offset left
			lol = pic.offset().left-left-bol;
			console.log(bol, lol);
			left--;
			right++;
		} while (lol < 0)
		
		do {
			bob = bound.offset().top+bound.height();
			lob = pic.offset().top+pic.height()+bottom-bob;
			console.log(bob, lob);
			top++;
			bottom--;
		} while (lob > -3)
		
		do {
			bot = bound.offset().top;
			lot = pic.offset().top-top-bot;
			console.log(bot, lot);
			top--;
			bottom++;
		} while (lot < 0)
		
		console.log(left, right, top, bottom);
		
		return {
			'left': -left,
			'right': -right,
			'top': -top,
			'bottom': -bottom
		}
	}
	/*$('#photos-slideshow-big').click(function() {
	$('#container').css('overflow', 'visible');
		$('#photos-slideshow').animate({
			width: $(document).width(),
			height: $(window).height(),
			left: -(($(document).width()-$('#content').width())/2),
			top: -($(document).height()-$('#content').height()-$('#footer').height())
		}, "slow", function() {
			$('#photos-slideshow').css('position', 'fixed').css('top', 0).css('left', 0).css('width', '100%').css('height', '100%');
			$('#close-slideshow').css('display', 'block').click(function() {
				window.location.reload();
			});
		});
	});*/
	if (!/android|iphone|ipod|series60|symbian|windows ce|blackberry/i.test(navigator.userAgent)) {
		jQuery(function($) {
			$("a[rel^='lightbox']").slimbox({/* Put custom options here */}, null, function(el) {
				return (this == el) || ((this.rel.length > 8) && (this.rel == el.rel));
			});
		});
	}
});
