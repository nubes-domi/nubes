class User < ApplicationRecord
  include PrettyId
  self.id_prefix = :usr

  has_secure_password validations: false

  has_many :user_sessions, dependent: :destroy

  validates :username, uniqueness: true

  def name
    "#{given_name} #{middle_name} #{family_name}".strip.gsub(/ +/, " ")
  end

  def list_name
    name.presence || username
  end

  def any_address?
    %i[street locality region postal_code country].any? { |part| send("address_#{part}") }
  end

  def pronouns
    case gender
    when "male" then "masculine"
    when "female" then "feminine"
    else super || "neutral"
    end
  end

  def next_birthday
    return unless birthdate

    date = birthdate.strftime("%m%d")
    date = "0228" if date == "0229" && !Date.today.leap?

    year = Date.today.year
    year += 1 if date < Date.today.strftime("%m%d")

    Date.parse("#{year}#{date}")
  end

  class << self
    def identify(identifier)
      return nil unless identifier.present?

      where(username: identifier).or(where(email: identifier)).first
    end
  end
end
