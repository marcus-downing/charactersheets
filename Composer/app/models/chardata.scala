package models

import java.io.File

object CharacterData {

  def positiveData(data: Map[String, String]): List[String] = {
    val keys = data.keys.toList
    keys.filter(key => data.get(key) == Some("on"))
  }

  def parse(data: Map[String, String], gameData: GameData, customIconic: Option[File]): CharacterData = {
    //val data2 = data.flatMap { case (key, list) => key -> list.headOption }
    //println("Data 2 "+data2)
    val positive = positiveData(data)

    // classes
    val classNames = positive
      .filter(_.startsWith("class-"))
      .map(_.substring(6))
    val baseClasses: List[BaseClass] = classNames.flatMap(name => gameData.classByName(name)).toList
    
    val classes: List[GameClass] = baseClasses.map { cls =>
      data.get("variant-"+cls.name) match {
        case Some(variantName) => 
          //println("Variant name "+variantName)
          //println("Variant found "+cls.variantByName(variantName))
          cls.variantByName(variantName).getOrElse(cls)
        case _ => cls
      }
    }

    // data
    CharacterData(
      classes, 
      colour = data.get("colour").getOrElse("normal"),
      spellbookSize = data.get("spellbook-size").getOrElse("medium"),
      inventoryStyle = data.get("inventory-style").getOrElse("auto"),
      inventoryIconic = data.get("inventory-iconic").getOrElse("default"),
      customIconic = customIconic,

      includeGM = positive.contains("gm"),
      partyDownload = positive.contains("party-download"),
      hideInventory = positive.contains("simple"),
      includeCharacterBackground = positive.contains("include-background"),
      includePartyFunds = positive.contains("include-party-funds"),

      watermark = if (positive.contains("has-watermark")) data.get("watermark").getOrElse("") else ""
      )
  }

  def parseParty(data: Map[String, String], gameData: GameData): List[CharacterData] = {
    val charids = data.get("charids").getOrElse("").split(",").map(_.trim).filter(_ != "").toList
    println("Stashed character IDs: "+charids.mkString(", "))
    val stashedCharacters = charids.map { charid =>
      val prefix = "char-"+charid+"-"
      val subdata: Map[String, String] = data.filterKeys(_.startsWith(prefix)).map { case (key, value) => key.substring(prefix.length) -> value } toMap;
      parse(subdata, gameData, None)
    }
    val finalCharacter = parse(data, gameData, None)
    stashedCharacters ::: (finalCharacter :: Nil)
  }

  def parseGM(data: Map[String, String], gameData: GameData): GMData = {
    val positive = positiveData(data)

    GMData(
      colour = data.get("colour").getOrElse("normal"),
      watermark = if (positive.contains("has-watermark")) data.get("watermark").getOrElse("") else ""
      )
  }
}

case class GMData (
  colour: String,
  watermark: String
  )

case class CharacterData (
  classes: List[GameClass],
  colour: String,
  spellbookSize: String,
  inventoryStyle: String,
  inventoryIconic: String,
  customIconic: Option[File],

  includeGM: Boolean,
  partyDownload: Boolean,
  hideInventory: Boolean,
  includeCharacterBackground: Boolean,
  includePartyFunds: Boolean,

  watermark: String
) {
  def iconic: Option[IconicImage] = controllers.Application.getIconic(inventoryIconic)
}

case class IconicImage(group: String, set: String, name: String) {
  import IconicImage.slug
  val path = group+"/"+set+"/"+name
  val id = slug(group)+"--"+slug(set)+"--"+slug(name)
  val largeFile = "public/images/iconics/"+group+"/"+set+"/Large/"+name+".png"
  val smallFile = "images/iconics/"+group+"/"+set+"/Small/"+name+".png"
}

object IconicImage {
  def withoutNumber(name: String): String = {
    val rex = """[0-9]+\s+(.*)""" r
    
    name match {
      case rex(rem) => rem
      case _ => name
    }
  }

  def slug(str: String): String = str.toLowerCase.replaceAll("[^a-z]+", "-")
}