package models

case class Quote (quote: String, author: String) {
  def noteClass = if (quote.length > 115) "long" else if (quote.length > 50) "medium" else "short"
}