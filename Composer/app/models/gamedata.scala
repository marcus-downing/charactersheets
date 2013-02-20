package models

import java.io.File
import scala.io.Source
import com.novus.salat._
import com.novus.salat.global._

object GameData {
  def load(game: String): GameData = {
    val file = new File("public/data/"+game+".json")
    val data = Source.fromFile(file).getLines().mkString
    grater[GameData].fromJSON(data)
  }
}

case class GameData (
  game: String,
  name: String,
  pages: List[Page],
  gm: GM,
  base: BaseData,
  layout: List[List[String]],
  books: List[Book],
  classes: List[BaseClass],
  iconics: List[String] = Nil
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
  maps: List[Page],
  apKingmaker: List[Page] = Nil,
  apSkullAndShackles: List[Page] = Nil
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
  variants: List[VariantClass] = Nil
) extends GameClass {
  def variantByName(name: String): Option[GameClass] = variants.filter(_.name == name).map(_.mergeInto(this)).headOption
}

case class VariantClass (
  name: String,
  pages: List[String]
) extends GameClass {
  def mergeInto(base: BaseClass) = new BaseClass(name, base.pages ::: pages)
}