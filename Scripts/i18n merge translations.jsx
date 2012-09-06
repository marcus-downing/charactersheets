#include Tools.jsxinc

/*
i18n merge translations
Produce a single translations file

 - master file: a newly exported file with the right strings, or a bigger file with the full set of strings
 - reference file: a file that contains useful translations but may not be complete or up to date
*/

var masterFile = File.openDialog( 'Select "master" file with correct origins', "*.csv" );
var referenceFile = File.openDialog( 'Select "reference" file with translations to merge', "*.csv" );
var outfile = File.saveDialog( 'Save merged translation file', "*.csv" );

var master = masterFile.readCSV().associate();
var reference = referenceFile.readCSV().associate();

for (var i = 0; i < master.length; i++) {
  var translation = false;
  for (var j = 0; j < reference.length; j++) {
    if (master[i]['Original'] == reference[j]['Original'] && master[i]['Part of'] == reference[j]['Part of']) {
      master[i]['Translation'] = reference[j]['Translation'];
      break;
    }
  }
}

outfile.writeCSV(master.dissociate());