/*
 * jQuery Color Animations
 * Copyright 2007 John Resig
 * Released under the MIT and GPL licenses.
 */

(function(jQuery){

    // We override the animation for all of these color styles
    jQuery.each(['backgroundColor', 'borderBottomColor', 'borderLeftColor', 'borderRightColor', 'borderTopColor', 'color', 'outlineColor'], function(i,attr){
        jQuery.fx.step[attr] = function(fx){
            if ( !fx.colorInit ) {
                fx.start = getColor( fx.elem, attr );
                fx.end = getRGB( fx.end );
                fx.colorInit = true;
            }

            fx.elem.style[attr] = "rgba(" + [
                Math.max(Math.min( parseInt((fx.pos * (fx.end[0] - fx.start[0])) + fx.start[0]), 255), 0),
                Math.max(Math.min( parseInt((fx.pos * (fx.end[1] - fx.start[1])) + fx.start[1]), 255), 0),
                Math.max(Math.min( parseInt((fx.pos * (fx.end[2] - fx.start[2])) + fx.start[2]), 255), 0),
                Math.max(Math.min( parseFloat((fx.pos * (fx.end[3] - fx.start[3])) + fx.start[3]), 1), 0)
            ].join(",") + ")";
        }
    });

    // Color Conversion functions from highlightFade
    // By Blair Mitchelmore
    // http://jquery.offput.ca/highlightFade/

    // Parse strings looking for color tuples [255,255,255]
    function getRGB(color) {
        var result;

        // Check if we're already dealing with an array of colors
        if ( color && color.constructor == Array && color.length == 3 )
            return color;

								// Look for rgba(num,num,num,num)
        if (result = /rgba\(\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([-+]?[0-9]*\.?[0-9]+)\s*\)/.exec(color))
            return [parseInt(result[1]), parseInt(result[2]), parseInt(result[3]), parseFloat(result[4])];

        // Look for rgb(num,num,num)
        if (result = /rgb\(\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*,\s*([0-9]{1,3})\s*\)/.exec(color))
            return [parseInt(result[1]), parseInt(result[2]), parseInt(result[3]), 1];

        // Look for rgb(num%,num%,num%)
        if (result = /rgb\(\s*([0-9]+(?:\.[0-9]+)?)\%\s*,\s*([0-9]+(?:\.[0-9]+)?)\%\s*,\s*([0-9]+(?:\.[0-9]+)?)\%\s*\)/.exec(color))
            return [parseFloat(result[1])*2.55, parseFloat(result[2])*2.55, parseFloat(result[3])*2.55, 1];

        // Look for #a0b1c2
        if (result = /#([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})/.exec(color))
            return [parseInt(result[1],16), parseInt(result[2],16), parseInt(result[3],16), 1];

        // Look for #fff
        if (result = /#([a-fA-F0-9])([a-fA-F0-9])([a-fA-F0-9])/.exec(color))
            return [parseInt(result[1]+result[1],16), parseInt(result[2]+result[2],16), parseInt(result[3]+result[3],16), 1];

        // Look for rgba(0, 0, 0, 0) == transparent in Safari 3
        if (result = /rgba\(0, 0, 0, 0\)/.exec(color))
            return colors['transparent'];

        // Otherwise, we're most likely dealing with a named color
        return colors[jQuery.trim(color).toLowerCase()];
    }

    function getColor(elem, attr) {
        var color;

        do {
            color = jQuery.curCSS(elem, attr);

            // Keep going until we find an element that has color, or we hit the body
            if ( color != '' && color != 'transparent' || jQuery.nodeName(elem, "body") )
                break;

            attr = "backgroundColor";
        } while ( elem = elem.parentNode );

        return getRGB(color);
    };

    // Some named colors to work with
    // From Interface by Stefan Petre
    // http://interface.eyecon.ro/

    var colors = {
        aqua:[0,255,255,1],
        azure:[240,255,255,1],
        beige:[245,245,220,1],
        black:[0,0,0,1],
        blue:[0,0,255,1],
        brown:[165,42,42,1],
        cyan:[0,255,255,1],
        darkblue:[0,0,139,1],
        darkcyan:[0,139,139,1],
        darkgrey:[169,169,169,1],
        darkgreen:[0,100,0,1],
        darkkhaki:[189,183,107,1],
        darkmagenta:[139,0,139,1],
        darkolivegreen:[85,107,47,1],
        darkorange:[255,140,0,1],
        darkorchid:[153,50,204,1],
        darkred:[139,0,0,1],
        darksalmon:[233,150,122,1],
        darkviolet:[148,0,211,1],
        fuchsia:[255,0,255,1],
        gold:[255,215,0,1],
        green:[0,128,0,1],
        indigo:[75,0,130,1],
        khaki:[240,230,140,1],
        lightblue:[173,216,230,1],
        lightcyan:[224,255,255,1],
        lightgreen:[144,238,144,1],
        lightgrey:[211,211,211,1],
        lightpink:[255,182,193,1],
        lightyellow:[255,255,224,1],
        lime:[0,255,0,1],
        magenta:[255,0,255,1],
        maroon:[128,0,0,1],
        navy:[0,0,128,1],
        olive:[128,128,0,1],
        orange:[255,165,0,1],
        pink:[255,192,203,1],
        purple:[128,0,128,1],
        violet:[128,0,128,1],
        red:[255,0,0,1],
        silver:[192,192,192,1],
        white:[255,255,255,1],
        yellow:[255,255,0,1],
        transparent: [255,255,255, 0]
    };

})(jQuery);
