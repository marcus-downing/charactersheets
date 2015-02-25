#include Tools.jsxinc

#include i18n_tools.jsxinc

#include i18n_translate.jsxinc


i18n.init();

var folders = [ "Pathfinder", "3.5", "All" ];

var languages = [ "Italian", "Spanish", "Polish", "French", "Portuguese", "Russian" ];


for (var i = 0; i < languages.length; i++) {
	var language = languages[i];

	i18n.loadTranslations(baseFolder + '/Languages/' + language + '.csv');

	for (var j = 0; j < folders; j++) {
		var folder = folders[j];
		var srcFolder = baseFolder + '/' + folder;
		var dstFolder = baseFolder + '/Languages/' + language + '/' + folder;
		i18.applyTranslationsFolder(srcFolder, dstFolder);
	}
}
