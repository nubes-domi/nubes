class Contact < ApplicationRecord
  include PrettyId
  self.id_prefix = :cnt

  has_many :contact_addresses
  has_many :contact_postal_addresses

  def name
    [title, given_name, middle_name, family_name, suffix].map(&:presence).compact.join(" ")
  end
end
