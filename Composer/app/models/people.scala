package models

  import com.mongodb.casbah.Imports._
  import com.novus.salat._
  import com.novus.salat.global._
  import com.novus.salat.dao._
  import Cryptable._

/* Person
 */

object Model {
  val mongo = MongoConnection()("charsheets")
}

object User extends SalatDAO[User, ObjectId](collection = Model.mongo("User")) {
  def get(id: ObjectId) = findOne(MongoDBObject("_id" -> id))
  def get(email: String) = findOne(MongoDBObject("email" -> email))
}

case class User (
  _id: ObjectId,
  email: String,
  pwd: String      // bcrypt hash of their password
) {
  def login(candidate: String): Boolean = pwd.verify(candidate)

  def setPassword(newPassword: String): User = {
    val newUser = this.copy(pwd = newPassword.bcrypt)
    User.save(newUser)
    newUser
  }
/*
  def requestReset(): LoginReset = {
    val secret = // random code
    val reset = new LoginReset(null, _id, DateTime.now(), secret.bcrypt)
    reset.save()
    secret
  }

  def credits: List[Credit] = CreditDAO.get(this)

  def hasCredit: Boolean = {
    val now = DateTime.now
    !credits.filter(c => c.start < now && c.end > now).isEmpty
  }

  def addCredit(): Credit = {
    val existing = credits
    val start = if (existing.isEmpty) DateTime.now() else existing.head.end
    val end = start + 1 year
    val credit = Credit(null, _id, start, end)
    credit.save()
  }

  def pings: List[Ping] = PingDAO.get(this)*/
}


/*
// Credit

object Credit extends SalatDAO[Credit](collection = Model.mongo("Credit")) {

}

case class Credit (
  _id: ObjectId,
  @Key("person") _person: ObjectID,
  start: DateTime,
  end: DateTime
) {
  def person: Person = Person.get(_person)
}


// Login Reset

object LoginReset extends SalatDAO[LoginReset](collection = Model.mongo("LoginReset")) {

}

case class LoginReset (
  _id: ObjectID,
  @Key("person") _person: ObjectID,
  at: DateTime,
  secret: String,    // bcrypt hash of the login reset key
) {
  def person: Person = Person.get(_person)

  def verify(candidate: String): Boolean = secret.verify(candidate)

  def redeem(candidate: String, password: String): Boolean = {
    if (secret.verify(candidate)) {
      // remove this

      // adjust person
      person.setPassword(password)
      true
    } else false
  }
}


// Pings

object Ping extends SalatDAO[Ping](collection = Model.mongo("Ping")) {

}

@Salat trait Ping {
  def person: Person,
  def at: LocalDate
}

case class Ping_v4 (
  _id: ObjectID,
  @Key("person") _person: ObjectID,
  at: LocalDate,
  ip: Int
) {
  def person: Person = Person.get(_person)
}

case class Ping_v6 (
  _id: ObjectID,
  @Key("person") _person: ObjectID,
  at: LocalDate,
  ip: String
) {
  def person: Person = Person.get(_person)
}
*/
