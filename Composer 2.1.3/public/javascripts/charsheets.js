$(function() {
  $("html, body").addClass("postload");

  $("#start-single").click(function () {
    $("#class-tab-link, #options-tab-link, #download-tab-link").show();
    $("#start-tab-link, #party-tab-link, #gm-tab-link").hide();
    $("#class-tab-link").show().click();
    $("#add-to-party").hide();
    $("#party-readout").hide();
    $("#start-type").val('single');
    return false;
  });

  $("#start-party").click(function () {
    $("#class-tab-link, #options-tab-link, #party-tab-link, #download-tab-link").show();
    $("#start-tab-link, #gm-tab-link").hide();
    $("#class-tab-link").click();
    $("#party-readout").show();
    $("#add-to-party").show();
    $("#start-type").val('party');
    return false;
  });

  $("#start-gm").click(function () {
    $("#gm-tab-link, #download-tab-link").show();
    $("#start-tab-link, #party-tab-link, #class-tab-link, #options-tab-link").hide();
    $("#gm-tab-link").click();
    $("#add-to-party").hide();
    $("#party-readout").hide();
    $("#start-type").val('gm');
  });

  $("#start-all").click(function () {
    $("#download-tab-link").show();
    $("#start-tab-link, #class-tab-link, #options-tab-link, #party-tab-link, #gm-tab-link").hide();
    $("#download-tab-link").show().click();
    $("#add-to-party").hide();
    $("#party-readout").hide();
    $(".wizardnav").hide();
    $("#start-type").val('all');
    return false;
  });

  $("#include-pathfinder-society").change(function () {
    if ($(this).is(":checked")) {
      $("#include-background").prop('checked', true);
    }
  });

  $("#simple").change(function () {
    if ($(this).is(":checked")) {
      $("#more").prop('checked', false);
    }
  });

  $("#more").change(function () {
    if ($(this).is(":checked")) {
      $("#simple").prop('checked', false);
    }
  });

  $("input[name=mini-size]").change(function () {

  });

  $("a.lightbox").click(function () {
    var id = $(this).attr('rel');
    var lightbox = $(id);
    if (lightbox) {
      var img = lightbox.find("img");
      var src = img.attr('src');
      img.attr('src', '');
      img.attr('src', src).load(function () {
        console.log("image loaded!");
        var outer = lightbox.innerHeight();
        var inner = img.outerHeight();
        var margin = (outer - inner) / 2;
        lightbox.find("> *").css("margin-top", margin+"px");
      });
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

      // april fool
      if ($("body").is(".april-fool") && cornify_add && Math.random() > 0.8) {
        cornify_add();
      }
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

      // april fool
      if ($("body").is(".april-fool") && cornify_add && Math.random() > 0.8) {
        cornify_add();
      }
      return false;
    }
    return true;
  });

  function update_character() {
    if ($("#simple-section").length) {
      if ($("#class-Barbarian, #class-Ranger, #class-Ardent, #class-Divine-Mind, #class-Lurk, #class-Psion, #class-Psychic-Warrior, #class-Soulknife, #class-Wilder").is(":checked")) {
        $("#simple-section").addClass("disabled");
        $("#simple").attr("disabled", true).removeAttr("checked");
      } else {
        $("#simple-section").removeClass("disabled");
        $("#simple").removeAttr("disabled");
      }

      if ($("#simple").is(":checked")) {
        $("#iconic-section").addClass("disabled");
        $("#inventory-iconic-set").val('default').attr("disabled", true);
      } else {
        $("#iconic-section").removeClass("disabled");
        $("#inventory-iconic-set").removeAttr("disabled");
      }

      update_iconic();
    }
  }
  $("#class-tab input, #class-tab select, #simple").change(update_character);

  // iconics
  $("#select-iconic-button").click(function () {
    $("#blanket, #iconic-select-dialog").fadeIn("fast");
  });

  $("#iconic-set-list a").click(function () {
    var setId = $(this).data('set-id');
    $("#iconic-set-list a").removeClass("selected");
    $(this).addClass("selected");
    $("#iconic-image-list > div").removeClass("selected");
    $("#iconic-image-list-"+setId).addClass("selected");
    $("#iconic-image-list-"+setId+" img").each(function () {
      $(this).attr("src", $(this).data("src"));
    });
  });

  $("#iconic-image-list a").click(function () {
    var iconic = $(this).data("id");
    if (iconic == "custom") {
      $("#iconic-image-list > div").removeClass("selected");
      $("#iconic-image-list-custom").addClass("selected");
    } else {
      $("#inventory-iconic").val(iconic);
      $("#iconic img").removeClass("selected");
      $("#iconic-"+iconic).addClass("selected").attr('src', $("#iconic-"+iconic).data('src'));
      // close
      $("#blanket, #download-thanks-dialog, #iconic-select-dialog").fadeOut("fast");
    }
  });

  $("#iconic-custom-file-ok-button").click(function () {
    $("#inventory-iconic").val("custom");
    $("#iconic img").removeClass("selected");
    $("#iconic-custom").addClass("selected");
    $("#blanket, #iconic-select-dialog").fadeOut("fast");
  })

  $("#iconic-custom-file-cancel-button").click(function () {
    $("#blanket, #iconic-select-dialog").fadeOut("fast");
  });

  // logos
  $("#select-logo-button").click(function () {
    $("#blanket, #logo-select-dialog").fadeIn("fast");
  });

  $("#logo-list a").click(function () {
    var logo = $(this).data("id");
    $("#logo-select").val(logo);
    $("#logo img").removeClass("selected");
    $("#logo-"+logo).addClass("selected");
    // close
    $("#blanket, #logo-select-dialog").fadeOut("fast");
  })


  var nextcharid = 1;
  $("#add-to-party").click(function () {
    var form = $("#build-my-character");
    var inputs = form.find("input").not("[data-charid]");
    var charid = nextcharid; nextcharid++;

    // collect all the character data
    var chardata = {};
    inputs.each(function () {
      var input = $(this);
      if (input.attr('type') == 'radio' && !input.is(":checked")) {
        return;
      }
      var name = input.attr('name');
      var value = input.attr('value');
      if (input.attr('type') == 'checkbox') {
        value = input.is(":checked") ? "on" : "";
      }
      chardata[name] = value;
    });

    // store the data in hidden fields
    for (name in chardata) {
      var value = chardata[name];
      $("<input type='hidden' name='char-"+charid+"-"+name+"' data-charid='"+charid+"' />").val(value).appendTo(form);
    }

    // interpret the data
    var classes = [];
    for (name in chardata) {
      if (name.slice(0, 6) == 'class-' && chardata[name] == 'on') {
        classes.push(name.slice(6));
      }
    }

    // add the character to the list
    var readout = $("#party-readout ul");
    var img = $("#inventory-iconic").val();
    var imgsrc = $("#iconic-"+img).attr('src');
    $("<li><img src='"+imgsrc+"'/><span>"+classes.join(", ")+"</span></li>").appendTo(readout);
    var charids = $("#charids");
    var ids = charids.val().split(",");
    ids.push(charid);
    charids.val(ids.join(","));

    // reset the data
    $("#class-tab input[type=checkbox]").prop("checked", false);
    $("#iconic img").removeClass("selected");
    $("#iconic-generic").addClass("selected");
    $("#inventory-iconic").val("generic");

    // move along
    $("#party-tab-link").click();
  });

  $("#build-my-character").submit(function () {/*
    var path = "";
    $("#class-tab input:checkbox:checked").each(function () {
      var name = $(this).data('classname');

      $("#variant-"+code+" option:selected").each(function () {
        name = $(this).attr('value');
      })

      path = path+"/"+name;
    });

    var url = "https://flattr.com/submit/auto?user_id=marcusdowning&url=http://charactersheets.minotaur.cc"+path;
    $("a#flattr").attr('href', url);*/

    $("#blanket, #download-thanks-dialog").fadeIn("fast");
  });

  $("#close").click(function () {
    $("#blanket, #download-thanks-dialog, #iconic-select-dialog").fadeOut("fast");
  });
});

    window._idl = {};
    _idl.variant = "banner";
    (function() {
        var idl = document.createElement('script');
        idl.type = 'text/javascript';
        idl.async = true;
        idl.src = ('https:' == document.location.protocol ? 'https://' : 'http://') + 'members.internetdefenseleague.org/include/?url=' + (_idl.url || '') + '&campaign=' + (_idl.campaign || '') + '&variant=' + (_idl.variant || 'banner');
        document.getElementsByTagName('body')[0].appendChild(idl);
    })();
