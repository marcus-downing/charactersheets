#include Tools.jsxinc
#include i18n_tools.jsxinc

/*
i18n combine strings
Add two CSVs to produce a single translations file

 - first file: the file with the most important strings
 - second file: a file with some extra translations
*/

i18n.init();

var firstFile = File.openDialog( 'Select the first messages file', "*.csv" );
var secondFile = File.openDialog( 'Select a second messages file', "*.csv" );
var outFile = File.saveDialog( 'Save merged messages file', "*.csv" );

i18n.messages = i18n.loadCSV(firstFile);
var second = i18n.loadCSV(secondFile);

i18n.combine(second);
i18n.renumber();
i18n.saveCSV(outFile);

alert("Done!");

/*
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
*/