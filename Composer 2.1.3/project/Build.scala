import sbt._
import Keys._
import play.Project._

object ApplicationBuild extends Build {

  val appName         = "charactersheets"
  val appVersion      = "2.1.3"

  val appDependencies = Seq(
    // Add your project dependencies here,
    // jdbc,
    // anorm

      "joda-time" % "joda-time" % "2.3",
      //"com.mongodb.casbah" %% "casbah" % "2.1.5-1",
      // "org.mongodb" %% "casbah" % "2.6.2",
      // "com.novus" %% "salat-core" % "1.9.1",
      // "com.novus" %% "salat-core" % "1.9.2-SNAPSHOT",
      //"com.lowagie" % "itext" % "2.1.7"
      "com.itextpdf" % "itextpdf" % "5.4.3",
      "org.apache.commons" % "commons-email" % "1.2"
  )


  val main = play.Project(appName, appVersion, appDependencies).settings(
    // Add your own project settings here      
    resolvers += Resolver.sonatypeRepo("snapshots")
  )

}
