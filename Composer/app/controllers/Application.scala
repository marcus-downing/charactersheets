package controllers

import play.api._
import play.api.mvc._

import java.io.File
import scala.io.Source
import org.joda.time._

import models._

object Application extends Controller {
  
  def isAprilFool = {
    val today = new DateTime().toLocalDate
    //println("Today is "+today.getDayOfMonth+" of "+today.getMonthOfYear)
    today.getDayOfMonth == 1 && today.getMonthOfYear == 4
    //true
  }

  // index
  def index = Action { Ok(views.html.index()) }

  // quotes
  lazy val quotes: List[Quote] = {
    val quotesFile = new File("public/quotes.txt")
    println(" * Quotes file at: "+quotesFile.getAbsolutePath())
    if (quotesFile.exists()) {
      val lines = Source.fromFile(quotesFile).getLines().toList
      val quoteLines = lines.map(_.trim).filter(s => s != "" && !s.startsWith("--") && s.indexOf(" --by-- ") != -1)
      //println(" * File has "+quoteLines.length+" quotes")

      val quotes = quoteLines.map { line => 
        //println(" * - "+line)
        val parts = line.split(" --by-- ")
        val quote = parts(0)
        val author = parts(1)
        //println(" * - [ "+quote+" ] by "+author)
        Quote(quote, author)
      }
      println(" * Made "+quotes.length+" quotes")
      quotes.filter(_.quote.length <= 200).toList
    } else Nil
  }

  def randomQuote: Quote = quotes(scala.util.Random.nextInt(quotes.length))

  //  build
  lazy val pathfinderData: GameData = GameData.load("pathfinder")
  lazy val dnd35Data: GameData = GameData.load("dnd35")
  lazy val testData: GameData = GameData.load("test")

  def buildPathfinder = Action { Ok(views.html.build(pathfinderData, iconics)) }
  def buildDnd35 = Action { Ok(views.html.build(dnd35Data, iconics)) }
  def buildTest = Action { Ok(views.html.build(testData, iconics)) }

  //  Group -> Set -> [] IconicImage
  def iconics: Map[String, Map[String, List[IconicImage]]] = {
    val iconicsFolder = new File("public/images/iconics")
    if (!iconicsFolder.isDirectory) return Map.empty

    val groups: List[(String, Map[String, List[IconicImage]])] = for (groupFolder <- iconicsFolder.listFiles.toList if groupFolder.isDirectory) yield {
      val group = groupFolder.getName
      val groupName = IconicImage.withoutNumber(group)

      val sets: List[(String, List[IconicImage])] = for (setFolder <- groupFolder.listFiles.toList if setFolder.isDirectory) yield {
        val set = setFolder.getName
        val setName = IconicImage.withoutNumber(set)

        val largeFolder = new File(setFolder.getAbsolutePath+"/Large")
        val imageFiles: List[IconicImage] = largeFolder.listFiles.toList.filter(_.isFile).flatMap { file =>
          val filename = file.getName
          if (filename.endsWith(".png")) {
            val imageName = filename.substring(0, filename.length - 4)
            Some(IconicImage(group, set, imageName))
          }
          else None
        }
        (setName, imageFiles)
      }
      (groupName, sets.toMap)
    }
    groups.toMap
  }

  def getIconic(path: String): Option[IconicImage] = {
    path.split("/").toList match {
      case group :: set :: name :: Nil => Some(IconicImage(group, set, name))
      case _ => None
    }
  }
}