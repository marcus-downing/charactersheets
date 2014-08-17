package models

import java.io.File
import scala.io.Source
// import com.novus.salat._
// import com.novus.salat.global._
import play.api.libs.json._

object GameData {
  def load(game: String): GameData = {
    val file = new File("public/data/"+game+".json")
    val data = Source.fromFile(file).getLines().mkString
    val json = Json.parse(data)
    parse(json)
  }

  def parse(json: JsValue) = GameData(
      game = (json \ "game").as[String],
      name = (json \ "name").as[String],
      pages = (json \ "pages").as[List[JsObject]].map(parsePage),
      gm = parseGM((json \ "gm").as[JsObject]),
      base = parseBaseData((json \ "base").as[JsObject]),
      layout = (json \ "layout").as[List[List[String]]],
      books = (json \ "books").as[List[JsObject]].map(parseBook),
      languages = (json \ "languages").as[List[JsObject]].map(parseLanguageInfo),
      classes = (json \ "classes").as[List[JsObject]].map(parseBaseClass)
    )

  def parsePage(json: JsObject) = Page(
    file = (json \ "file").as[String],
    page = (json \ "page").asOpt[Int].getOrElse(1),
    slot = (json \ "slot").asOpt[String].getOrElse(""),
    name = (json \ "name").asOpt[String].getOrElse(""),
    variant = (json \ "variant").asOpt[String],
    position = (json \ "position").asOpt[Int]
  )

  def parseGM(json: JsObject) = GM(
    campaign = (json \ "campaign").as[List[JsObject]].map(parsePage),
    maps = parseMaps((json \ "maps").as[JsObject]),
    aps = (json \ "aps").asOpt[List[JsObject]].getOrElse(Nil).map(parseAP)
  )

  def parseMaps(json: JsObject) = Maps(
    maps2d = (json \ "2d").as[List[JsObject]].map(parsePage),
    maps3d = (json \ "3d").as[List[JsObject]].map(parsePage)
  )

  def parseAP(json: JsObject) = AP(
    name = (json \ "name").as[String],
    code = (json \ "code").as[String],
    pages = (json \ "pages").as[List[JsObject]].map(parsePage)
  )

  def parseBaseData(json: JsObject) = BaseData(
    pages = (json \ "pages").as[List[String]]
  )

  def parseBook(json: JsObject) = Book(
    name = (json \ "name").as[String],
    classes = (json \ "classes").as[List[String]]   
  )

  def parseLanguageInfo(json: JsObject) = LanguageInfo(
    code = (json \ "code").as[String],
    short = (json \ "short").as[String],
    name = (json \ "name").as[String],
    ready = (json \ "ready").as[List[Float]]
  )

  def parseBaseClass(json: JsObject) = BaseClass(
    name = (json \ "name").as[String],
    pages = (json \ "pages").as[List[String]],
    variants = (json \ "variants").asOpt[List[JsObject]].getOrElse(Nil).map(parseVariant),
    axes = (json \ "axes").asOpt[List[List[String]]].getOrElse(Nil)
  )

  def parseVariant(json: JsObject) = VariantClass(
    name = (json \ "name").as[String],
    pages = (json \ "pages").as[List[String]],
    axes = (json \ "axes").asOpt[List[String]].getOrElse(Nil)
  )
}

case class GameData (
  game: String,
  name: String,
  pages: List[Page],
  gm: GM,
  base: BaseData,
  layout: List[List[String]],
  books: List[Book],
  languages: List[LanguageInfo],
  classes: List[BaseClass]
) {
  def isPathfinder = game == "pathfinder"
  def isDnd35 = game == "dnd35"
  def isNeoexodus = game == "neoexodus"
  def isTest = game == "test"
  def classByName(name: String) = classes.filter(_.name == name).headOption
  def bookByName(name: String) = books.filter(_.name == name).headOption

  def slugOf(str: String) = str.toLowerCase.replaceAll("[^a-z]+", " ").trim.replace(" ", "-")
}

case class GM (
  campaign: List[Page],
  maps: Maps,
  aps: List[AP] = Nil
  )

case class Maps (
  maps2d: List[Page],
  maps3d: List[Page]
  )

case class AP (
  name: String,
  code: String,
  pages: List[Page]
  )

case class Page (
  file: String,
  page: Int = 1,
  slot: String = "",
  name: String = "",
  variant: Option[String] = None,
  position: Option[Int] = None
) {
  def pagePosition = position.getOrElse(page)
}

case class BaseData (
  pages: List[String]
)

case class Book (
  name: String,
  classes: List[String]
)

trait GameClass {
  def name: String
  def pages: List[String]
  def code = name.replaceAll("[^a-zA-Z]+", "-")
}

case class BaseClass (
  name: String,
  pages: List[String],
  variants: List[VariantClass] = Nil,
  axes: List[List[String]] = Nil
) extends GameClass {
  def variantByName(name: String): Option[GameClass] = variants.filter(_.name == name).map(_.mergeInto(this)).headOption
  def axisValues: List[List[String]] = axes.zipWithIndex.map { case (axisValues,index) =>
    if (!axisValues.isEmpty) axisValues 
    else variants.map(_.axes(index)).distinct
  }
  def variantByAxes(axisValues: List[String]): Option[GameClass] = variants.filter(_.axes == axisValues).map(_.mergeInto(this)).headOption
}

case class VariantClass (
  name: String,
  pages: List[String],
  axes: List[String] = Nil
) extends GameClass {
  def mergeInto(base: BaseClass) = new BaseClass(name, base.pages ::: pages)
}

case class LanguageInfo (
  code: String,
  short: String,
  name: String,
  ready: List[Float]
)