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
      if (cls.axes.isEmpty)
        data.get("variant-"+cls.name) match {
          case Some(variantName) => 
            //println("Variant name "+variantName)
            //println("Variant found "+cls.variantByName(variantName))
            cls.variantByName(variantName).getOrElse(cls)
          case _ => cls
        }
      else {
        val axisValues = (Range(0, cls.axes.length) flatMap { i => data.get("variant-"+cls.name+"-axis-"+i) }).toList
        println("Axis values: "+axisValues.mkString(", "))
        cls.variantByAxes(axisValues).getOrElse(cls)
      }
    }

    val variantRules: List[String] = positive.filter(_.startsWith("variant-")).map(_.substring("variant-".length))
    val inventoryIconic = data.get("inventory-iconic").getOrElse("default")

    println("Given iconic: "+data.get("inventory-iconic").getOrElse("(none)"))
    if (customIconic.isDefined) println("Custom iconic uploaded")

    // data
    CharacterData(
      classes, 
      colour = data.get("colour").getOrElse("normal"),
      spellbookSize = data.get("spellbook-size").getOrElse("medium"),
      inventoryStyle = data.get("inventory-style").getOrElse("auto"),
      inventoryIconic = inventoryIconic,
      customIconic = if (inventoryIconic == "custom") customIconic else None,
      logo = Logo.get(data.get("logo").getOrElse(gameData.game)),

      includeGM = positive.contains("gm"),
      partyDownload = positive.contains("party-download"),
      hideInventory = positive.contains("simple"),
      includeCharacterBackground = positive.contains("include-background"),
      includePartyFunds = positive.contains("include-party-funds"),
      includeAnimalCompanion = positive.contains("include-animal-companion"),

      watermark = if (positive.contains("has-watermark")) data.get("watermark").getOrElse("") else "",

      variantRules = variantRules
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

    val aps = for (ap <- gameData.gm.aps; if positive.contains("ap-"+ap.code)) yield ap.code
    println("Game APs: "+gameData.gm.aps.map(_.code).mkString(", "))

    GMData(
      colour = data.get("colour").getOrElse("normal"),
      watermark = if (positive.contains("has-watermark")) data.get("watermark").getOrElse("") else "",
      gmCampaign = positive.contains("gm-campaign"),
      maps = positive.contains("maps"),
      maps3d = data.get("maps-view").getOrElse("3d") == "3d",
      aps = aps
      )
  }
}

case class GMData (
  colour: String,
  watermark: String,
  gmCampaign: Boolean,
  maps: Boolean,
  maps3d: Boolean,
  aps: List[String]
  )

case class CharacterData (
  classes: List[GameClass],
  colour: String,
  spellbookSize: String,
  inventoryStyle: String,
  inventoryIconic: String,
  customIconic: Option[File],
  logo: Option[Logo],

  includeGM: Boolean,
  partyDownload: Boolean,
  hideInventory: Boolean,
  includeCharacterBackground: Boolean,
  includePartyFunds: Boolean,
  includeAnimalCompanion: Boolean,

  watermark: String,

  variantRules: List[String]
) {
  def hasCustomIconic = customIconic.isDefined
  def iconic: Option[IconicImage] = IconicImage.get(inventoryIconic)
}


//  Iconics

case class IconicSet(filePath: String, nicePath: String) {
  val sortableName = filePath
  val (groupName, setName) = IconicImage.splitPath(nicePath)
  val id = IconicImage.slug(filePath)
  val groupDisplayName = IconicImage.withoutNumber(groupName)
  val setDisplayName = IconicImage.withoutNumber(setName)

  lazy val iconics: List[IconicImage] = IconicImage.iconics.filter(_.set == this).sortBy(_.sortableName)
}

case class IconicImage(set: IconicSet, fileName: String, niceName: String) {
  import IconicImage.slug
  val path = set.filePath+"/"+fileName
  val id = set.id+"-"+slug(fileName)
  val sortableName = id
  val largeFile = "public/images/iconics/large/"+set.filePath+"/"+fileName+".png"
  val smallFile = "public/images/iconics/large/"+set.filePath+"/"+fileName+".png"
  val url = ("/images/iconics/small/"+set.filePath+"/"+fileName+".png").replaceAll(" ", "+")
}

object IconicImage {
  lazy val iconics: List[IconicImage] = {
    val iconicsFolder = new File("public/images/iconics")
    if (!iconicsFolder.isDirectory) Nil
    else {
      val iconicsList = new File("public/images/iconics/iconics.txt")
      val lines = scala.io.Source.fromFile(iconicsList).getLines.toList
      val iconics = lines.flatMap { line =>
        try {
          val filePath :: nicePath :: _ = line.split("=").toList
          val (fileBase, fileName) = splitPath(filePath)
          val (niceBase, niceName) = splitPath(nicePath)
          val set = IconicSet(fileBase, niceBase)
          println(" - Found iconic: "+set.nicePath+" / "+niceName)
          Some(IconicImage(set, fileName, niceName))
        } catch {
          case _: Exception => None
        }
      }

      println("Found "+iconics.length+" iconics")

      iconics.sortBy(_.sortableName)
    }
  }

  def get(code: String): Option[IconicImage] = iconics.filter(_.id == code).headOption

  lazy val sets: List[IconicSet] = iconics.map(_.set).distinct.sortBy(_.sortableName)

  def withoutNumber(name: String): String = {
    val rex = """[0-9]+\s+(.*)""" r
    
    name match {
      case rex(rem) => rem
      case _ => name
    }
  }

  def splitPath(path: String): (String, String) = {
    val reversePath: List[String] = path.split("/").toList.reverse
    val head = reversePath.head
    val tail = reversePath.tail.reverse
    val base = tail.mkString("/")
    (base, head)
  }

  def slug(str: String): String = str.toLowerCase.replaceAll("/", "--").replaceAll("[^a-z0-9/]+", "-").replaceAll("(\\.|-)png$", "").replaceAll("^-+", "").replaceAll("-+$", "")
}

case class Logo(code: String, name: String) {
  lazy val fileName: Option[String] = {
    if (new File("public/images/logos/"+code+".png").exists()) Some(code+".png")
    else if (new File("public/images/logos/"+code+".jpg").exists()) Some(code+".jpg")
    else None
  }
  def url = fileName.map("/images/logos/"+_)
  def filePath = fileName.map("public/images/logos/"+_)
  def ifExists: Option[Logo] = fileName.map(f => this)
}

object Logo {
  lazy val logos: List[Logo] = {
    val logosFolder = new File("public/images/logos")
    if (!logosFolder.isDirectory) Nil
    else {
      val logosList = new File("public/images/logos/logos.txt")
      val lines = scala.io.Source.fromFile(logosList).getLines.toList

      val logos = lines.flatMap { line =>
        try {
          val code :: name :: _ = line.split("=").toList
          println(" - Found logo: "+code+" / "+name)
          Logo(code, name).ifExists
        } catch {
          case _: Exception => None
        }
      }

      println("Found "+logos.length+" logos")
      logos
    }
  }

  def get(code: String): Option[Logo] = logos.filter(_.code == code).headOption
  def get(codes: List[String]): Option[Logo] = codes.flatMap(code => logos.filter(_.code == code)).headOption
}
