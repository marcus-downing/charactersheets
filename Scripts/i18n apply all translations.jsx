#include Tools.jsxinc

#include i18n_tools.jsxinc

#include i18n_translate.jsxinc


i18n.init();

var folders = [ "Pathfinder", "3.5", "All", "Extra" ];
var languages = [ "Italian", "Spanish", "Polish", "German", "French", "Portuguese", "Russian", "US English" ];
log("i18n: Preparing to translate into", languages);

for (var i = 0; i < languages.length; i++) {
	var language = languages[i];

	i18n.loadTranslations(new File(baseFolder + 'Languages/' + language + '.csv'));

	for (var j = 0; j < folders.length; j++) {
		var folder = folders[j];
		var srcFolder = new Folder(baseFolder + folder);
		var dstFolder = new Folder(baseFolder + 'Languages/' + language + '/' + folder);
		i18n.applyTranslationsFolder(srcFolder, dstFolder);
	}
}

log("i18n: Translated "+i18n.countTranslatedLines+" strings from "+i18n.countTranslatedFiles+" files");
alert("Done!");
