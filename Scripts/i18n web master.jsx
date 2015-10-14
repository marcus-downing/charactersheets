#include Tools.jsxinc
#include i18n_tools.jsxinc

/*
i18n web master
Walk the entire repo producing a CSV of all translatable strings in all files,
in a format suitable for the Translator app
*/

i18n.init();

var webMaster = {
  entries: {},
  volume: "",
  game: "",
  level: 1,

  baseURI: new Folder(baseFolder).absoluteURI+"/",

  clear: function() {
    entries = {};
  },

  pushEntry: function(text, partOf, filename) {
    var replacement = i18n.getTranslationReplacement(text, partOf, filename);
    text = replacement[0];
    partOf = replacement[1];

    var normal = i18n.normalise(text);
    // log("Push entry", text, normal, partOf, filename);
    text = normal;
    if (text.length <= 1) return;

    var simpletext = text.replace(/[^a-zA-Z]*/g, '');
    if (simpletext.length == 0) return;

    partOf = i18n.normalise(partOf);
    if (partOf == text) partOf = '';

    filename = filename.substring(0, filename.length - 3);

    var key = text+"%%%"+partOf+"%%%"+filename
    if (this.entries[key]) {
      this.entries[key].Count++;
    } else {
      this.entries[key] = entry = {
        'Original': text,
        'Part of': partOf,
        'Count': 1,
        'File': filename,
        'Volume': this.volume,
        'Game': this.game,
        'Level': this.level
      };
    }
  },

  extractFile: function(file) {
    var num = 0;
    try {
      var filename = file.absoluteURI.replace(webMaster.baseURI, '');
      log("Extracting from file", file.fullName, filename)
      var doc = app.open(file);

      var frames = doc.textFrames;
      for ( var j = 0; j < frames.length; j++ ) {
        var frame = frames[j];
        var partspushed = 0;
        var fullrange = frame.textRange;
        var fullstr = fullrange.contents;

        // split range based on continuous font, size and colour
        var str = '';
        var ranges = frame.textRanges;
        var prev = false;
        for ( var k = 0; k < ranges.length; k++ ) {
          var range = ranges[k];
          if (prev == false || 
            (  isEqual(range.characterAttributes.fillColor, prev.characterAttributes.fillColor)
            && isEqual(range.characterAttributes.textFont, prev.characterAttributes.textFont)
            && isEqual(range.characterAttributes.size, prev.characterAttributes.size)
            )) {
            str = str+String(range.contents);
          } else {
            this.pushEntry(str, fullstr, filename);
            partspushed++;
            num++;
            str = String(range.contents);
          }
          prev = range;
        }
        if (str !== '') {
          this.pushEntry(str, fullstr, filename);
          partspushed++;
          num++;
        }

        if (partspushed == 0) {
          this.pushEntry(fullstr, "", filename);
            num++;
        }
      }
      doc.close();

      var additions = i18n.getPageAdditions(filename);
      for ( var i = 0; i < additions.length; i++ ) {
        var add = additions[i]
        this.pushEntry(add, add, filename)
      }
    } catch (e) {
      log("Error in file", file, { "Error": e.message } );
    }
    i18n.tick();
    return num;
  },

  extractFolder: function(folder) {
    var num = 0;
    var files = folder.getAllFiles();
    log('Scanning '+files.length+' files for translatable strings', folder);

    for ( var i = 0; i < files.length; i++ ) {
      var file = files[i];
      log("Scanning file "+(i+1)+" of "+files.length, file);
      num += this.extractFile(file);
    }
    log("Extracted "+num+" strings from "+files.length+" files");
  },

  saveCSV: function(file) {
    var entries = [];
    for(var key in this.entries) {
      entries.push(this.entries[key]);
    }

    data = entries.dissociate(['Original', 'Part of', 'Count', 'File']);
    file.writeCSV(data);
  }
}

log("Base URI", webMaster.baseURI)

webMaster.game = "Pathfinder"
webMaster.level = 1;
webMaster.volume = "Core Rulebook";
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Core'));
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Extra'));
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/GM'));

webMaster.level = 2;
webMaster.volume = "Advanced Players Guide"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Advanced'));
webMaster.volume = "Ultimate Magic"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Ultimate Magic'));
webMaster.volume = "Ultimate Combat"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Ultimate Combat'));
webMaster.volume = "Advanced Class Guide"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Advanced Class Guide'));

webMaster.level = 3;
webMaster.volume = ""
webMaster.extractFolder(new Folder(baseFolder+'All'));
webMaster.extractFolder(new Folder(baseFolder+'Extra'));
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Archetypes'));
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Prestige Classes'));
webMaster.volume = "Mythic Adventures"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Mythic'));

webMaster.level = 4;
webMaster.volume = "Psionics Unleashed"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Psionics'));
webMaster.volume = "Tome of Secrets"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/Tome of Secrets'));
webMaster.volume = "NeoExodus"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/NeoExodus'));
webMaster.volume = "TPK"
webMaster.extractFolder(new Folder(baseFolder+'Pathfinder/TPK'));


webMaster.game = "3.5";
webMaster.level = 1;
webMaster.volume = "Players Handbook"
webMaster.extractFolder(new Folder(baseFolder+'3.5/Core'));
webMaster.extractFolder(new Folder(baseFolder+'3.5/Barbarian'));
webMaster.extractFolder(new Folder(baseFolder+'3.5/Variants'));
webMaster.extractFolder(new Folder(baseFolder+'3.5/DM'));

webMaster.level = 2;
webMaster.volume = "";
webMaster.extractFolder(new Folder(baseFolder+'3.5/Extended'));

webMaster.level = 3;
webMaster.volume = ""
webMaster.extractFolder(new Folder(baseFolder+'All'));
webMaster.extractFolder(new Folder(baseFolder+'Extra'));
webMaster.volume = "Dragon Compendium";
webMaster.extractFolder(new Folder(baseFolder+'3.5/Dragon Compendium'));

webMaster.level = 4;
webMaster.volume = "";
webMaster.extractFolder(new Folder(baseFolder+'3.5/Psionics'));
webMaster.extractFolder(new Folder(baseFolder+'3.5/Tomes'));
webMaster.extractFolder(new Folder(baseFolder+'3.5/Incarnum'));


webMaster.game = "Extra"
webMaster.level = 2;
webMaster.volume = "";
webMaster.extractFolder(new Folder(baseFolder+'All'));
webMaster.extractFolder(new Folder(baseFolder+'Extra'));


webMaster.saveCSV(new File(baseFolder+'Languages/Master.csv'));
webMaster.clear();
alert("Done");