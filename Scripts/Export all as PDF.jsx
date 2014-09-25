#include Tools.jsxinc


// var sourceFolder = new Folder( '/Users/Marcus Downing/Documents/GitHub/charactersheets/Pathfinder/Archetypes/Monk' );
// var destinationFolder = new Folder( '/Users/Marcus Downing/Documents/GitHub/charactersheets/Composer/public/pdf/pathfinder/Archetypes/Monk' );
var sourceFolder = baseFolder;
var destinationFolder = baseFolder+'Composer 2.1.3/public/pdf/'

exportFolderAsPDF(new Folder(sourceFolder+'Pathfinder'), new Folder(destinationFolder+'pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'3.5'), new Folder(destinationFolder+'dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'All'), new Folder(destinationFolder+'pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'All'), new Folder(destinationFolder+'dnd35'));

alert("Done!");