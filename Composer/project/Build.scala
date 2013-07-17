import sbt._
import Keys._
import play.Project._

object ApplicationBuild extends Build {

    val appName         = "Character Sheet Composer"
    val appVersion      = "2.0"

    //val novusRels = "repo.novus rels" at "http://repo.novus.com/releases/"
    //val jerksonRels = "repo.codahale.com" at "http://repo.codahale.com"

    // val typesafe = "Typesafe Repo" at "http://repo.typesafe.com/typesafe/releases/"
    // val typesafeSnapshots = "Typesafe Snaps Repo" at "http://repo.typesafe.com/typesafe/snapshots/"
    
    val appDependencies = Seq(
      // Add your project dependencies here
      //"org.mindrot" % "jbcrypt" % "0.3m",
      "joda-time" % "joda-time" % "2.2",
      //"com.mongodb.casbah" %% "casbah" % "2.1.5-1",
      "org.mongodb" %% "casbah" % "2.6.2",
      "com.novus" %% "salat-core" % "1.9.1",
      //"com.lowagie" % "itext" % "2.1.7"
      "com.itextpdf" % "itextpdf" % "5.4.2",
      "org.apache.commons" % "commons-email" % "1.3.1"
      //"net.debasishg" %% "sjson" % "0.17"
      //"com.codahale" %% "jerkson" % "0.5.0"
    )

    val main = play.Project(appName, appVersion, appDependencies).settings(
      // Add your own project settings here      
    )

}
