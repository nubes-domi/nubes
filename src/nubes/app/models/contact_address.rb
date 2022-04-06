class ContactAddress < ApplicationRecord
  include PrettyId
  self.id_prefix = :caddr

  belongs_to :contact
end
