#include Tools.jsxinc


// Create the PDFSaveOptions object to set the PDF options
var pdfSaveOpts = new PDFSaveOptions();

// Setting PDFSaveOptions properties. Please see the JavaScript Reference
// for a description of these properties.
// Add more properties here if you like
pdfSaveOpts.acrobatLayers = false;
pdfSaveOpts.colorBars = false;
pdfSaveOpts.compatibility = PDFCompatibility.ACROBAT5;
pdfSaveOpts.colorCompression = CompressionQuality.AUTOMATICJPEGHIGH;
//pdfSaveOpts.compatibility = 
pdfSaveOpts.compressArt = true; //default
pdfSaveOpts.embedICCProfile = false;
pdfSaveOpts.enablePlainText = true;
pdfSaveOpts.generateThumbnails = false; // default
pdfSaveOpts.optimization = false; // specifically, optimise for fast web streaming
pdfSaveOpts.pageInformation = false;
pdfSaveOpts.preserveEditability = false;
pdfSaveOpts.viewAfterSaving = false;

pdfSaveOpts.printerResolution = 800.0;
pdfSaveOpts.monochromeDownsampling = 300.0;
pdfSaveOpts.grayscaleDownsampling = 150.0;
pdfSaveOpts.colorDownsampling = 150.0;
  
var originalInteractionLevel = userInteractionLevel;
userInteractionLevel = UserInteractionLevel.DISPLAYALERTS;

var sourceFolder = Folder.selectDialog( 'Select the folder of Illustrator files you want to export as PDFs' );
var destinationFolder = Folder.selectDialog( 'Select the destination folder into which PDFs will be saved' );
var files = sourceFolder.getAllFiles();

log("Exporting "+files.length+" Illustrator files as PDFs", sourceFolder, { "Destination": destinationFolder });

userInteractionLevel = UserInteractionLevel.DONTDISPLAYALERTS;

var success = 0;
var failure = 0;
for ( var i = 0; i < files.length; i++ ) {
  try {
    var file = files[i];
    var doc = app.open(file);

    var filename = file.fullName;
    var targetName = destinationFolder.fullName+filename.substring(sourceFolder.fullName.length, filename.length - 3)+".pdf";
    var targetFile = new File(targetName);

    log("Exporting file as PDF", file);
    doc.saveAs( targetFile, pdfSaveOpts );
    doc.close();
    success++;
  } catch (e) {
    log("Error in file", file, { "Error": e.message } );
    failure++;
  }
}

userInteractionLevel = originalInteractionLevel;

log("Exported "+files.length+" files", false, {'Success': success, 'Failed': failure});
alert("Done!");