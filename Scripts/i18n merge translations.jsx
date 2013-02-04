#include Tools.jsxinc
#include i18n_tools.jsxinc

/*
i18n merge translations
Reduce two CSVs to produce a single translations file

 - master file: a newly exported file with the right strings, or a bigger file with the full set of strings
 - reference file: a file that contains useful translations but may not be complete or up to date
*/

i18n.init();

var masterFile = File.openDialog( 'Select "master" file with correct origins', "*.csv" );
var referenceFile = File.openDialog( 'Select "reference" file with translations to merge', "*.csv" );
var outFile = File.saveDialog( 'Save merged translation file', "*.csv" );

i18n.messages = i18n.loadCSV(masterFile);
var reference = i18n.loadCSV(referenceFile);

i18n.reduce(reference);
i18n.saveCSV(outFile);

alert("Done!");

/*
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
*/