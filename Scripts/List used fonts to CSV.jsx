#include "tools.jsxinc"

var sourceFolder = Folder.selectDialog( 'Select the folder of Illustrator files in which you want to find fonts' );
var outfile = File.saveDialog("Save CSV file", "*.csv");

var files = sourceFolder.getAllFiles();
log('Scanning '+files.length+' files for fonts', sourceFolder);

var data = [];
var success = 0;
var failure = 0;
for ( var i = 0; i < files.length; i++ ) {
  var file = files[i];
  log('Scanning file'+(i+1)+' of '+files.length, file);
  try {
    var doc = app.open(file);

    var frames = doc.textFrames;
    for ( var j = 0; j < frames.length; j++ ) {
      var ranges = frames[j].textRanges;

      range_loop:
      for ( var k = 0; k < ranges.length; k++ ) {
        var range = ranges[k];
        var font = range.characterAttributes.textFont;
        var entry = {
          'Font name': font.name, 
          'Family': font.family, 
          'Style': font.style, 
          'Uses': 1,
          'First use': files[i].name
        };

        // check if the font's been lodged before
        var elen = data.length;
        for (var e = 0; e < elen; e++) {
          var dentry = data[e];
          if (entry['Font name'] == dentry['Font name']) {
            dentry['Uses']++;
            continue range_loop;
          }
        }
        data.push(entry);
      }
    }

    doc.close();
    success++;
  } catch (e) {
    log("Error in file", file, { "Error": e.message } );
    failure++;
  }
}

var data = data.dissociate(['Font name', 'Family', 'Style']);
outfile.writeCSV(data);
log("Extracted "+(data.length-1)+" fonts from "+files.length+" files", false, {'Success': success, 'Failed': failure});
alert("Done!");