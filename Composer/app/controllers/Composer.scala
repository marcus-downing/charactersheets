package controllers

import play.api._
import play.api.mvc._
import play.api.data.{Form, Mapping}

import java.io.{File,FileInputStream,ByteArrayOutputStream}
import scala.io.Source
import com.itextpdf.text.pdf._
import com.itextpdf.text.{Paragraph, BaseColor, Document, Image, Element}

import models._
import controllers.Application.isAprilFool

object Composer extends Controller {
  lazy val pathfinderData = Application.pathfinderData
  lazy val dnd35Data = Application.dnd35Data
  lazy val testData = Application.testData

  def downloadPathfinder = downloadAction(pathfinderData)
  def downloadDnd35 = downloadAction(dnd35Data)
  def downloadTest = downloadAction(testData)

  def downloadAction(gameData: GameData) = Action(parse.multipartFormData) { request =>
    val iconic = request.body.file("iconic-custom-file").map{ filepart =>
      for (contentType <- filepart.contentType)
        println("File uploaded with content type: "+contentType)
      filepart.ref.file
    }
    val bodydata = request.body.asFormUrlEncoded
    //val data: Map[String, String] = bodydata.flatMap { case (key, list) => key -> list.headOption } toMap
    val data: Map[String, String] = bodydata.mapValues { _.head }
    val sourceFolder = new File("public/pdf/"+gameData.game)

    data.get("start-type") match {
      case Some("single") =>
        val character = CharacterData.parse(data, gameData, iconic)

        val pdf = composePDF(character, gameData, sourceFolder)
        val filename = character.classes.toList.map(_.name).mkString(", ")+".pdf"

        Ok(pdf).as("application/pdf").withHeaders(
          "Content-disposition" -> ("attachment; filename=\""+filename+"\"")
        )

      case Some("party") =>
        val characters = CharacterData.parseParty(data, gameData)
        val pdf = composeParty(characters, gameData, sourceFolder)
        val filename = characters.map(_.classes.toList.map(_.name).mkString("-")).mkString(", ")+".pdf"

        Ok(pdf).as("application/pdf").withHeaders(
          "Content-disposition" -> ("attachment; filename=\""+filename+"\"")
        )

      case Some("gm") =>
        val gmdata = CharacterData.parseGM(data, gameData)
        val pdf = composeGM(gmdata, gameData, sourceFolder)
        Ok(pdf).as("application/pdf").withHeaders(
          "Content-disposition" -> ("attachment; filename=\""+(if(gameData.isDnd35) "Dungeon Master" else "Game Master")+"\"")
        )

      case _ => NotFound
    }
  }

  def composeGM(gmdata: GMData, gameData: GameData, folder: File): Array[Byte] = {
    val out = new ByteArrayOutputStream()
    val document = new Document
    val writer = PdfWriter.getInstance(document, out)
    writer.setRgbTransparencyBlending(true)
    document.open

    val colour = gmdata.colour

    val pages = gameData.gm
    for (page <- pages) {
      val pageFile = new File(folder.getPath+"/"+page.file)
      val fis = new FileInputStream(pageFile)
      val reader = new PdfReader(fis)

      // get the right page size
      val template = writer.getImportedPage(reader, 1)
      val pageSize = reader.getPageSize(1)
      document.setPageSize(pageSize)
      document.newPage

      //  fill with white so the blend has something to work on
      val canvas = writer.getDirectContent
      val baseLayer = new PdfLayer("Character Sheet", writer);
      canvas.beginLayer(baseLayer)
      canvas.setColorFill(BaseColor.WHITE)
      canvas.rectangle(0f, 0f, 1000f, 1000f)
      canvas.fill

      val defaultGstate = new PdfGState
      defaultGstate.setBlendMode(PdfGState.BM_NORMAL)
      defaultGstate.setFillOpacity(1.0f)
      canvas.setGState(defaultGstate)

      //  the page
      canvas.addTemplate(template, 0, 0)
      writeCopyright(canvas, writer, gameData)
      writeColourOverlay(canvas, colour)

      //  done
      canvas.endLayer()
      fis.close
    }
    document.close
    out.toByteArray
  }

