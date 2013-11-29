#include Tools.jsxinc
#include i18n_tools.jsxinc

i18n.tickThreshold = 2500;

var files = ['Pathfinder - Core', 'Pathfinder - Advanced', 'Pathfinder - Everything', '3.5 - Core', '3.5 - Everything', 'Everything'];
for (var i = 0; i < files.length; i++) {
	var filename = files[i];

	i18n.init();
	i18n.messages = i18n.loadCSV(new File(baseFolder+'Languages/Template/'+filename+'.csv'));
	
	for (var j = 0; j < i18n.messages.length; j++) {
    	i18n.messages[j]['Translation'] = i18n.messages[j]['Frame number'];
    	i18n.tick();
	}
	i18n.saveCSV(new File(baseFolder+'Languages/Numbered/'+filename+'.csv'));
}

alert("Done");