class Contact < ApplicationRecord
  include PrettyId
  self.id_prefix = :cnt

  has_many :addresses, class_name: "ContactAddress"
  has_many :postal_addresses, class_name: "ContactPostalAddress"

  validates :given_name, presence: true

  ## GraphQL definition
  include Graphqlable

  field :addresses, [ContactAddress], readonly: true
  field :postal_addresses, [ContactPostalAddress], readonly: true

  field :title, String, description: "This contact title. (eg: Mr, Miss, Dr, Señor)"
  field :given_name, String, description: "Given or first name."
  field :middle_name, String, description: "Middle name"
  field :family_name, String, description: "Family name or surname"
  field :suffix, String, description: "Name suffixes. (eg: III, OBE, PhD)"
  field :nickname, String, description: "Alternative short name that this person goes by"
  field :gender, String, description: "Gender identity of this person. Not restricted to male and female."
  field :pronouns, Types::Models::Pronouns, description: "Pronoun set for this person."

  # field :birthdate, String
  # field :birthdate_has_year, String
  # field :birthdate_age_based, String

  field :formatted_name, String,
        description: "Full name, based on the other name fields. (eg: Mr Anthony Dean Smith III)",
        readonly: true

  def formatted_name
    [title, given_name, middle_name, family_name, suffix].map(&:presence).compact.join(" ")
  end
end
