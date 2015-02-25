#include Tools.jsxinc


// var sourceFolder = new Folder( '/Users/Marcus Downing/Documents/GitHub/charactersheets/Pathfinder/Archetypes/Monk' );
// var destinationFolder = new Folder( '/Users/Marcus Downing/Documents/GitHub/charactersheets/Composer/public/pdf/pathfinder/Archetypes/Monk' );
var sourceFolder = baseFolder;
var destinationFolder = baseFolder+'../..//Character Sheets Website/public/pdf/'

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/Pathfinder'), new Folder(destinationFolder+'languages/italian/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/3.5'), new Folder(destinationFolder+'languages/italian/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/All'), new Folder(destinationFolder+'languages/italian/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/All'), new Folder(destinationFolder+'languages/italian/dnd35'));
alert("Italian done");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/Pathfinder'), new Folder(destinationFolder+'languages/spanish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/3.5'), new Folder(destinationFolder+'languages/spanish/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/All'), new Folder(destinationFolder+'languages/spanish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/All'), new Folder(destinationFolder+'languages/spanish/dnd35'));
alert("Spanish done");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/Pathfinder'), new Folder(destinationFolder+'languages/polish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/3.5'), new Folder(destinationFolder+'languages/polish/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/All'), new Folder(destinationFolder+'languages/polish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/All'), new Folder(destinationFolder+'languages/polish/dnd35'));
alert("Polish done");