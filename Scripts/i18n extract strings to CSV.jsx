#include Tools.jsxinc
#include i18n_tools.jsxinc

i18n.init();

var sourceFolder = Folder.selectDialog( 'Select the folder of Illustrator files in which you want to find text' );
var outfile = File.saveDialog("Save translation CSV file", "*.csv");

i18n.init();
i18n.extractFolder(sourceFolder);
i18n.saveCSV(outfile);

alert("Done!");