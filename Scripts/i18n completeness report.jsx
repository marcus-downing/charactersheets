#include Tools.jsxinc
#include i18n_tools.jsxinc

/*
i18n completeness report
Produce a report on how complete a set of translations is

 - master file: a newly exported file with the right strings, or a bigger file with the full set of strings
 - reference file: a file that contains useful translations but may not be complete or up to date
*/

i18n.init();

var masterFile = File.openDialog( 'Select "master" file with correct origins', "*.csv" );
var referenceFile = File.openDialog( 'Select a CSV file containing translations', "*.csv" );

// collate messages by file
var master = i18n.loadCSV(masterFile);
var fileMessages = {};
for (var i = 0; i < master.length; i++) {
	var message = master[i];
	var messageID = {
		'Original': message['Original'],
		'Part of': message['Part of'],
	};
	var files = message['Files'];
	for (var j = 0; j < files.length; j++) {
		var file = files[j];
		if (fileMessages.hasOwnProperty(file))
			fileMessages[file].push(messageID);
		else
			fileMessages[file] = [messageID];
	}
}
log('Found '+fileMessages.length+' origin files from '+master.length+' messages');

// compare and count
var reference = referenceFile.readCSV().associate();
var completeness = [];

for (var file in fileMessages) {
	var messages = fileMessages[file];
	var c = 0;
	var n = 0;
	messages:
	for (var j = 0; j < messages.length; j++) {
		var message = messages[j];
		if (!i18n.isTranslatable(message['Original']))
			continue;
		n++
		for (var k = 0; k < reference.length; k++) {
			var ref = reference[k];
			if (message['Original'] == ref['Original'] && message['Part of'] == ref['Part of']) {
				if (ref['Translation'] != '')
					c++;
				continue messages;
			}
		}
	}
	completeness.push({
		'File': file,
		'Count': n,
		'Complete': c,
		'Completeness %': Math.round(100.0 * c / n)
	});
}

//  report out
var outFile = File.saveDialog( 'Save completeness report', "*.csv" );
outFile.writeCSV(completeness.dissociate(['File', 'Count', 'Complete', 'Completeness %']));
alert("Done!");