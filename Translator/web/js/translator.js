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
			input.closest("tr.entry").removeClass("untranslated");
			$("#saved-notice").stop(true).show().fadeTo(0, 1.0).fadeOut(2500);
		});
	}).focus(function () {
		$(this).closest('td.translation-parts').find('.my-translation-arrow-score').addClass('focus');
	}).blur(function () {
		$(this).closest('td.translation-parts').find('.my-translation-arrow-score').removeClass('focus');
	});

	$("body.translate a.vote").click(function () {
		var a = $(this);
		var up = a.is('.vote-up');
		var active = a.is('.active');

		var original = a.closest('tr.entry').find('input.entry-original').first().val();
		var partOf = a.closest('tr.entry').find('input.entry-partof').first().val();
		var translation = a.closest('.other-translation').find('label.part').first().text();

		$.get('/api/vote', {
			language: $("#current-language").val(),
			original: original,
			partOf: partOf,
			translation: translation,
			up: (up && !active),
			down: (!up && !active)
		});

		if (up) {
			a.closest("td").find("a.vote-up").removeClass('active btn-success');
		}

		var inverse = a.closest('.btn-group').find('a.vote-'+(up ? 'down' : 'up'));
		inverse.removeClass('active btn-success btn-danger');
		
		if (active) {
			a.removeClass('active btn-success btn-danger');
		} else {
			a.addClass('active btn-'+(up ? 'success' : 'danger'));
		}
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
		$(this).closest("tr").addClass("with-translation").find("p.my-translation input").first().focus();
	});
});