  def composeParty(characters: List[CharacterData], gameData: GameData, folder: File): Array[Byte] = {
    val out = new ByteArrayOutputStream()
    val document = new Document
    val writer = PdfWriter.getInstance(document, out)
    writer.setRgbTransparencyBlending(true)
    document.open

    for (character <- characters) {
      println("START OF CHARACTER")
      addCharacterPages(character, gameData, folder, document, writer)
    }

    document.close
    out.toByteArray
  }

  def composePDF(character: CharacterData, gameData: GameData, folder: File): Array[Byte] = {
    val out = new ByteArrayOutputStream()
    val document = new Document
    val writer = PdfWriter.getInstance(document, out)
    writer.setRgbTransparencyBlending(true)
    document.open

    addCharacterPages(character, gameData, folder, document, writer)

    document.close
    out.toByteArray
  }

  def addCharacterPages(character: CharacterData, gameData: GameData, folder: File, document: Document, writer: PdfWriter) {
    val iconic = character.iconic
    val pages = new CharacterInterpretation(gameData, character).pages

    val colour = character.colour
    for (page <- pages) {
      val pageFile = new File(folder.getPath+"/"+page.file)
      val fis = new FileInputStream(pageFile)
      val reader = new PdfReader(fis)

      // get the right page size
      val template = writer.getImportedPage(reader, 1)
      val pageSize = reader.getPageSize(1)
      document.setPageSize(pageSize)
      document.newPage

      //  fill with white so the blend has something to work on
      val canvas = writer.getDirectContent
      val baseLayer = new PdfLayer("Character Sheet", writer);
      canvas.beginLayer(baseLayer)
      canvas.setColorFill(BaseColor.WHITE)
      canvas.rectangle(0f, 0f, 1000f, 1000f)
      canvas.fill

      val defaultGstate = new PdfGState
      defaultGstate.setBlendMode(PdfGState.BM_NORMAL)
      defaultGstate.setFillOpacity(1.0f)
      canvas.setGState(defaultGstate)

      //  the page
      canvas.addTemplate(template, 0, 0)

      //  copyright notice
      writeCopyright(canvas, writer, gameData)

      //  generic image
      if (page.slot == "inventory" && !character.iconic.isDefined && character.customIconic.isEmpty) {
        canvas.setGState(defaultGstate)
        val imgLayer = new PdfLayer("Iconic image", writer)
        canvas.beginLayer(imgLayer)
        val imgFile = "public/images/iconics/generic.png"
        val awtImage = java.awt.Toolkit.getDefaultToolkit().createImage(imgFile)
        val img = Image.getInstance(awtImage, null)
        img.scaleToFit(200f,220f)
        img.setAbsolutePosition(315f - (img.getScaledWidth() / 2), 410f)
        canvas.addImage(img)
        canvas.endLayer
      }

      // april fool
      if (page.slot == "core" && isAprilFool) {
        val pageFile = new File(folder.getPath+"/Extra/Special Overlays/Character Info.pdf")
        val fis = new FileInputStream(pageFile)
        val reader = new PdfReader(fis)
        val template = writer.getImportedPage(reader, 1)
        canvas.setGState(defaultGstate)

        //  the page
        canvas.addTemplate(template, 0, 0)
      }

      //  watermark
      if (character.watermark != "") {
        writeWatermark(canvas, writer, character.watermark)
      }

      writeColourOverlay(canvas, colour)

      canvas.endLayer()

      //  logo
      if (page.slot == "core" || page.slot == "eidolon") {
        canvas.setGState(defaultGstate)
        val imgLayer = new PdfLayer("Logo image", writer)
        canvas.beginLayer(imgLayer)
        val imgFile = logoImage(gameData, character)
        println("Adding logo: "+imgFile)
        val awtImage = java.awt.Toolkit.getDefaultToolkit().createImage(imgFile)
        val img = Image.getInstance(awtImage, null)
        img.scaleToFit(170f,35f)
        img.setAbsolutePosition(127f - (img.getScaledWidth() / 2), 785f)
        canvas.addImage(img)
        canvas.endLayer()
      }

      //  iconics
      if (page.slot == "inventory") {
        if (isAprilFool) {
            println("Adding april fool image")
            canvas.setGState(defaultGstate)
            val imgLayer = new PdfLayer("Iconic image", writer)
            canvas.beginLayer(imgLayer)
            val imgFile = "public/images/iconics/1 Paizo/3 Advanced Races/Large/Goblin - d20.png"
            println("Image: "+imgFile)
            val awtImage = java.awt.Toolkit.getDefaultToolkit().createImage(imgFile)
            val img = Image.getInstance(awtImage, null)
            img.scaleToFit(190f,220f)
            img.setAbsolutePosition(315f - (img.getScaledWidth() / 2), 410f)
            canvas.addImage(img)
            canvas.endLayer()
        } else if (character.iconic.isDefined) {
          for (i <- character.iconic) {
            println("Adding inventory image")
            canvas.setGState(defaultGstate)
            val imgLayer = new PdfLayer("Iconic image", writer)
            canvas.beginLayer(imgLayer)
            val imgFile = i.largeFile
            println("Image: "+imgFile)
            val awtImage = java.awt.Toolkit.getDefaultToolkit().createImage(imgFile)
            val img = Image.getInstance(awtImage, null)
            img.scaleToFit(190f,220f)
            img.setAbsolutePosition(315f - (img.getScaledWidth() / 2), 410f)
            canvas.addImage(img)
            canvas.endLayer()
          }
        } else if (!character.customIconic.isEmpty) {
          for (i <- character.customIconic) {
            println("Adding custom inventory image")
            canvas.setGState(defaultGstate)
            val imgLayer = new PdfLayer("Custom iconic image", writer)
            canvas.beginLayer(imgLayer)
            val filename = i.getAbsolutePath
            println("Image: "+filename)
            val awtImage = java.awt.Toolkit.getDefaultToolkit().createImage(filename)
            val img = Image.getInstance(awtImage, null)
            img.scaleToFit(180f,220f)
            img.setAbsolutePosition(315f - (img.getScaledWidth() / 2), 410f)
            canvas.addImage(img)
            canvas.endLayer()
          }
        }
      }
    fis.close
    }
  }

