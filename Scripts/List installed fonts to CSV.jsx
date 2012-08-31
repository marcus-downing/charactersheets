/*
List installed fonts to CSV
Creates a new CSV file with a complete list of fonts installed on your system
*/

#include Tools.jsxinc

//var data = [['Font name', 'Font family', 'Style']];
var data = [];
var len = textFonts.length;
for (var i=0; i < len; i++) {
  var line = {
    'Font name': textFonts[i].name, 
    'Family': textFonts[i].family, 
    'Style': textFonts[i].style 
  };
  data.push(line);
}

var file = File.saveDialog("Save CSV file", "*.csv");
data = data.dissociate();
file.writeCSV(data);