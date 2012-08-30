package models

  import org.mindrot.jbcrypt.BCrypt

object Cryptable {
  implicit def string2cryptable(str: String) = new Cryptable(str)
}

class Cryptable(str: String) {
  def bcrypt: String = BCrypt.hashpw(str, BCrypt.gensalt(12))
  def verify(candidate: String): Boolean = BCrypt.checkpw(candidate, str)
}