  def writeCopyright(canvas: PdfContentByte, writer: PdfWriter, gameData: GameData) {
    //  copyright notice
    canvas.setColorFill(new BaseColor(0.5f, 0.5f, 0.5f))
    val font = BaseFont.createFont(BaseFont.HELVETICA, BaseFont.CP1252, BaseFont.EMBEDDED)

    canvas.beginText
    val copyrightLayer = new PdfLayer("Iconic image", writer)
    canvas.beginLayer(copyrightLayer)
    canvas.setFontAndSize(font, 5)
    canvas.showTextAligned(Element.ALIGN_LEFT, "Copyright \u00A9 Marcus Downing 2012        http://charactersheets.minotaur.cc", 30, 21, 0)
    if (gameData.isPathfinder) {
      canvas.setFontAndSize(font, 4)

      canvas.showTextAligned(Element.ALIGN_LEFT, "This character sheet uses trademarks and/or copyrights owned by Paizo Publishing, LLC, which are used under Paizo's Community Use Policy. We are expressly prohibited from charging you to use or", 206, 21, 0)
      canvas.showTextAligned(Element.ALIGN_LEFT, "access this content. This character sheet is not published, endorsed, or specifically approved by Paizo Publishing. For more information about Paizo's Community Use Policy, please visit paizo.com/communityuse. For more information about Paizo Publishing and Paizo products, please visit paizo.com.", 30, 16, 0)
    } else if (gameData.isDnd35) {
      canvas.setFontAndSize(font, 4)

      canvas.showTextAligned(Element.ALIGN_LEFT, "This character sheet is not affiliated with, endorsed, sponsored, or specifically approved by Wizards of the Coast LLC. This character sheet may use the trademarks and other intellectual property of", 206, 21, 0)
      canvas.showTextAligned(Element.ALIGN_LEFT, "Wizards of the Coast LLC, which is permitted under Wizards' Fan Site Policy. For example, DUNGEONS & DRAGONS®, D&D®, PLAYER'S HANDBOOK 2®, and DUNGEON MASTER'S GUIDE® are trademark[s] of Wizards of the Coast and D&D® core rules, game mechanics, characters and their distinctive likenesses are the", 30, 16, 0)
      canvas.showTextAligned(Element.ALIGN_LEFT, "property of the Wizards of the Coast. For more information about Wizards of the Coast or any of Wizards' trademarks or other intellectual property, please visit their website.", 30, 11, 0)
    }
    canvas.endLayer
    canvas.endText
  }

