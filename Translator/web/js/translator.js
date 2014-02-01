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
		}, function () {
			input.closest("tr").removeClass("untranslated");
		});
	});

	var pageOptions = $("form.page-options");
	pageOptions.find("input, select").change(function () {
		pageOptions.submit();
	});

	$("a.api").click(function () {
		var a = $(this);
		var href = a.attr('href');
		$.get(href, a.data(), function () {
			if (a.is(".reload")) {
				location.reload(true);
			}
		});
		return false;
	});

	$("a.reveal-my-translation").click(function () {
		$(this).closest("tr").addClass("my-translation");
	});
});