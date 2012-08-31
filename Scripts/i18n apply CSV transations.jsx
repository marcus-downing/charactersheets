#include Tools.jsxinc

//var sourceFolder = new Folder( '/Users/Marcus Downing/Documents/Projects/Character Sheets/Pathfinder/Core/X' );
//var destinationFolder = new Folder( '/Users/Marcus Downing/Documents/Projects/Character Sheets/Pathfinder/Core/X2' );
var sourceFolder = Folder.selectDialog( 'Select the folder of Illustrator files in which you want to replace text' );
var destinationFolder = Folder.selectDialog('Select a destination folder into which to save translated files');
var messagesFile = File.openDialog("Translation CSV file", "*.csv");
var messages = messagesFile.readCSV().associate();

for (var i = 0; i < messages.length; i++) {
  messages[i]['Original'] = normalise(messages[i]['Original']);
  messages[i]['Translation'] = normalise(messages[i]['Translation']);
  messages[i]['Part of'] = normalise(messages[i]['Part of']);

  // if (messages[i]['Translation'] && messages[i]['Translation'].length > 0)
  //   alert("Message: "+messages[i]['Original']+" ("+messages[i]['Part of']+") -> "+messages[i]['Translation']);
}

var files = sourceFolder.getAllFiles();
log("Translating strings in "+files.length+" files.", sourceFolder, { 'Messages': messages, 'Destination': destinationFolder });

function trailingWhitespace(text) {
  var text = String(text);
  var trimmed = text.rtrim();
  return text.substring(trimmed.length);
}

function normalise(text) {
  if (typeof text === "undefined") return '';
  text = String(text).trim();
  text = text.replaceAll('\n', '|');
  text = text.replaceAll('\r', '|');
  return text;
}

function denormalise(text) {
  if (typeof text === "undefined") return '';
  text = String(text).trim();
  text = text.replaceAll('|', '\r');
  return text;
}

function patternise(text) {
  text = denormalise(text);
  text = text.replaceAll('/', '\/');
  return '/'+text+'/';
}

function translate(message, partof) {
  message = normalise(message);
  partof = normalise(partof);
  for (var i = 0; i < messages.length; i++) {
    if (messages[i]['Original'] == message && messages[i]['Part of'] == partof) {
      var translation = denormalise(messages[i]['Translation']);
      if (translation.length > 0)
        return translation;
    }
  }
  return false;
}


var count = 0;
var frameNum = 1;
for ( var i = 0; i < files.length; i++ ) {
  var file = files[i];
  try {
    var destinationFile = new File(destinationFolder.fullName+file.fullName.substring(sourceFolder.fullName.length));
    var destinationFolder = destinationFile.parent;
    if (!destinationFolder.exists) destinationFolder.create();

    var doc = app.open(file);

    var frames = doc.textFrames;
    for ( var j = 0; j < frames.length; j++ ) {
      var frame = frames[j];
      var fullstr = frame.contents;
      var fulltranslation = translate(fullstr);
      if (fulltranslation) {
        frame.contents = fulltranslation;
        count++;
        continue;
      }

      // split range based on continuous font, size and colour
      // always keep the splitting rules in sync with the other script!
      var fullrange = frame.textRange;
      var str = '';
      var ranges = frame.textRanges;
      var spanranges = [];
      var prev = false;
      for ( var k = 0; k < ranges.length; k++ ) {
        var range = ranges[k];
        if (prev == false || 
          (  isEqual(range.characterAttributes.fillColor, prev.characterAttributes.fillColor)
          && isEqual(range.characterAttributes.textFont, prev.characterAttributes.textFont)
          && isEqual(range.characterAttributes.size, prev.characterAttributes.size)
          )) {
          str = str+String(range.contents);
          spanranges.push(range);
        } else {
          var translation = translate(str, fullstr);
          if (translation) {
            var span = spanranges[0];
            var trailing = trailingWhitespace(str);
            for (var l = 1; l < spanranges.length; l++) {
              spanranges[l].remove();
            }
            span.characters.addBefore(translation+trailing);
            span.contents = span.contents.substring(0, span.contents.length - 1);

            count++;
          }
          str = String(range.contents);
          spanranges = [ range ];
        }
        prev = range;
      }
      if (str !== '') {
        var translation = translate(str, fullstr);
        if (translation) {
          var span = spanranges[0];
          var trailing = trailingWhitespace(str);
          for (var l = 1; l < spanranges.length; l++) {
            spanranges[l].remove();
          }
          span.characters.addBefore(translation+trailing);
          span.contents = span.contents.substring(0, span.contents.length - 1);

          count++;
        }
      }
      frameNum++;
    }

    var destinationFile = new File(destinationFilename);
    doc.saveAs(destinationFile);
    doc.close();
  } catch (e) {
    log("Error in file", file, { "Error": e.message } );
  }
}

log("Translated "+count+" strings from "+files.length+" files");
alert("Done!");