  def writeWatermark(canvas: PdfContentByte, writer: PdfWriter, watermark: String) {
    println("Adding watermark: "+watermark)

    val watermarkGstate = new PdfGState
    watermarkGstate.setBlendMode(PdfGState.BM_NORMAL)
    watermarkGstate.setFillOpacity(0.1f)
    canvas.setGState(watermarkGstate)

    canvas.beginText
    val watermarkLayer = new PdfLayer("Watermark", writer)
    canvas.beginLayer(watermarkLayer)
    val font = BaseFont.createFont(BaseFont.HELVETICA, BaseFont.CP1252, BaseFont.EMBEDDED)
    canvas.setFontAndSize(font, (900f / watermark.length).toInt)
    canvas.setColorFill(new BaseColor(0f, 0f, 0f))
    canvas.showTextAligned(Element.ALIGN_CENTER, watermark, 365f, 400f, 60f)
    canvas.endLayer
    canvas.endText
  }

  def writeColourOverlay(canvas: PdfContentByte, colour: String) {
    if (colour == "black") {
      val gstate = new PdfGState
      
      gstate.setBlendMode(PdfGState.BM_OVERLAY)
      //gstate.setFillOpacity(0.5f)
      canvas.setGState(gstate)
      canvas.setColorFill(new BaseColor(0.1f, 0.1f, 0.1f))
      canvas.rectangle(0f, 0f, 1000f, 1000f)
      canvas.fill
      
      val gstate2 = new PdfGState
      gstate2.setBlendMode(PdfGState.BM_COLORDODGE)
      gstate2.setFillOpacity(0.5f)
      canvas.setGState(gstate2)
      canvas.setColorFill(new BaseColor(0.2f, 0.2f, 0.2f))
      canvas.rectangle(0f, 0f, 1000f, 1000f)
      canvas.fill
      
      //  ...correct hilights...
    } else if (colour != "normal") {
      val gstate = new PdfGState
      gstate.setBlendMode(colour match {
          case "light" => PdfGState.BM_SCREEN
          case "dark" => PdfGState.BM_OVERLAY
          case "black" => PdfGState.BM_COLORBURN
          case _ => new PdfName("Color")
      })
      canvas.setGState(gstate)
      canvas.setColorFill(interpretColour(colour))
      canvas.rectangle(0f, 0f, 1000f, 1000f)
      canvas.fill
    }
  }

  def interpretColour(colour: String): BaseColor = colour match {
    case "light" => new BaseColor(0.3f, 0.3f, 0.3f)
    case "dark" => new BaseColor(0.35f, 0.35f, 0.35f)
    case "black" => new BaseColor(0f, 0f, 0f)
    case "red" => new BaseColor(0.60f, 0.2f, 0.2f)
    case "orange" => new BaseColor(0.72f, 0.47f, 0.30f)
    case "yellow" => new BaseColor(1.0f, 0.92f, 0.55f)
    case "lime" => new BaseColor(0.77f, 0.85f, 0.55f)
    case "green" => new BaseColor(0.5f, 0.7f, 0.5f)
    case "cyan" => new BaseColor(0.6f, 0.75f, 0.75f)
    case "blue" => new BaseColor(0.55f, 0.63f, 0.80f)
    case "purple" => new BaseColor(0.80f, 0.6f, 0.70f)
    case "pink" => new BaseColor(1.0f, 0.60f, 0.65f)
  }

