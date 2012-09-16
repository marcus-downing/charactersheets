#include Tools.jsxinc

//var sourceFolder = Folder.selectDialog( 'Select the folder with Illustrator files in which you want to replace fonts');
var sourceFolder = Folder("~/Documents/Projects/charactersheets/Pathfinder/Core/Barbarian");
//var csvfile = File.openDialog("Select a CSV file of substitutions", "*.csv");
var csvfile = new File("~/Documents/Projects/charactersheets/Scripts/font substitutions.csv")
var substitutions = csvfile.readCSV();
//new File("~/Desktop/dump.csv").writeCSV(substitutions);

substitutions = substitutions.associate();

for (var i = 0; i < substitutions.length; i++) {
  log('Candidate substitution: '+substitutions[i]['Family']+" / "+substitutions[i]['Style']+" -> "+substitutions[i]['To family']+" / "+substitutions[i]['To style']+" ("+substitutions[i]['To scale']+", "+substitutions[i]['To tracking']+")");
}

//  map the substitutions to real fonts
var substitutions2 = [];
for (var i = 0; i < textFonts.length; i++) {
  var font = textFonts[i];
  for (var j = 0; j < substitutions.length; j++) {
    if (substitutions[j]['To family'] == font.family && substitutions[j]['To style'] == font.style) {
      substitutions[j].textFont = font;
      substitutions2.push(substitutions[j]);
    }
  }
}
substitutions = substitutions2;

//  clean up the scale and kerning values
function cleanNumber(value, defaultValue) {
  try {
    if (value == '')
      return defaultValue;
    value = parseFloat(value);
    if (isNaN(value))
      return defaultValue;
    // if (value < 0.01)
    //   return defaultValue;
    return value;
  } catch(e) {
    return defaultValue;
  }
}

for (var i = 0; i < substitutions.length; i++) {
  substitutions[i]['To scale'] = cleanNumber(substitutions[i]['To scale'], false);
  substitutions[i]['To tracking'] = cleanNumber(substitutions[i]['To tracking'], false);

  var scale = " ("+substitutions[i]['To scale']+", "+substitutions[i]['To tracking']+")";
  if (scale == " (false, false)") scale = '';
  log('Substitution: '+substitutions[i]['Family']+" / "+substitutions[i]['Style']+" -> "+substitutions[i]['To family']+" / "+substitutions[i]['To style']+scale);
}
new File("~/Desktop/dump.csv").writeCSV(substitutions);


//  walk through the folder
var files = sourceFolder.getAllFiles();
log('Replacing fonts in '+files.length+' files', sourceFolder);

var success = 0;
var failure = 0;
for ( var i = 0; i < files.length; i++ ) {
  var file = files[i];
  log('Replacing fonts in file', file);
  try {
    var doc = app.open(file);

    // replace the font
    var frames = doc.textFrames;
    for ( var j = 0; j < frames.length; j++ ) {
      var ranges = frames[j].textRanges;
      for ( var k = 0; k < ranges.length; k++ ) {
        var range = ranges[k];
        for (var f = 0; f < substitutions.length; f++) {
          if (range.characterAttributes.textFont.family == substitutions[f]['Family'] 
            && range.characterAttributes.textFont.style == substitutions[f]['Style']) {
            range.characterAttributes.textFont = substitutions[f].textFont;
            if (substitutions[f]['To scale'])
              range.characterAttributes.size = range.characterAttributes.size * substitutions[f]['To scale'];
            if (substitutions[f]['To tracking'])
              range.characterAttributes.tracking = substitutions[f]['To tracking'];
          }
        }
      }
    }
    redraw();
          
    // Save the file
    doc.save();
    doc.close();
    success++;
  } catch (e) {
    log("Error in file", file, { "Error": e.message } );
    failure++;
  }
}

log("Replaced fonts from "+files.length+" files", false, {'Success': success, 'Failed': failure});
alert("Done!");