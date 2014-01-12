jQuery(function ($) {
	$("body.translate input").change(function () {
		var input = $(this);
		var name = input.attr('name');
		var translation = input.val();
		$.get('/api/translate', {
			language: $("#current-language").val(),
			name: name,
			translation: translation
		});
	});
});