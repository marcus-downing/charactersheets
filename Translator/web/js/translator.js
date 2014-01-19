jQuery(function ($) {
	$("body.translate input").change(function () {
		var input = $(this);
		var name = input.attr('name');
		var original = $("#"+name+"-original").val();
		var partOf = $("#"+name+"-partof").val();
		var translation = input.val();
		$.get('/api/translate', {
			language: $("#current-language").val(),
			original: original,
			partOf: partOf,
			translation: translation
		});
	});

	var pageOptions = $("form.page-options");
	pageOptions.find("input, select").change(function () {
		pageOptions.submit();
	});
});