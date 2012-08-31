/*
List installed fonts to AI
Creates a new A3 sized document and display a list of available fonts until the document is full.
Almost direct copy of the Illustrator Scripting Reference, p216
*/

var edgeSpacing = 10;
var columnSpacing = 230;
var docPreset = new DocumentPreset;
docPreset.width = 1191.0;
docPreset.height = 842.0
var doc = documents.addDocument(DocumentColorSpace.CMYK, docPreset);
var sFontNames = "";
var x = edgeSpacing;
var y = (doc.height - edgeSpacing);

var len = textFonts.length;
for (var i=0; i < len; i++) {
  sFontNames = textFonts[i].family + " / " + textFonts[i].style;
  var textRef = doc.textFrames.add();
  textRef.textRange.characterAttributes.size = 10;
  textRef.contents = sFontNames;
  textRef.top = y;
  textRef.left = x;
  
  // check wether the text frame will go off the edge of the document
  if ((x + textRef.width)> doc.width){
    textRef.remove();
    len = i;
    break;
  } else{
    // display text frame
    textRef.textRange.characterAttributes.textFont =
    textFonts.getByName(textFonts[i].name);
    redraw();
    if( (y-=(textRef.height)) <= 20 ) {
      y = (doc.height - edgeSpacing);
      x += columnSpacing;
    }
  }
}