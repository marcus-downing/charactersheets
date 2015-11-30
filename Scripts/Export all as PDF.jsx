#include Tools.jsxinc

#include i18n_tools.jsxinc

// var sourceFolder = new Folder( '/Users/Marcus Downing/Documents/GitHub/charactersheets/Pathfinder/Archetypes/Monk' );
// var destinationFolder = new Folder( '/Users/Marcus Downing/Documents/GitHub/charactersheets/Composer/public/pdf/pathfinder/Archetypes/Monk' );
var sourceFolder = baseFolder;
var destinationFolder = baseFolder+'../Character Sheets PDF/'

i18n.init();

exportFolderAsPDF(new Folder(sourceFolder+'Pathfinder'), new Folder(destinationFolder+'pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'3.5'), new Folder(destinationFolder+'dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'All'), new Folder(destinationFolder+'pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'All'), new Folder(destinationFolder+'dnd35'));

alert("English done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/Pathfinder'), new Folder(destinationFolder+'languages/italian/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/3.5'), new Folder(destinationFolder+'languages/italian/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/All'), new Folder(destinationFolder+'languages/italian/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Italian/All'), new Folder(destinationFolder+'languages/italian/dnd35'));

alert("Italian done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/Pathfinder'), new Folder(destinationFolder+'languages/spanish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/3.5'), new Folder(destinationFolder+'languages/spanish/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/All'), new Folder(destinationFolder+'languages/spanish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Spanish/All'), new Folder(destinationFolder+'languages/spanish/dnd35'));

alert("Spanish done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/Pathfinder'), new Folder(destinationFolder+'languages/polish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/3.5'), new Folder(destinationFolder+'languages/polish/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/All'), new Folder(destinationFolder+'languages/polish/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Polish/All'), new Folder(destinationFolder+'languages/polish/dnd35'));

alert("Polish done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Portuguese/Pathfinder'), new Folder(destinationFolder+'languages/portuguese/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Portuguese/3.5'), new Folder(destinationFolder+'languages/portuguese/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Portuguese/All'), new Folder(destinationFolder+'languages/portuguese/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Portuguese/All'), new Folder(destinationFolder+'languages/portuguese/dnd35'));

alert("Portuguese done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/French/Pathfinder'), new Folder(destinationFolder+'languages/french/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/French/3.5'), new Folder(destinationFolder+'languages/french/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/French/All'), new Folder(destinationFolder+'languages/french/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/French/All'), new Folder(destinationFolder+'languages/french/dnd35'));

alert("French done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/German/Pathfinder'), new Folder(destinationFolder+'languages/german/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/German/3.5'), new Folder(destinationFolder+'languages/german/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/German/All'), new Folder(destinationFolder+'languages/german/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/German/All'), new Folder(destinationFolder+'languages/german/dnd35'));

alert("German done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/Russian/Pathfinder'), new Folder(destinationFolder+'languages/russian/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Russian/3.5'), new Folder(destinationFolder+'languages/russian/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Russian/All'), new Folder(destinationFolder+'languages/russian/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/Russian/All'), new Folder(destinationFolder+'languages/russian/dnd35'));

alert("Russian done...");

exportFolderAsPDF(new Folder(sourceFolder+'Languages/US English/Pathfinder'), new Folder(destinationFolder+'languages/american/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/US English/3.5'), new Folder(destinationFolder+'languages/american/dnd35'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/US English/All'), new Folder(destinationFolder+'languages/american/pathfinder'));
exportFolderAsPDF(new Folder(sourceFolder+'Languages/US English/All'), new Folder(destinationFolder+'languages/american/dnd35'));

alert("US English done...");