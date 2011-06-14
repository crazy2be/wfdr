var uploadErrorsHelp = {
	500: 'Make sure you have selected an album from the list above.',
	404: 'You likely don\'t have permission to upload to this album.'
};
if (!window.console) window.console = {log: function(){}};

// Uploader code.
(function($) {
	var options = {
		uploaded: function(files, index, response) {
			return options.defaultUploaded(files, index, response);
		},
		defaultUploaded: function(files, index, response) {
			console.log(index, response);
			file = {};
			file.elm = $('#'+files[index].id);
			$('#photos-upload-progress').get(0).value = index+1;
			if (response == '1') {
				file.elm.addClass('dropped-image-complete');
			} else {
				file.elm.addClass('dropped-image-error');
				// Split into error number and textual description.
				resp = response.split('\n', 2);
				file.elm.find('.error-message').html('<h3>' + response + '</h3>' + uploadErrorsHelp[resp[0]]);
				file.elm.find('.dropped-image-actions').fadeIn();
			}
		},
		progress: function(files, index, done, total) {
			return options.defaultProgress(files, index, done, total);
		},
		defaultProgress: function(files, index, done, total) {
			var file = $('#'+files[index].id+'-progress').get(0);
			var totalProgress = $('#photos-upload-progress').get(0);
			totalProgress.max = files.length;
			totalProgress.value = index+(done/total);
			totalProgress.nextSibling.innerHTML = 'Uploading ' +  files[index].file.name + '...';
			file.value = done;
			file.max = total;
			console.log(file);
			console.log(totalProgress);
		},
		upload: function(files, index) {
			return options.defaultUpload(files, index);
		},
		defaultUpload: function(files, index) {
			$('#photos-upload-progress-outer').slideDown();
			// Actual uploading
			var url = '/picasa/upload';
			var xhr = new XMLHttpRequest();
			//console.log(files, index);
			var file = files[index];
			xhr.file = file.file;
			
			xhr.addEventListener('progress', function(e) {
				var done = e.position || e.loaded;
				var total = e.totalSize || e.total;
				console.log('xhr progress: ' + (Math.floor(done/total*1000)/10) + '%');
				options.progress(files, index, done, total);
			}, false);
				
				if ( xhr.upload ) {
					xhr.upload.onprogress = function(e) {
						var done = e.position || e.loaded;
						var total = e.totalSize || e.total;
						console.log('xhr.upload progress: ' + done + ' / ' + total + ' = ' + (Math.floor(done/total*1000)/10) + '%');
						options.progress(files, index, done, total);
					};
				}
				
				xhr.onreadystatechange = function(e) {
					//console.log(['XHR onreadystatechange', e]);
					if ( 4 == this.readyState ) {
						console.log(['xhr upload complete', e]);
						var response = this.responseText;
						// Allows only one upload at a time.
						if (files.length > index+1) {
							options.upload(files, index+1);
						}
						options.uploaded(files, index, response);
					}
				};
				
				xhr.open('post', url, true);
				xhr.setRequestHeader('X-Album-Name', $('#album-name').val());
				xhr.setRequestHeader('Slug', file.file.name);
				xhr.setRequestHeader('Content-Type', file.file.type);
				xhr.send(file.file);
		},
		dropped: function(event) {
			return options.defaultDropped(event);
		},
		defaultDropped: function(event) {
			$('#photos-drop').removeClass('drop-over');
			$('#photos-drop h1').remove();
			var dt = event.dataTransfer;
			var origFiles = dt.files;
			console.log('# of files dropped: ', origFiles.length);
			var ids = [];
			var files = [];
			// TODO: A text-based fallback for > 100 files.
			for (var i = 0; i < origFiles.length; i++) {
				var file = origFiles[i]
				if (file.type == 'image/jpeg') {
					files[files.length] = {};
					files[files.length-1].file = file;
					var id = parseInt(Math.random()*1000000);
					files[files.length-1].id = id;
				} else {
					// TODO: A user-visible error.
					console.log('Only JPEG supported');
				}
			}
			//console.log(files);
			// WARNING: This could cause a race condition in the unlikely case that the uploader function chews through pictures faster than the show function.
			options.show(files, 0);
			// The user has not selected an album
			if ($('#album-name').val().length < 3) {
				$('#tip').fadeOut(10).html('<b>Error:</b> Please select an album from the list above!').addClass('error').fadeIn(1000);
				$('#album-name').change(function() {
					$('#tip').fadeOut(10).html('Your photos are now uploading.').removeClass('error').fadeIn(1000);
					options.upload(files, 0);
				});
			} else {
				options.upload(files, 0);
			}
			// Stop the browser from opening the image.
			event.stopPropagation();
			event.preventDefault();
		},
		show: function(files, index) {
			return options.defaultShow(files, index);
		},
		defaultShow: function(files, index) {
			console.log('show() called');
			var file = files[index];
			var id = file.id;
			droppedImage = $(
				'<div class="dropped-image" id="'+id+'">\
					<div class="dropped-image-inner clearfix">\
						<div class="dropped-image-actions">\
							<img class="close"  alt="Close/Remove from list" src="/img/close.png">\
							<img class="retry" alt="Retry uploading this image" src="/img/retry.png">\
						</div>\
						<progress class="photo-upload-progress" id="'+id+'-progress" value="0"></progress>\
						<div class="dropped-image-description">\
							<div class="dropped-image-name">'
								+file.file.name+
							'</div>\
							<div class="dropped-image-size">'
							// Multiply inside, divide outside allows us to get two decimal places of precision.
							+(Math.round(file.file.size/1000/1000*100)/100)+'MB'+
							'</div>\
						</div>\
						<div class="error-message"></div>\
					</div>\
				</div>');
			droppedImage.data('file', file);
			$('#photos-drop').append(droppedImage);
			//img.classList.add('dropped-image');
			//img.id = id;
			//document.getElementById(id).appendChild(img);
			//(function(aImg) { return function(e) { aImg.src = e.target.result; }; })(img);

			//var reader = new FileReader();
			//reader.onload = function(event) {
				// For some reason, img doesn't seem to have the correct value, but id does. Likely a browser bug.
				//$('#'+id+' .dropped-image-thumbnail').attr('src', event.target.result);
				//document.getElementById(id).firstChild.firstChild.nextSibling.nextSibling.nextSibling.src = event.target.result; 
				//console.log(id);
				//console.log("in show() function: ", index, id);
				// We wait for the first image to load before attempting to load the next one. Anything else just causes massive lag on the user's system.
				//if (index+1 < files.length) {
					//options.show(files, index+1);
				//}
			//}
			
			//console.log(thing);
			//reader.readAsDataURL(file.file);
			if (index+1 < files.length) {
				options.show(files, index+1);
			}
		},
		init : function() {
			// Is the API supported?
			if (!Modernizr.draganddrop) {
				return
			}
			// Initialization.
			$('#photos-drop').get(0).addEventListener('drop', options.dropped, false);
			$('#photos-drop').get(0).addEventListener('dragover', function(event) {
				$('#photos-drop').addClass('drop-over');
				event.stopPropagation();
				event.preventDefault();
			}, false);
			$('#photos-drop').get(0).addEventListener('dragenter', function(event) {
				$('#photos-drop').addClass('drop-over');
				event.stopPropagation();
				event.preventDefault();
			}, false);
			$('#photos-drop').get(0).addEventListener('dragleave', function(event) {
				$('#photos-drop').removeClass('drop-over');
				event.stopPropagation();
				event.preventDefault();
			}, false);
			$('.close').live('click', function(event) {
				//TODO: Should this remove it from the files array as well?
				$(this).closest('.dropped-image').remove();
			});
			$('.retry').live('click', function(event) {
				file = $(this).closest('.dropped-image').data('file');
				console.log(file);
				files = [file];
				options.upload(files, 0);
			});
		},
		// Actual "options". All the ones above are overridable methods.
		path: '/picasa/upload'
	};
	var file = {
		id: 0,
		fileObj: 0,
		init: function(id) {
			
		}
	};
	$.fn.draguploader = function( method ) {
		// Method calling logic
		if ( options[method] ) {
			return options[ method ].apply( this, Array.prototype.slice.call( arguments, 1 ));
		} else if ( typeof method === 'object' || ! method ) {
			return options.init.apply( this, arguments );
		} else {
			$.error( 'Method ' +  method + ' does not exist on jQuery.draguploader' );
		}
	}
})(jQuery);
$(document).ready(function() {
	if (navigator.appName == 'Microsoft Internet Explorer') {
		// For now until they add support.
		$('html').addClass('no-draganddrop').removeClass('draganddrop');
		Modernizr.draganddrop = false;
	}
	$('#photos-drop').draguploader();
});