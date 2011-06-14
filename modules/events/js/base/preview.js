// events-preview.js
//
// Works with live preview-related functions for the event editor page.
//

// Prepares the live preview for use by updating the appropriate fields.
//
function prepareLivePreview() {
    // Event name
    $('#event-title').html($('#title').keyup(updateEventName).val());
		console.log($('#event-title').val());

    // Description
    $('#event-desc').html($('#desc').keyup(updateEventDesc).val());
		
    // Importance
    $('#importance-dd').change(updateEventImportance);
    updateEventImportance();
	
    // TODO: Image (async upload), date
}

// Updates the title of the event preview.
//
function updateEventName()
{
    $('#event-title').html($('#title').val());

    // TODO: Only do this if the user has not edited the id manually.
    $('#id').val(getURLFromName($('#title').val()));
}

// Automatically generates the URL of the event's supporting page based on the
// title of the event.
//
// Parameter: name
//   The name to use for the automatically generated URL.
//
function getURLFromName(name) {
	lower = name.toLowerCase();

	// Replace all whitespace characters and groups thereof with dashes.
	nospaces = lower.replace(/\s+/gim, "-");
	return encodeURIComponent(nospaces);
}

// Updates the importance of the event preview.
//
function updateEventImportance()
{
	importance = $('#importance-dd').val();
	oldImportance = (importance == 1)?2:1;
	console.log(importance, oldImportance);
	$('#event-base').removeClass('importance-'+oldImportance).addClass('importance-'+importance);
}

// Updates the event description.
//
function updateEventDesc()
{
    $('#event-desc').html($('#desc').val());
}
