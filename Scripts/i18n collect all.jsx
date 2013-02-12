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
i18n.extractFolder(new Folder(base+'Pathfinder/NeoExodus'));

i18n.extractFolder(new Folder(base+'All'));
i18n.extractFolder(new Folder(base+'Extra'));
i18n.saveCSV(new File(base+'Languages/Template/Pathfinder - Everything.csv'));


// 3.5 core
i18n.init();

i18n.extractFile(new File(base+'3.5/Core/Character Info - Simple.ai'));
i18n.extractFile(new File(base+'3.5/Core/Character Info.ai'));
i18n.extractFile(new File(base+'3.5/Core/Blank Character Info.ai'));
i18n.extractFile(new File(base+'3.5/Core/Combat - Simple.ai'));
i18n.extractFile(new File(base+'3.5/Core/Combat.ai'));
i18n.extractFolder(new Folder(base+'3.5/Core'));
i18n.extractFolder(new Folder(base+'3.5/Barbarian'));
i18n.saveCSV(new File(base+'Languages/Template/3.5 - Core.csv'));

i18n.extractFolder(new Folder(base+'3.5/Extended'));
i18n.extractFolder(new Folder(base+'3.5/Variants'));
i18n.extractFolder(new Folder(base+'3.5/Dragon Compendium'));
i18n.extractFolder(new Folder(base+'3.5/Incarnum'));
i18n.extractFolder(new Folder(base+'3.5/Psionics'));
i18n.extractFolder(new Folder(base+'3.5/Tomes'));
i18n.saveCSV(new File(base+'Languages/Template/3.5 - Everything.csv'));


alert("Done");