  def logoImage(gameData: GameData, character: CharacterData) = gameData.game match {
    case "pathfinder" =>
      if (character.classes.exists(_.pages.exists(_.startsWith("core/neoexodus"))))
        "public/images/neoexodus-logo.png"
      else
        "public/images/pathfinder-logo.png"
    case "dnd35" => "public/images/dnd35-logo.png"
    case _ => ""
  }
}

class CharacterInterpretation(gameData: GameData, character: CharacterData) {
  case class PageSlot(slot: String, variant: Option[String]) {
    lazy val page: Option[Page] = {
      val ps = gameData.pages.toList.filter { p => p.slot == slot && p.variant == variant }
      ps.headOption
    }
    override def toString = variant match {
      case Some(v) => slot+" / "+v 
      case None => slot
    }
  }

  def pageSlot(name: String) = 
    name.split("/", 2).toList match {
      case page :: Nil => PageSlot(page, None)
      case page :: variant :: _ => PageSlot(page, Some(variant))
      case _ => throw new Exception("Wow. I guess that match really wasn't exhaustive.")
    }

  def selectCharacterPages(classes: List[GameClass]): List[Page] = {
    //println(" -- Classes: "+classes.map(_.name).mkString(", "))
    val basePages = gameData.base.pages.toList.map(pageSlot)
    val classPages = classes.flatMap(_.pages).map(pageSlot)

    //  additional pages
    var pages = basePages ::: classPages
    if (character.includeCharacterBackground)
      pages = pages ::: List(PageSlot("background", None))
    if (character.includePartyFunds)
      pages = pages ::: List(PageSlot("partyfunds", None))

    println(" -- Base pages: "+basePages.map(_.toString).mkString(", "))
    println(" -- Class pages: "+classPages.map(_.toString).mkString(", "))
    var slotNames = pages.map(_.slot).distinct
    println(" -- Distinct slots: "+slotNames.mkString(", "))

    //  special cases
    if (character.hideInventory) {
      pages = PageSlot("core", Some("simple")) :: PageSlot("combat", Some("simple")) :: pages
      println("Slot names (before): "+slotNames.mkString(", "))
      slotNames = slotNames.filter(_ != "inventory")
      println("Slot names (simplified): "+slotNames.mkString(", "))
    }

    if (slotNames.contains("spellbook")) {
      val spellbookPage = character.spellbookSize match {
        case "small" => PageSlot("spellbook", Some("small"))
        case "medium" => PageSlot("spellbook", None)
        case "large" => PageSlot("spellbook", Some("large"))
      }
      pages = pages.filter(_.page != "spellbook") ::: (spellbookPage :: Nil)
    }
    if (character.inventoryStyle != "auto") {
      val page = PageSlot("inventory", Some(character.inventoryStyle))
      if (page.page != None)
        pages = pages.filter(_.page != "inventory") ::: (page :: Nil)
    }

    println(" -- Final slots: "+slotNames.mkString(", "))
    pages = for (slotName <- slotNames) yield {
      val pageInstances = pages.filter( _.slot == slotName)
      val overridingInstances = pageInstances.filter(v => v.variant != None)
      val selectedInstance =
        if (overridingInstances == Nil)
          pageInstances.head
        else
          overridingInstances.head

      println("Page: "+slotName+" ~ " + selectedInstance.variant.getOrElse(""))
      selectedInstance
    }
    
    println(" -- Selected pages: "+pages.map(_.toString).mkString(", "))
    val printPages = pages.toList.flatMap(_.page)
    printPages
    //printPages.sortWith((a,b) => a.pagePosition < b.pagePosition)
  }

  def pages = {
    var clsPages =
      if (character.partyDownload)
        character.classes.flatMap( cls => selectCharacterPages(List(cls)) )
      else
        selectCharacterPages(character.classes)

    var pages = 
      if (character.includeGM) 
        clsPages ::: gameData.gm
      else
        clsPages

    pages
  }
}
