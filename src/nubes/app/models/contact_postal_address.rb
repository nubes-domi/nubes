class ContactPostalAddress < ApplicationRecord
  include PrettyId
  self.id_prefix = :cpaddr

  belongs_to :contact
end
