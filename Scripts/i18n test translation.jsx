#include Tools.jsxinc

#include i18n_tools.jsxinc

#include i18n_translate.jsxinc


i18n.init();

(function () {
	var srcFile = new File(baseFolder+'3.5/Core/Cleric.ai');
	var dstFile = new File(baseFolder+'Languages/Italian/3.5/Core/Cleric.ai');
	var messagesFile = new File(baseFolder+'Languages/Italian.csv');

	i18n.loadTranslations(messagesFile);
	i18n.applyTranslationsFile(srcFile, dstFile);

	log("i18n: Translated "+i18n.countTranslatedLines+" strings from "+i18n.countTranslatedFiles+" files");
	alert("Done!");
})();