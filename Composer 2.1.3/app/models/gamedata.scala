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
  languages: List[LanguageInfo],
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
  aps: List[AP] = Nil
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
  name: String,
  ready: List[Float]
)