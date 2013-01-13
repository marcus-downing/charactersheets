#include Tools.jsxinc


/*
i18n combine strings
Add two CSVs to produce a single translations file

 - first file: the file with the most important strings
 - second file: a file with some extra translations
*/

var firstFile = File.openDialog( 'Select the first messages file', "*.csv" );
var secondFile = File.openDialog( 'Select a second messages file', "*.csv" );
var outfile = File.saveDialog( 'Save merged messages file', "*.csv" );

var messages = firstFile.readCSV().associate();
var second = secondFile.readCSV().associate();

// merge the files
second:
for (var i = 0; i < second.length; i++) {
  for (var j = 0; j < messages.length; j++) {
    if (messages[j]['Original'] == second[i]['Original'] && messages[j]['Part of'] == second[i]['Part of']) {
      continue second;
    }
  }
  var add = second[i];
  add['Frame number'] += 10000;
  messages.push(add);
}

// renumber frames
var frames = [];
var nframe = 0;
var lastframe = null;
for (var i = 0; i < messages.length; i++) {
	if (messages[i]['Frame number'] !== lastframe)
		nframe++;
	log("Comparing frame numbers", messages[i]['Frame number'], lastframe);
	lastframe = messages[i]['Frame number'];
	messages[i]['Frame number'] = nframe;
}

outfile.writeCSV(messages.dissociate());
alert("Done!");