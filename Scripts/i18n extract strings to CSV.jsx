#include Tools.jsxinc

var sourceFolder = Folder.selectDialog( 'Select the folder of Illustrator files in which you want to find text' );
var outfile = File.saveDialog("Save translation CSV file", "*.csv");

var files = sourceFolder.getAllFiles();
log('Scanning '+files.length+' files for translatable strings', sourceFolder);

function normalise(text) {
  if (typeof text === "undefined") return '';
  text = String(text).trim();
  text = text.replaceAll('\n', '|');
  text = text.replaceAll('\r', '|');
  return text;
}

var data = [];
function pushEntry(frameNum, text, partOf, filename) {
  text = normalise(text);
  if (text.length <= 1) return;

  var simpletext = text.replace(RegExp('[^a-zA-Z]*', 'g'), '');
  if (simpletext.length == 0) return;

  partOf = normalise(partOf);
  if (partOf == text) partOf = '';

  filename = filename.substring(0, filename.length - 3);

  var entry = {
    'Frame number': frameNum,
    'Original': text,
    'Translation': '',
    'Part of': partOf,
    'Count': 1,
    'Files': [ filename ]
  };

  // check if the item has been stored before, and merge if necessary
  var len = data.length;
  for (var e = 0; e < len; e++) {
    if (text === data[e]['Original'] && partOf === data[e]['Part of']) {
      data[e]['Count']++;
      if (!data[e]['Files'].contains(filename)) data[e]['Files'].push( filename );
      return;
    }
  }
  data.push(entry);
}



var frameNum = 1;
for ( var i = 0; i < files.length; i++ ) {
  var file = files[i];
  log("Scanning file "+(i+1)+" of "+files.length, file);
  try {
    var filename = file.name;
    var doc = app.open(file);

    var frames = doc.textFrames;
    for ( var j = 0; j < frames.length; j++ ) {
      var frame = frames[j];
      var fullrange = frame.textRange;
      var fullstr = fullrange.contents;
      pushEntry(frameNum, fullstr, '', filename);

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
          pushEntry(frameNum, str, fullstr, filename);
          str = String(range.contents);
        }
        prev = range;
      }
      if (str !== '') pushEntry(frameNum, str, fullstr, filename);
      frameNum++;
    }
    doc.close();
  } catch (e) {
    log("Error in file", file, { "Error": e.message } );
  }
}

for (var i = 0; i < data.length; i++) {
  data[i]['File count'] = data[i]['Files'].length;
  data[i]['Files'] = any2string(data[i]['Files']);
}

data = data.dissociate(['Frame number', 'Original', 'Translation', 'Count']);
outfile.writeCSV(data);
log("Extracted "+(data.length-1)+" strings from "+files.length+" files", false, {'Success': files.length});
alert("Done!");