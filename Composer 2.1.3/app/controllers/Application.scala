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

  def useLanguages = true

  // index
  def index = Action { Ok(views.html.index()) }

  //  how to
  def howto = Action { Ok(views.html.howto()) }

  //  legal
  def legal = Action { Ok(views.html.legal()) }

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

  def buildPathfinder = Action { Ok(views.html.build(pathfinderData, iconics, iconicSets, logos)) }
  def buildDnd35 = Action { Ok(views.html.build(dnd35Data, iconics, iconicSets, logos)) }
  def buildTest = Action { Ok(views.html.build(testData, iconics, iconicSets, logos)) }

  //  messages

  def leaveMessageForm = Action { Ok(views.html.messageForm()) }
  def leaveMessagePost = Action { request =>
    //  get the data
    val post: Map[String, Seq[String]] = request.body.asFormUrlEncoded.getOrElse(Map.empty)
    val message = post("message").head
    val author = post("author").head
    val email = post("email").head
    val human = post("human").head

    //  send it
    try {
      if (human == "sure") {
        import org.apache.commons.mail.SimpleEmail

        val mail = new SimpleEmail()
        mail.setHostName("localhost")
        mail.addTo("marcus@basingstokeanimesociety.com")
        mail.setFrom("charactersheets@minotaur.cc")
        mail.setSubject("Charactersheets: Message from "+author)
        if (email != "") {
          mail.addReplyTo(email)
          mail.setMsg(message+"\n\n\nFrom: "+email)
        } else {
          mail.setMsg(message)
        }
                    
        println("Sending message: "+author+" <"+email+">: "+message)
        mail.send()
      }

      //  say thank you
      Ok(views.html.messageThanks(message, author))
    } catch {
      case x => 
        x.printStackTrace()
        Ok(views.html.messageError(message, author))
    }
  }

  def iconics: List[IconicImage] = IconicImage.iconics
  def iconicSets: List[IconicSet] = IconicImage.sets

  def getIconic(path: String): Option[IconicImage] = IconicImage.get(path)

  def logos = Logo.logos
}