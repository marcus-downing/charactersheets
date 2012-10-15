package models

object CharacterData {
  def parse(data: Map[String, Seq[String]], gameData: GameData): CharacterData = {
    val data2 = data.map { case (key, list) => key -> list.head }
    //println("Data 2 "+data2)
    val keys = data.keys.toList
    val positive = keys.filter(key => data2.get(key) == Some("on"))

    // classes
    val classNames = positive
      .filter(_.startsWith("class-"))
      .map(_.substring(6))
    val baseClasses: List[BaseClass] = classNames.flatMap(name => gameData.classByName(name)).toList
    
    val classes: List[GameClass] = baseClasses.map { cls =>
      data2.get("variant-"+cls.name) match {
        case Some(variantName) => 
          //println("Variant name "+variantName)
          //println("Variant found "+cls.variantByName(variantName))
          cls.variantByName(variantName).getOrElse(cls)
        case _ => cls
      }
    }
    
    //val classes = baseClasses
    //println("Class names: "+classes.map(_.name).mkString(", "))

    // data
    CharacterData(
      classes, 
      colour = data2.get("colour").getOrElse("normal"),
      spellbookSize = data2.get("spellbook-size").getOrElse("medium"),
      inventoryStyle = data2.get("inventory-style").getOrElse("auto"),
      inventoryIconic = data2.get("inventory-iconic").getOrElse("default"),

      includeGM = positive.contains("gm"),
      partyDownload = positive.contains("party-download"),
      hideInventory = false,
      includeCharacterBackground = positive.contains("include-background"),
      includePartyFunds = positive.contains("include-party-funds"),

      watermark = if (positive.contains("has-watermark")) data2.get("watermark").getOrElse("") else ""
      )
  }
}

case class CharacterData (
  classes: List[GameClass],
  colour: String,
  spellbookSize: String,
  inventoryStyle: String,
  inventoryIconic: String,

  includeGM: Boolean,
  partyDownload: Boolean,
  hideInventory: Boolean,
  includeCharacterBackground: Boolean,
  includePartyFunds: Boolean,

  watermark: String
) {
  def iconic: Option[IconicImage] = controllers.Application.getIconic(inventoryIconic)
}

case class IconicImage(group: String, name: String) {
  val path = group+"/"+name
  val id = (group+"--"+name).toLowerCase.replaceAll("[^a-z]+", " ").trim.replace(" ", "-")
  val largeFile = "public/images/iconics/"+group+"/Large/"+name+".png"
  val smallFile = "images/iconics/"+group+"/Small/"+name+".png"
}