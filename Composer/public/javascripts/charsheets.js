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
    var rel = $("#inventory-iconic option:selected").attr('rel');
    $("#iconic img").removeClass("selected");
    var img = $(rel).addClass("selected");
    img.attr('src', img.attr('data-src'));
  }
  $("#inventory-iconic").change(update_iconic);
  update_iconic();
});