import sbt._
import Keys._
import PlayProject._

object ApplicationBuild extends Build {

    val appName         = "Character Sheet Composer"
    val appVersion      = "2.0"

    //val novusRels = "repo.novus rels" at "http://repo.novus.com/releases/"
    val typesafe = "Typesafe Repo" at "http://repo.typesafe.com/typesafe/releases/"
    val typesafeSnapshots = "Typesafe Snaps Repo" at "http://repo.typesafe.com/typesafe/snapshots/"
    //val jerksonRels = "repo.codahale.com" at "http://repo.codahale.com"
    
    val appDependencies = Seq(
      // Add your project dependencies here
      "org.mindrot" % "jbcrypt" % "0.3m",
      "joda-time" % "joda-time" % "2.1",
      "com.mongodb.casbah" %% "casbah" % "2.1.5-1",
      "com.novus" %% "salat-core" % "1.9.0",
      //"com.lowagie" % "itext" % "2.1.7"
      "com.itextpdf" % "itextpdf" % "5.2.0"
      //"net.debasishg" %% "sjson" % "0.17"
      //"com.codahale" %% "jerkson" % "0.5.0"
    )

    val main = PlayProject(appName, appVersion, appDependencies, mainLang = SCALA).settings(
      // Add your own project settings here      
    )

}
