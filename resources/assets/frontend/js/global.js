var Octocv = Octocv||{};

Octocv.FacebookLogin = {
	initialized: false,
	setup: function(){
		if(!this.initialized) {
			this.initialized = true;
			window.fbAsyncInit = function() {
				FB.init({
					appId: "221184025235764",
					cookie: true,
					xfbml: true,
					version: "v2.10"
				})
			};
			
			try {
				e=document,i="script",n="facebook-jssdk",o=e.getElementsByTagName(i)[0],e.getElementById(n)||((t=e.createElement(i)).id=n,t.src="https://connect.facebook.net/en_US/sdk.js",o.parentNode.insertBefore(t,o))
			} catch(c){}
			
			var e,i,n,t,o;
			$("#auth_facebook").on("click",function(i){
				i.preventDefault(),FB.login(function(e){
					e.authResponse&&(window.location=i.currentTarget.href)
				},{scope:"public_profile,email,user_friends"})
			})
		}
	}
};

(function($, window, undefined) {
	var pluginName = 'form-required';

	function Plugin(element, options) {
		this.element = $(element);
		this.options = $.extend({}, $.fn[pluginName].defaults, options);
		this.init();
	}

	Plugin.prototype = {
		init: function() {
			var fields = this.element.data("form-required");
			this.element.on("submit", function (e) {
				var res  = true;
				$.each(fields.split("|"), function(k, v){
					if ($("input[name=" + v + "]").val() == "") {
						$("label[for=" + v + "]").css({color: "#f66"})
						$("input[name=" + v + "]").css({border: "solid 1px #f66"})
						res = false;
					}
				});
				return res
			})
		}
	};

	$.fn[pluginName] = function(options, params) {
		return this.each(function() {
			var instance = $.data(this, pluginName);
			if (!instance) {
				$.data(this, pluginName, new Plugin(this, options));
			} else if (instance[options]) {
				instance[options](params);
			} else {
				window.console && console.log(options ? options + ' method is not exists in ' + pluginName : pluginName + ' plugin has been initialized');
			}
		});
	};

	$.fn[pluginName].defaults = {
		option: 'value'
	};

	$(function() {
		$('[data-' + pluginName + ']')[pluginName]();
	});

}(jQuery, window));