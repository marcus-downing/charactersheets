#include Tools.jsxinc

var sourceFolder = Folder.selectDialog( 'Select the folder of Illustrator files to check', baseFolder);
var files = sourceFolder.getAllFiles();

log('Starting check all files', sourceFolder);
var success = 0;
var failed = 0;

for (var i = 0; i < files.length; i++) {
  var file = files[i];
  log('Checking file '+(i+1)+' of '+files.length, file);
  try {
    var doc = app.open(file);
    doc.close();
    success++;
  } catch (e) {
    log('File may be broken!', file);
    failed++;
  }
}

log('Finished check all files', false, {'Success': success, 'Failed': failed});
alert("Done!");