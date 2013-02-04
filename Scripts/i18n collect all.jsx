#include Tools.jsxinc
#include i18n_tools.jsxinc

var base = '/Users/Marcus Downing/Documents/GitHub/charactersheets/';
i18n.init();

// Pathfinder Core
i18n.extractFile(new File(base+'Pathfinder/Core/Character Info.ai'));
i18n.extractFile(new File(base+'Pathfinder/Core/Combat.ai'));
i18n.extractFolder(new Folder(base+'Pathfinder/Core'));
i18n.saveCSV(new File(base+'Languages/Template/Pathfinder - Core.csv'));

// Pathfinder Advanced
i18n.extractFolder(new Folder(base+'Pathfinder/Advanced'));
i18n.extractFolder(new Folder(base+'Pathfinder/Ultimate Magic'));
i18n.extractFolder(new Folder(base+'Pathfinder/Ultimate Combat'));
i18n.extractFolder(new Folder(base+'Pathfinder/Extra'));
i18n.saveCSV(new File(base+'Languages/Template/Pathfinder - Advanced.csv'));

// Pathfinder Everything
i18n.extractFolder(new Folder(base+'Pathfinder/Archetypes'));
i18n.extractFolder(new Folder(base+'Pathfinder/GM'));

i18n.extractFolder(new Folder(base+'Pathfinder/Psionics'));
i18n.extractFolder(new Folder(base+'Pathfinder/Tome of Secrets'));
i18n.saveCSV(new File(base+'Languages/Template/Pathfinder - Everything.csv'));

alert("Done");