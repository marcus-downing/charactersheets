$(function() {
  $("html, body").addClass("postload");

  $("a.lightbox").click(function () {
    var id = $(this).attr('rel');
    var lightbox = $(id);
    if (lightbox) {
      lightbox.fadeIn("fast");
      return false;
    }
    return true;
  });

  $("div.lightbox").click(function () {
    $(this).fadeOut("fast");
  });
  $("div.lightbox .note").click(function () {
    return false;
  });

  $("nav.tabs a").click(function () {
    var rel = $(this).attr('rel');
    var target = $(rel);
    if (target.is("section.tab")) {
      // show the tab pane
      $("section.tab").removeClass('selected');
      target.addClass('selected');

      // select the tab label
      $("nav.tabs a").removeClass('selected');
      $(this).addClass('selected');
      return false;
    }
    return true;
  });

  $("a[href^='#']").click(function () {
    var href = $(this).attr('href');
    var target = $(href);
    if (target.is("section.tab")) {
      // show the tab pane
      $("section.tab").removeClass('selected');
      target.addClass('selected');

      // select the tab label
      $("nav.tabs a").removeClass('selected');
      $("nav.tabs a[rel='"+href+"']").addClass('selected');
      return false;
    }
    return true;
  });

  function update_iconic() {
    $("#options-tab .inventory-iconic-set").hide();
    $("#iconic img").removeClass("selected");
    $('#inventory-iconic-custom').hide();

    var option = $("#inventory-iconic-set option:selected");
    var value = option.val();
    if (value == "default") {
      $('#inventory-iconic').val('default');
      $('#iconic-default').addClass("selected");
      return;
    }

    if (value == "custom") {
      $('#inventory-iconic').val('custom');
      $('#inventory-iconic-custom').show();
      return;
    }

    var setSelect = $(option.attr('rel'));
    setSelect.show();

    option = setSelect.find("option:selected")
    var rel = option.attr('rel');
    var img = $(rel).addClass("selected");
    img.attr('src', img.attr('data-src'));

    $("#inventory-iconic").val(option.val());
  }
  $("#inventory-iconic-set, .inventory-iconic-set").change(update_iconic);
  update_iconic();
});